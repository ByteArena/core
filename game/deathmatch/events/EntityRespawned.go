package events

import "github.com/bytearena/ecs"

type EntityRespawned struct {
	Entity        ecs.EntityID
	StartingPoint [2]float64
}

func (ev EntityRespawned) Topic() string { return "gameplay:entity:respawned" }
