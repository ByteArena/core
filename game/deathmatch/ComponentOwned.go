package deathmatch

import "github.com/bytearena/ecs"

type Owned struct {
	owner ecs.EntityID
}

func (o Owned) GetOwner() ecs.EntityID {
	return o.owner
}

func (o *Owned) SetOwner(owner ecs.EntityID) *Owned {
	o.owner = owner
	return o
}
