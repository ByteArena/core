package events

import "github.com/bytearena/ecs"

type EntityRespawning struct {
	Entity     ecs.EntityID
	RespawnsIn int
}

func (ev EntityRespawning) Topic() string { return "gameplay:entity:respawning" }
