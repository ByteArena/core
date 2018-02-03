package deathmatch

import (
	"fmt"

	"github.com/bytearena/box2d"
	"github.com/bytearena/ecs"

	commontypes "github.com/bytearena/core/common/types"
	"github.com/bytearena/core/common/types/mapcontainer"
	"github.com/bytearena/core/common/utils"
)

func (deathmatch *DeathmatchGame) NewEntitySensor(polygon mapcontainer.MapPolygon, name string) *ecs.Entity {

	sensor := deathmatch.manager.NewEntity()

	bodydef := box2d.MakeB2BodyDef()
	bodydef.Type = box2d.B2BodyType.B2_staticBody

	body := deathmatch.PhysicalWorld.CreateBody(&bodydef)
	vertices := make([]box2d.B2Vec2, len(polygon.Points))

	for i := 0; i < len(polygon.Points); i++ {
		vertices[i].Set(polygon.Points[i].GetX(), polygon.Points[i].GetY()*-1) // TODO(jerome): invert axes in transform, not here
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("\n\nERROR - Sensor " + name + " is not valid; perhaps some vertices are duplicated?\n\n")
			panic(r)
		}
	}()

	prev := len(vertices) - 1
	for cur := 0; cur < len(vertices); cur++ {
		shape := box2d.MakeB2EdgeShape()
		shape.Set(vertices[prev], vertices[cur])
		body.CreateFixture(&shape, 0.0)

		prev = cur
	}

	body.SetUserData(commontypes.MakePhysicalBodyDescriptor(
		commontypes.PhysicalBodyDescriptorType.Sensor,
		sensor.GetID(),
	))

	return sensor.
		AddComponent(deathmatch.physicalBodyComponent, &PhysicalBody{
			body:   body,
			static: true,
		}).
		AddComponent(deathmatch.collidableComponent, &Collidable{
			collisiongroup: CollisionGroup.Ground,
			collideswith: utils.BuildTag(
				CollisionGroup.Agent,
			),
		}) //.
	// AddComponent(deathmatch.sensorComponent, &Sensor{
	// 	onSensed: func() {
	// 		log.Println("SENSED!!!")
	// 	},
	// })
}
