package deathmatch

import (
	"fmt"

	"github.com/bytearena/box2d"
	"github.com/bytearena/ecs"

	commontypes "github.com/bytearena/core/common/types"
	"github.com/bytearena/core/common/types/mapcontainer"
	"github.com/bytearena/core/common/utils"
)

func (deathmatch *DeathmatchGame) NewEntityGround(polygon mapcontainer.MapPolygon, name string) *ecs.Entity {
	return newEntityGroundOrObstacle(deathmatch, polygon, commontypes.PhysicalBodyDescriptorType.Ground, name).
		AddComponent(deathmatch.collidableComponent, &Collidable{
			collisiongroup: CollisionGroup.Ground,
			collideswith: utils.BuildTag(
				CollisionGroup.Agent,
			),
		})
}

func (deathmatch *DeathmatchGame) NewEntityObstacle(polygon mapcontainer.MapPolygon, name string) *ecs.Entity {
	return newEntityGroundOrObstacle(deathmatch, polygon, commontypes.PhysicalBodyDescriptorType.Obstacle, name).
		AddComponent(deathmatch.collidableComponent, &Collidable{
			collisiongroup: CollisionGroup.Obstacle,
			collideswith: utils.BuildTag(
				CollisionGroup.Agent,
				CollisionGroup.Projectile,
			),
		})
}

func newEntityGroundOrObstacle(deathmatch *DeathmatchGame, polygon mapcontainer.MapPolygon, obstacletype string, name string) *ecs.Entity {

	obstacle := deathmatch.manager.NewEntity()

	bodydef := box2d.MakeB2BodyDef()
	bodydef.Type = box2d.B2BodyType.B2_staticBody

	body := deathmatch.PhysicalWorld.CreateBody(&bodydef)
	vertices := make([]box2d.B2Vec2, len(polygon.Points))

	for i := 0; i < len(polygon.Points); i++ {
		vertices[i].Set(polygon.Points[i].GetX(), polygon.Points[i].GetY()*-1) // TODO(jerome): invert axes in transform, not here
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("\n\nERROR - Obstacle or ground (type " + obstacletype + ") " + name + " is not valid; perhaps some vertices are duplicated?\n\n")
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
		obstacletype,
		obstacle.GetID(),
	))

	return obstacle.
		AddComponent(deathmatch.physicalBodyComponent, &PhysicalBody{
			body:   body,
			static: true,
		})
}
