package agent

import (
	"github.com/bytearena/core/common/types"
)

type AgentProxyLocalInterface interface {
	AgentProxyInterface
}

type AgentProxyLocal struct {
	AgentProxyGeneric
	DebugNbPutPerception int
}

func MakeLocalAgentImp() AgentProxyLocal {
	return AgentProxyLocal{
		AgentProxyGeneric: MakeAgentProxyGeneric(),
	}
}

func (agent AgentProxyLocal) String() string {
	return "<LocalAgentImp(" + agent.GetProxyUUID().String() + ")>"
}

func (agent AgentProxyLocal) SetPerception(perception []byte, comm types.AgentCommunicatorInterface) error {
	return nil
}
