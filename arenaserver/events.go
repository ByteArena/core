package arenaserver

type EventStatusGameUpdate struct{ Status string }
type EventClose struct{}
type EventLog struct{ Value string }
type EventError struct{ Err error }
type EventDebug struct{ Value string }
type EventWarn struct{ Err error }
type EventAgentLog struct{ Value string }
type EventOrchestratorLog struct{ Value string }
type EventRawComm struct{ Value []byte }
