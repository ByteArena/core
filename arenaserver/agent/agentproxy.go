package agent

import (
	uuid "github.com/satori/go.uuid"

	"github.com/bytearena/core/common/types"
	"github.com/bytearena/ecs"
)

type AgentProxyInterface interface {
	GetProxyUUID() uuid.UUID
	GetEntityId() ecs.EntityID
	SetPerception(perceptionjson []byte, comm types.AgentCommunicatorInterface) error // abstract method
	SendAgentWelcome(message []byte, comm types.AgentCommunicatorInterface) error     // abstract method
	String() string
}

type AgentProxyGeneric struct {
	proxyUUID uuid.UUID
	entityID  ecs.EntityID
}

func MakeAgentProxyGeneric() AgentProxyGeneric {
	return AgentProxyGeneric{
		proxyUUID: uuid.NewV4(), // random uuid
	}
}

func (agent AgentProxyGeneric) GetProxyUUID() uuid.UUID {
	return agent.proxyUUID
}

func (agent *AgentProxyGeneric) SetEntityId(id ecs.EntityID) {
	agent.entityID = id
}

func (agent AgentProxyGeneric) GetEntityId() ecs.EntityID {
	return agent.entityID
}

func (agent AgentProxyGeneric) String() string {
	return "<AgentImp(" + agent.GetProxyUUID().String() + ")>"
}

func (agent AgentProxyGeneric) SetPerception(perceptionjson []byte, comm types.AgentCommunicatorInterface) error {
	// I'm abstract, override me !
	return nil
}

func (agent AgentProxyGeneric) SendAgentWelcome(bytes []byte, comm types.AgentCommunicatorInterface) error {
	return nil
}
