package arenaserver

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	notify "github.com/bitly/go-notify"
	"github.com/bytearena/core/arenaserver/agent"
	"github.com/bytearena/core/arenaserver/comm"
	"github.com/phayes/freeport"
	uuid "github.com/satori/go.uuid"

	bettererrors "github.com/xtuc/better-errors"

	"github.com/bytearena/core/common/mq"
	"github.com/bytearena/core/common/types"
	"github.com/bytearena/core/common/types/mapcontainer"
	"github.com/bytearena/core/common/utils"
	commongame "github.com/bytearena/core/game/common"
)

const (
	POURCENT_LEFT_BEFORE_QUIT    = 50 // %
	CLOSE_CONNECTION_BEFORE_KILL = 1 * time.Second
	LOG_ENTRY_BUFFER             = 100
)

type Server struct {
	host            string
	port            int
	arenaServerUUID string
	tickspersec     int

	stopticking  chan bool
	nbhandshaked int

	currentturn uint32

	tearDownCallbacks      []types.TearDownCallback
	tearDownCallbacksMutex *sync.Mutex

	containerorchestrator types.ContainerOrchestrator
	commserver            *comm.CommServer
	mqClient              mq.ClientInterface

	gameDescription types.GameDescriptionInterface

	agentproxies           map[uuid.UUID]agent.AgentProxyInterface
	agentproxiesmutex      *sync.Mutex
	agentproxieshandshakes map[uuid.UUID]struct{}
	agentimages            map[uuid.UUID]string

	pendingmutations []types.AgentMutationBatch
	mutationsmutex   *sync.Mutex

	tickdurations []int64

	///////////////////////////////////////////////////////////////////////
	// Game logic
	///////////////////////////////////////////////////////////////////////

	game          commongame.GameInterface
	gameIsRunning bool
	gameDuration  *time.Duration
	gameStartTime *time.Time
	gameOver      bool

	isDebug bool

	events chan interface{}
}

func NewServer(
	host string,
	orch types.ContainerOrchestrator,
	gameDescription types.GameDescriptionInterface,
	game commongame.GameInterface,
	arenaServerUUID string,
	mqClient mq.ClientInterface,
	gameDuration *time.Duration,
	isDebug bool,
) *Server {

	gamehost := host

	if host == "" {
		host, err := orch.GetHost()
		utils.Check(err, "Could not determine arena-server host/ip.")

		gamehost = host
	}

	port, err := freeport.GetFreePort()
	utils.Check(err, "Unable to allocate a port") // Fatal

	tickspersec := gameDescription.GetTps()

	s := &Server{
		host:            gamehost,
		port:            port,
		arenaServerUUID: arenaServerUUID,
		tickspersec:     tickspersec,

		stopticking:  make(chan bool),
		nbhandshaked: 0,

		tearDownCallbacks:      make([]types.TearDownCallback, 0),
		tearDownCallbacksMutex: &sync.Mutex{},

		containerorchestrator: orch,
		commserver:            nil, // initialized in Listen()
		mqClient:              mqClient,

		gameDescription: gameDescription,

		// agents here: proxy to agent in container
		agentproxies:           make(map[uuid.UUID]agent.AgentProxyInterface),
		agentproxiesmutex:      &sync.Mutex{},
		agentproxieshandshakes: make(map[uuid.UUID]struct{}),
		agentimages:            make(map[uuid.UUID]string),

		pendingmutations: make([]types.AgentMutationBatch, 0),
		mutationsmutex:   &sync.Mutex{},

		tickdurations: make([]int64, 0),

		///////////////////////////////////////////////////////////////////////
		// Game logic
		///////////////////////////////////////////////////////////////////////

		game:          game,
		gameIsRunning: false,
		gameDuration:  gameDuration,
		gameStartTime: nil,
		gameOver:      false,

		events: make(chan interface{}, LOG_ENTRY_BUFFER),

		isDebug: isDebug,
	}

	return s
}

func (s Server) getNbExpectedagents() int {
	return len(s.GetGameDescription().GetAgents())
}

///////////////////////////////////////////////////////////////////////////////
// Public API
///////////////////////////////////////////////////////////////////////////////

func (server *Server) Start() (chan interface{}, error) {

	block := server.listen()
	err := server.startAgentContainers()

	if err != nil {
		return nil, bettererrors.New("Failed to start agent containers").With(err)
	}

	server.AddTearDownCall(func() error {
		//server.Log(EventLog{"Publish game state (" + server.arenaServerUUID + "stopped)"})

		game := server.GetGameDescription()

		err := server.mqClient.Publish("game", "stopped", types.NewMQMessage(
			"arena-server",
			"Arena Server "+server.arenaServerUUID+", game "+game.GetId()+" stopped",
		).SetPayload(types.MQPayload{
			"id":              game.GetId(),
			"arenaserveruuid": server.arenaServerUUID,
		}))

		return err
	})

	return block, nil
}

func (server *Server) Stop() {
	server.gameIsRunning = false

	server.Log(EventDebug{"TearDown from stop"})
	server.TearDown()
}

func (s *Server) SubscribeStateObservation() chan interface{} {
	ch := make(chan interface{})
	notify.Start("app:stateupdated", ch)
	return ch
}

func (s *Server) SendLaunched() {
	payload := types.MQPayload{
		"id":              s.GetGameDescription().GetId(),
		"arenaserveruuid": s.arenaServerUUID,
	}

	s.mqClient.Publish("game", "launched", types.NewMQMessage(
		"arena-server",
		"Arena Server "+s.arenaServerUUID+" launched",
	).SetPayload(payload))

	payloadJson, _ := json.Marshal(payload)

	s.Log(EventLog{"Send game launched: " + string(payloadJson)})
}

func (s Server) GetGameDescription() types.GameDescriptionInterface {
	return s.gameDescription
}

func (s Server) GetGame() commongame.GameInterface {
	return s.game
}

func (s Server) GetTicksPerSecond() int {
	return s.tickspersec
}

func (server *Server) onAgentsReady() {
	server.Log(EventLog{"Agents are ready; starting in 100 ms"})
	time.Sleep(time.Duration(time.Millisecond * 100))

	server.startTicking()
}

func (server *Server) startTicking() {

	server.gameIsRunning = true
	now := time.Now()
	server.gameStartTime = &now

	tickduration := time.Duration((1000000 / time.Duration(server.tickspersec)) * time.Microsecond)
	ticker := time.Tick(tickduration)

	server.AddTearDownCall(func() error {
		server.stopticking <- true
		close(server.stopticking)
		return nil
	})

	go func() {
		for {
			<-ticker

			if server.gameOver {
				return
			}

			server.doTick()
		}
	}()

	if server.gameDuration != nil {
		server.Log(EventHeadsUp{"Game will run for " + server.gameDuration.String()})
		go func() {
			<-time.After(*server.gameDuration)
			server.gameOver = true
			server.Log(EventHeadsUp{"Game ended after " + server.gameDuration.String()})
			notify.Post("app:stopticking", true) // gameover: true
		}()
	} else {
		server.Log(EventHeadsUp{"Game will run indefinitely"})
	}

	go func() {
		<-server.stopticking
		server.gameOver = true
		server.Log(EventLog{"Received stop ticking signal"})
		notify.Post("app:stopticking", false) // gameover: false
	}()
}

func (server *Server) popMutationBatches() []types.AgentMutationBatch {
	server.mutationsmutex.Lock()
	mutations := server.pendingmutations
	server.pendingmutations = make([]types.AgentMutationBatch, 0)
	server.mutationsmutex.Unlock()

	return mutations
}

func (server *Server) doTick() {

	//watch := utils.MakeStopwatch("doTick")
	//watch.Start("global")

	begin := time.Now()

	turn := int(server.currentturn) // starts at 0
	atomic.AddUint32(&server.currentturn, 1)

	dolog := (turn%server.tickspersec) == 0 || server.isDebug

	///////////////////////////////////////////////////////////////////////////
	// Updating Game
	///////////////////////////////////////////////////////////////////////////

	timeStep := 1.0 / float64(server.GetTicksPerSecond())
	mutations := server.popMutationBatches()
	server.game.Step(turn, timeStep, mutations)

	///////////////////////////////////////////////////////////////////////////
	// Refreshing perception for every agent
	///////////////////////////////////////////////////////////////////////////

	arenamap := server.GetGameDescription().GetMapContainer()
	for _, agentproxy := range server.agentproxies {
		go func(server *Server, agentproxy agent.AgentProxyInterface, arenamap *mapcontainer.MapContainer) {

			err := agentproxy.SetPerception(
				server.GetGame().GetAgentPerception(agentproxy.GetEntityId()),
				server,
			)

			if err != nil && server.gameIsRunning {
				berror := bettererrors.
					New("Failed to send perception").
					SetContext("agent", agentproxy.GetProxyUUID().String()).
					With(bettererrors.NewFromErr(err))

				server.Log(EventError{berror})
			}

		}(server, agentproxy, arenamap)
	}

	///////////////////////////////////////////////////////////////////////////
	// Pushing updated state to viz
	///////////////////////////////////////////////////////////////////////////

	//watch.Stop("global")
	//fmt.Println(watch.String())

	notify.Post("app:stateupdated", nil)

	var lastduration int64 = time.Now().UnixNano() - begin.UnixNano()
	nbsamplesToKeep := server.GetTicksPerSecond() * 1
	if len(server.tickdurations) < nbsamplesToKeep {
		server.tickdurations = append(server.tickdurations, lastduration)
	} else {
		server.tickdurations[turn%nbsamplesToKeep] = lastduration
	}

	if dolog {
		var totalDuration int64 = 0
		for _, duration := range server.tickdurations {
			totalDuration += duration
		}
		meanTick := float64(totalDuration) / float64(len(server.tickdurations))
		server.Log(EventStatusGameUpdate{fmt.Sprintf(
			"Tick %d; %.3f ms mean; %.3f ms last; %d goroutines",
			turn,
			meanTick/1000000.0,
			float64(lastduration)/1000000.0,
			runtime.NumGoroutine(),
		)})
	}

}

func (s *Server) AddTearDownCall(fn types.TearDownCallback) {
	s.tearDownCallbacksMutex.Lock()
	defer s.tearDownCallbacksMutex.Unlock()

	s.tearDownCallbacks = append(s.tearDownCallbacks, fn)
}

func (server *Server) closeAllAgentConnections() {
	server.agentproxiesmutex.Lock()

	var wg sync.WaitGroup

	for _, agentproxy := range server.agentproxies {
		if netAgent, ok := agentproxy.(agent.AgentProxyNetworkInterface); ok {

			if conn := netAgent.GetConn(); conn != nil {
				wg.Add(1)

				err := conn.Close()
				<-time.After(CLOSE_CONNECTION_BEFORE_KILL)

				wg.Done()

				if err != nil {
					berror := bettererrors.
						New("Could not close agent connection").
						With(bettererrors.NewFromErr(err))

					server.Log(EventWarn{berror})
				}
			}
		}
	}

	wg.Wait()

	server.agentproxiesmutex.Unlock()
}

func (server *Server) TearDown() {
	server.events <- EventDebug{"teardown"}

	server.tearDownCallbacksMutex.Lock()

	for i := len(server.tearDownCallbacks) - 1; i >= 0; i-- {
		//server.events <- EventLog{"Executing TearDownCallback"}
		server.tearDownCallbacks[i]()
	}

	// Reset to avoid calling teardown callback multiple times
	server.tearDownCallbacks = make([]types.TearDownCallback, 0)
	server.tearDownCallbacksMutex.Unlock()

	// Close communication with agents
	server.closeAllAgentConnections()

	// Stop running container
	server.containerorchestrator.TearDownAll()

	server.events <- EventClose{}
}

func (server *Server) Events() chan interface{} {
	return server.events
}

func (server *Server) Log(l interface{}) {
	select {
	case server.events <- l:
		{
		}
	default:
		fmt.Println("[gameserver] Log dropped because buffer full")
	}
}
