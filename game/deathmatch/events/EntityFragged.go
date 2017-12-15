package events

import "github.com/bytearena/ecs"

type EntityFragged struct {
	Entity    ecs.EntityID
	FraggedBy ecs.EntityID
}

func (ev EntityFragged) Topic() string { return "gameplay:entity:fragged" }
