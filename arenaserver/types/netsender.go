package types

import "net"

// Interface to avoid circular dependencies between server and agent

type NetSenderInterface interface {
	NetSend(message []byte, conn net.Conn) error
}
