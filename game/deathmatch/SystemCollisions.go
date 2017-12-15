package deathmatch

import (
	"github.com/bytearena/box2d"
	"github.com/bytearena/ecs"

	commontypes "github.com/bytearena/core/common/types"
	"github.com/bytearena/core/common/utils/vector"
)

type collision struct {
	entityIDA         ecs.EntityID
	entityIDB         ecs.EntityID
	collidableAspectA *Collidable
	collidableAspectB *Collidable
	point             vector.Vector2
	collisionAngleA   float64
	collisionAngleB   float64
	// normal            vector.Vector2
	// toi               float64
	// friction          float64
	// restitution       float64
}

func systemCollisions(deathmatch *DeathmatchGame) []collision {

	collisions := make([]collision, 0)

	for _, coll := range deathmatch.collisionListener.PopCollisions() {

		A, ok := coll.GetFixtureA().GetBody().GetUserData().(commontypes.PhysicalBodyDescriptor)
		if !ok {
			continue
		}

		B, ok := coll.GetFixtureB().GetBody().GetUserData().(commontypes.PhysicalBodyDescriptor)
		if !ok {
			continue
		}

		// linearVelocityA := contact.GetFixtureA().GetBody().GetLinearVelocity()
		// linearVelocityB := contact.GetFixtureB().GetBody().GetLinearVelocity()

		worldManifold := box2d.MakeB2WorldManifold()
		coll.GetWorldManifold(&worldManifold)

		velA := vector.FromB2Vec2(coll.GetFixtureA().GetBody().GetLinearVelocityFromWorldPoint(worldManifold.Points[0]))
		velB := vector.FromB2Vec2(coll.GetFixtureB().GetBody().GetLinearVelocityFromWorldPoint(worldManifold.Points[0]))
		collisionAngleA := velB.Sub(velA).Angle()
		collisionAngleB := velA.Sub(velB).Angle()

		entityResultA := deathmatch.getEntity(A.ID, deathmatch.collidableComponent)
		entityResultB := deathmatch.getEntity(B.ID, deathmatch.collidableComponent)

		if entityResultA == nil || entityResultB == nil {
			// Should never happen; this case is filtered in deathmatch.collisionFilter
			continue
		}

		collidableAspectA := entityResultA.Components[deathmatch.collidableComponent].(*Collidable)
		collidableAspectB := entityResultB.Components[deathmatch.collidableComponent].(*Collidable)

		compiledCollision := collision{
			entityIDA:         A.ID,
			entityIDB:         B.ID,
			collidableAspectA: collidableAspectA,
			collidableAspectB: collidableAspectB,
			point:             vector.FromB2Vec2(worldManifold.Points[0]).Transform(deathmatch.physicalToAgentSpaceTransform),
			collisionAngleA:   collisionAngleA,
			collisionAngleB:   collisionAngleB,
			//normal:            vector.FromB2Vec2(worldManifold.Normal).Transform(deathmatch.physicalToAgentSpaceTransform),
			// toi:               coll.GetTOI(),
			// friction:          coll.GetFriction(),
			// restitution:       coll.GetRestitution(),
		}

		collisions = append(collisions, compiledCollision)

		collidableAspectA.CollisionScript(deathmatch, A.ID, B.ID, collidableAspectA, collidableAspectB, compiledCollision.point)
		collidableAspectB.CollisionScript(deathmatch, B.ID, A.ID, collidableAspectB, collidableAspectA, compiledCollision.point)
	}

	return collisions
}
