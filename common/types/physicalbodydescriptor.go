package types

import "github.com/bytearena/ecs"

// PhysicalBodyDescriptor is set as UserData on Box2D Physical bodies to be able to determine collider and collidee from Box2D contact callbacks
type PhysicalBodyDescriptor struct {
	Type string
	ID   ecs.EntityID
}

var PhysicalBodyDescriptorType = struct {
	Obstacle   string
	Agent      string
	Ground     string
	Projectile string
}{
	Obstacle:   "o",
	Agent:      "a",
	Ground:     "g",
	Projectile: "p",
}

func MakePhysicalBodyDescriptor(type_ string, id ecs.EntityID) PhysicalBodyDescriptor {
	return PhysicalBodyDescriptor{
		Type: type_,
		ID:   id,
	}
}
