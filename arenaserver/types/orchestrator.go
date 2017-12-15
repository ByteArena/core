package types

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	uuid "github.com/satori/go.uuid"

	t "github.com/bytearena/core/common/types"
)

type ContainerOrchestrator interface {
	StartAgentContainer(ctner *AgentContainer, addTearDownCall func(t.TearDownCallback)) error
	RemoveAgentContainer(ctner *AgentContainer) error
	Wait(ctner *AgentContainer) (<-chan container.ContainerWaitOKBody, <-chan error)
	TearDown(container *AgentContainer)
	CreateAgentContainer(agentid uuid.UUID, host string, port int, dockerimage string) (*AgentContainer, error)
	GetHost() (string, error)
	SetAgentLogger(container *AgentContainer) error
	TearDownAll() error
	GetCli() *client.Client
	GetContext() context.Context
	GetRegistryAuth() string
	AddContainer(*AgentContainer)
	RemoveContainer(*AgentContainer)
	Events() chan interface{}
}
