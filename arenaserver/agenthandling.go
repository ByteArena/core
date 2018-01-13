package arenaserver

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/bytearena/core/common/types"
	"github.com/bytearena/core/common/utils/vector"

	arenaserveragent "github.com/bytearena/core/arenaserver/agent"
	containertypes "github.com/bytearena/core/arenaserver/container"
	uuid "github.com/satori/go.uuid"
	bettererrors "github.com/xtuc/better-errors"
)

func (s *Server) RegisterAgent(agent *types.Agent, spawningVector *vector.Vector2) {
	agentimage := agent.Manifest.Id

	///////////////////////////////////////////////////////////////////////////
	// Building the agent entity (gameplay related aspects of the agent)
	///////////////////////////////////////////////////////////////////////////

	if spawningVector == nil {

		arenamap := s.GetGameDescription().GetMapContainer()
		agentSpawnPointIndex := len(s.agentproxies)

		if agentSpawnPointIndex >= len(arenamap.Data.Starts) {
			berror := bettererrors.
				New("Cannot spawn agent").
				SetContext("image", agent.Manifest.Id).
				SetContext("number of spawns", strconv.Itoa(len(arenamap.Data.Starts))).
				With(bettererrors.New("No starting point left"))

			s.Log(EventError{berror})
			return
		}

		agentSpawningPos := arenamap.Data.Starts[agentSpawnPointIndex]

		vector := vector.MakeVector2(agentSpawningPos.Point.GetX(), agentSpawningPos.Point.GetY())
		spawningVector = &vector
	}

	agententity := s.game.NewEntityAgent(agent, *spawningVector)

	///////////////////////////////////////////////////////////////////////////
	// Building the agent proxy (concrete link with container and communication pipe)
	///////////////////////////////////////////////////////////////////////////

	agentproxy := arenaserveragent.MakeAgentProxyNetwork()
	agentproxy.SetEntityId(agententity.GetID())

	s.setAgentProxy(agentproxy)
	s.agentimages[agentproxy.GetProxyUUID()] = agentimage

	agent.EntityID = agententity.GetID()
	agent.UUID = agentproxy.GetProxyUUID()

	// Keep last spawning point in case we will respawn it (via ReloadAgent)
	s.agentspawnedvector[agentproxy.GetProxyUUID()] = spawningVector
}

func (s *Server) ReloadAgent(agent *types.Agent) error {
	// Ignore future communication error if any
	s.gameIsRunning = false

	defer func() {
		s.gameIsRunning = true
	}()

	container, hasContainer := s.agentcontainers[agent.UUID]

	if !hasContainer {
		return bettererrors.
			New("Container not found").
			SetContext("agent", agent.Manifest.Id)
	}

	// Remove from ecs
	s.game.RemoveEntityAgent(agent)

	// Close connection
	proxy, _ := s.agentproxies[agent.UUID]

	if netAgent, ok := proxy.(arenaserveragent.AgentProxyNetworkInterface); ok {

		// Remove the connection and the entity from our states
		s.removeAgentConn(netAgent.GetConn())
	}

	// Stop and remove container
	s.containerorchestrator.TearDown(container)
	s.containerorchestrator.RemoveContainer(container)

	// Wait until the containers exists to continue the process
	exist, waiterr := s.containerorchestrator.Wait(container)

	select {
	case <-exist: // ok
	case <-waiterr: // ok, probably already removed
	}

	// Re-register it
	lastSpawnedVector, _ := s.agentspawnedvector[agent.UUID]
	s.RegisterAgent(agent, lastSpawnedVector)

	// Re-start it
	newProxy, _ := s.agentproxies[agent.UUID]
	err := s.startAgentContainer(newProxy)

	if err != nil {
		return bettererrors.
			New("Could not start agent").
			SetContext("agent", agent.Manifest.Id).
			With(err)
	}

	// Wait for handshake
	s.nbhandshaked--

	return nil
}

func (s *Server) startAgentContainer(
	agentproxy arenaserveragent.AgentProxyInterface,
) error {
	dockerimage := s.agentimages[agentproxy.GetProxyUUID()]

	arenaHostnameForAgents, err := s.containerorchestrator.GetHost()

	if err != nil {
		return bettererrors.
			New("Failed to fetch arena hostname for agents").
			With(bettererrors.NewFromErr(err))
	}

	container, err1 := s.containerorchestrator.CreateAgentContainer(
		agentproxy.GetProxyUUID(),
		arenaHostnameForAgents,
		s.port,
		dockerimage,
	)

	if err1 != nil {
		return bettererrors.
			New("Failed to create docker container").
			With(err1).
			SetContext("id", agentproxy.String())
	}

	err = s.containerorchestrator.StartAgentContainer(container, s.AddTearDownCall)

	if err != nil {
		return bettererrors.
			New("Failed to start docker container").
			With(bettererrors.NewFromErr(err)).
			SetContext("id", agentproxy.String())
	}

	go func() {
		wait, err := s.containerorchestrator.Wait(container)

		select {
		case msg := <-wait:

			if !s.gameOver {
				berror := bettererrors.
					New("Agent terminated").
					SetContext("code", strconv.FormatInt(msg.StatusCode, 10))

				if msg.Error != nil {
					berror.SetContext("error", msg.Error.Message)
				}

				s.Log(EventWarn{berror})
			} else {
				s.Log(EventHeadsUp{"Agent " + container.ImageName + " has stopped."})
			}

			s.containerorchestrator.RemoveContainer(container)

			s.agentproxiesmutex.Lock()
			s.removeAgent(agentproxy.GetProxyUUID())
			s.agentproxiesmutex.Unlock()
		case <-err:
			panic(err)
		}
	}()

	go func() {

		for {
			msg := <-s.containerorchestrator.Events()

			switch t := msg.(type) {
			case containertypes.EventDebug:
				s.Log(EventLog{t.Value})
			case containertypes.EventAgentLog:
				line := fmt.Sprintf("[%s] %s", t.AgentName, t.Value)
				s.Log(EventAgentLog{line})
			default:
				msg := fmt.Sprintf("Unsupported Orchestrator message of type %s", reflect.TypeOf(msg))
				panic(msg)
			}
		}
	}()

	// Keep a ref into agentcontainers
	s.agentcontainers[agentproxy.GetProxyUUID()] = container

	return nil
}

func (s *Server) startAgentContainers() error {

	for _, agentproxy := range s.agentproxies {
		err := s.startAgentContainer(agentproxy)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) setAgentProxy(agent arenaserveragent.AgentProxyInterface) {
	s.agentproxiesmutex.Lock()
	defer s.agentproxiesmutex.Unlock()
	s.agentproxies[agent.GetProxyUUID()] = agent
}

func (s *Server) getAgentProxy(agentid string) (arenaserveragent.AgentProxyInterface, error) {
	var emptyagent arenaserveragent.AgentProxyInterface

	foundkey, err := uuid.FromString(agentid)
	if err != nil {
		return emptyagent, err
	}

	s.agentproxiesmutex.Lock()
	if foundagent, ok := s.agentproxies[foundkey]; ok {
		s.agentproxiesmutex.Unlock()
		return foundagent, nil
	}
	s.agentproxiesmutex.Unlock()

	return emptyagent, errors.New("Agent" + agentid + " not found")
}
