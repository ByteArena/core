package events

import (
	"github.com/bytearena/ecs"
)

type EntityHit struct {
	Entity     ecs.EntityID
	HitBy      ecs.EntityID
	ComingFrom float64 // absolute azimuth in radian
	Damage     float64 // absolute azimuth in radian
}

func (ev EntityHit) Topic() string { return "gameplay:entity:hit" }
