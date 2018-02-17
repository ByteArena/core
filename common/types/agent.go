package types

import (
	"github.com/bytearena/ecs"
	uuid "github.com/satori/go.uuid"
)

type Agent struct {
	Manifest AgentManifest `json:"manifest"`
	EntityID ecs.EntityID  `json:"id"`

	UUID uuid.UUID `json:"-"`
}
