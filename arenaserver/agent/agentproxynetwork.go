package agent

import (
	"net"

	"github.com/bytearena/core/common/types"
)

type AgentProxyNetworkInterface interface {
	AgentProxyInterface
	SetConn(conn net.Conn) AgentProxyNetworkInterface
	GetConn() net.Conn
}

type AgentProxyNetwork struct {
	AgentProxyGeneric
	conn net.Conn
}

func MakeAgentProxyNetwork() AgentProxyNetwork {
	return AgentProxyNetwork{
		AgentProxyGeneric: MakeAgentProxyGeneric(),
	}
}

func (agent AgentProxyNetwork) String() string {
	return "<NetAgentImp(" + agent.GetProxyUUID().String() + ")>"
}

func (agent AgentProxyNetwork) SetPerception(perceptionjson []byte, comm types.AgentCommunicatorInterface) error {
	message := []byte("{\"method\":\"perception\",\"payload\":" + string(perceptionjson) + "}\n")
	return comm.NetSend(message, agent.GetConn())
}

func (agent AgentProxyNetwork) SendAgentWelcome(bytes []byte, comm types.AgentCommunicatorInterface) error {
	message := []byte("{\"method\":\"welcome\",\"payload\":" + string(bytes) + "}\n")
	return comm.NetSend(message, agent.GetConn())
}

func (agent AgentProxyNetwork) SetConn(conn net.Conn) AgentProxyNetworkInterface {
	agent.conn = conn
	return agent
}

func (agent AgentProxyNetwork) GetConn() net.Conn {
	return agent.conn
}
