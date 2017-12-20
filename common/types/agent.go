package types

import (
	"github.com/bytearena/ecs"
)

type Agent struct {
	Manifest AgentManifest `json:"manifest"`
	EntityID ecs.EntityID  `json:"id"`
}
