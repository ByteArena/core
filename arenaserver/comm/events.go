package comm

import "net"

type EventLog struct{ Value string }
type EventError struct{ Err error }
type EventWarn struct{ Err error }

type EventRawComm struct {
	Buffer []byte
	From   string
}

type EventConnDisconnected struct {
	Err  error
	Conn net.Conn
}
