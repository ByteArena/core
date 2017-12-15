package common

import (
	"github.com/bytearena/ecs"

	"github.com/bytearena/core/arenaserver/types"
	commontypes "github.com/bytearena/core/common/types"
	"github.com/bytearena/core/common/utils/vector"
)

type GameEventSubscription int32

type GameInterface interface {
	ImplementsGameInterface()
	Subscribe(event string, cbk func(data interface{})) GameEventSubscription
	Unsubscribe(subscription GameEventSubscription)
	Step(tickturn int, dt float64, mutations []types.AgentMutationBatch)
	NewEntityAgent(contestant commontypes.Agent, pos vector.Vector2) *ecs.Entity

	GetAgentPerception(entityid ecs.EntityID) []byte
	GetAgentWelcome(entityid ecs.EntityID) []byte

	GetVizFrameJson() []byte
}
