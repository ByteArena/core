package visualization

import (
	"context"
	"net"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/bytearena/schnapps/utils"

	"github.com/bytearena/core/common/mappack"
	"github.com/bytearena/core/common/recording"
	apphandler "github.com/bytearena/core/common/visualization/handler"
	"github.com/bytearena/core/common/visualization/types"
)

type FetchArenasCbk func() ([]*types.VizGame, error)
type EventLog struct{ Value string }

type VizService struct {
	addr          string
	webclientpath string
	mapkey        string
	fetchGames    FetchArenasCbk
	listener      *http.Server
	recordStore   recording.RecordStoreInterface
	mappack       *mappack.MappackInMemoryArchive

	events chan interface{}
}

const (
	LOG_ENTRY_BUFFER = 100
)

func NewVizService(addr string, webclientpath string, mapkey string, fetchArenas FetchArenasCbk, recordStore recording.RecordStoreInterface, mappack *mappack.MappackInMemoryArchive) *VizService {
	return &VizService{
		addr:          addr,
		webclientpath: webclientpath,
		mapkey:        mapkey,
		fetchGames:    fetchArenas,
		recordStore:   recordStore,
		mappack:       mappack,

		events: make(chan interface{}, LOG_ENTRY_BUFFER),
	}
}

type logger struct {
	LogFn func(interface{})
}

func (l logger) Write(p []byte) (n int, err error) {
	l.LogFn(EventLog{string(p)})

	return len(p), nil
}

func (viz *VizService) Start() chan struct{} {
	logger := logger{viz.Log}
	router := mux.NewRouter()

	router.PathPrefix("/mappack/").Handler(handlers.CombinedLoggingHandler(
		logger,
		http.StripPrefix("/mappack/", viz.mappack),
	))

	router.Handle("/", handlers.CombinedLoggingHandler(
		logger,
		http.HandlerFunc(apphandler.Home(viz.fetchGames)),
	)).Methods("GET")

	router.Handle("/record/{recordId:[a-zA-Z0-9\\-]+}", handlers.CombinedLoggingHandler(
		logger,
		http.HandlerFunc(apphandler.Replay(viz.recordStore, viz.webclientpath)),
	)).Methods("GET")

	router.Handle("/record/{recordId:[a-zA-Z0-9\\-]+}/ws", handlers.CombinedLoggingHandler(
		logger,
		http.HandlerFunc(apphandler.ReplayWebsocket(viz.recordStore, viz.webclientpath)),
	)).Methods("GET")

	router.Handle("/arena/{id:[a-zA-Z0-9\\-]+}", handlers.CombinedLoggingHandler(
		logger,
		http.HandlerFunc(apphandler.Game(viz.fetchGames, viz.mappack)),
	)).Methods("GET")

	router.Handle("/arena/{id:[a-zA-Z0-9\\-]+}/ws", handlers.CombinedLoggingHandler(
		logger,
		http.HandlerFunc(apphandler.Websocket(viz.fetchGames)),
	)).Methods("GET")

	viz.Log(EventLog{"VIZ Listening on " + viz.addr})

	listener, err := net.Listen("tcp4", viz.addr)
	if err != nil {
		utils.Check(err, err.Error())
	}

	viz.listener = &http.Server{
		Handler: router,
	}

	block := make(chan struct{})

	go func(block chan struct{}) {
		_ = viz.listener.Serve(listener)
		//utils.Check(err, "Failed to listen on "+viz.addr)
		close(block)
	}(block)

	return block
}

func (viz *VizService) Stop() {
	viz.mappack.Close()
	viz.listener.Shutdown(context.TODO())
}

func (viz *VizService) Events() chan interface{} {
	return viz.events
}

func (viz *VizService) Log(x interface{}) {
	viz.events <- x
}
