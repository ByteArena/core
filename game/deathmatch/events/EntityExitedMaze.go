package events

import "github.com/bytearena/ecs"

type EntityExitedMaze struct {
	Entity ecs.EntityID
	Exit   ecs.EntityID
}

func (ev EntityExitedMaze) Topic() string { return "gameplay:entity:exitedmaze" }
