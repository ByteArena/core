package deathmatch

import (
	"math"

	"github.com/bytearena/box2d"
	"github.com/bytearena/ecs"

	"github.com/bytearena/core/common/types"
	"github.com/bytearena/core/common/utils"
	"github.com/bytearena/core/common/utils/vector"
)

func (deathmatch *DeathmatchGame) NewEntityBallisticProjectile(ownerid ecs.EntityID, position vector.Vector2, velocity vector.Vector2) *ecs.Entity {

	ownerAspects := deathmatch.getEntity(ownerid,
		deathmatch.shootingComponent,
	)

	if ownerAspects == nil {
		// Should never happen
		return nil
	}

	///////////////////////////////////////////////////////////////////////////
	tps := deathmatch.gameDescription.GetTps()

	timeScaleIn := float64(tps)
	timeScaleOut := 1 / timeScaleIn
	///////////////////////////////////////////////////////////////////////////

	shootingAspect := ownerAspects.Components[deathmatch.shootingComponent].(*Shooting)

	bodyRadius := 0.3                                   // meters
	projectilespeed := shootingAspect.ProjectileSpeed   // m/tick
	projectiledamage := shootingAspect.ProjectileDamage // amount of life consumed on impact
	projectilerange := shootingAspect.ProjectileRange   // in meter

	projectilettl := 0

	if projectilespeed > 0 {
		projectilettl = int(math.Ceil(projectilerange / projectilespeed))
	}

	projectile := deathmatch.manager.NewEntity()

	bodydef := box2d.MakeB2BodyDef()
	bodydef.Type = box2d.B2BodyType.B2_dynamicBody
	bodydef.AllowSleep = true
	bodydef.FixedRotation = true

	physicalReferentialPosition := position.Transform(deathmatch.physicalToAgentSpaceInverseTransform)
	bodydef.Position.Set(physicalReferentialPosition.GetX(), physicalReferentialPosition.GetY())

	physicalReferentialVelocity := velocity.
		SetMag(projectilespeed).
		Scale(timeScaleIn).
		Transform(deathmatch.physicalToAgentSpaceInverseTransform)

	bodydef.LinearVelocity = box2d.MakeB2Vec2(physicalReferentialVelocity.GetX(), physicalReferentialVelocity.GetY())

	body := deathmatch.PhysicalWorld.CreateBody(&bodydef)
	body.SetBullet(true)
	body.SetLinearDamping(0.0) // no aerodynamic drag
	body.SetUserData(types.MakePhysicalBodyDescriptor(
		types.PhysicalBodyDescriptorType.Projectile,
		projectile.GetID(),
	))

	shape := box2d.MakeB2CircleShape()
	shape.SetRadius(bodyRadius * deathmatch.physicalToAgentSpaceInverseScale)

	fixturedef := box2d.MakeB2FixtureDef()
	fixturedef.Shape = &shape
	fixturedef.Density = 20.0

	body.CreateFixtureFromDef(&fixturedef)

	return projectile.
		AddComponent(deathmatch.physicalBodyComponent, &PhysicalBody{
			body:               body,
			maxSpeed:           projectilespeed,
			maxAngularVelocity: 0,
			dragForce:          0,

			pointTransformIn:  deathmatch.physicalToAgentSpaceInverseTransform,
			pointTransformOut: deathmatch.physicalToAgentSpaceTransform,

			distanceScaleIn:  deathmatch.physicalToAgentSpaceInverseScale, // same as transform matrix, but scale only (for 1D transforms of length)
			distanceScaleOut: deathmatch.physicalToAgentSpaceScale,        // same as transform matrix, but scale only (for 1D transforms of length)

			timeScaleIn:  timeScaleIn,  // m/tick to m/s; => ticksPerSecond
			timeScaleOut: timeScaleOut, // m/s to m/tick; => 1 / ticksPerSecond
			skipThisTurn: true,
		}).
		AddComponent(deathmatch.renderComponent, &Render{
			type_:  "projectile",
			static: false,
		}).
		AddComponent(deathmatch.lifecycleComponent, &Lifecycle{
			tickBirth: deathmatch.ticknum,
			maxAge:    projectilettl,
		}).
		AddComponent(deathmatch.ownedComponent, &Owned{ownerid}).
		AddComponent(deathmatch.impactorComponent, &Impactor{
			damage: projectiledamage,
		}).
		AddComponent(deathmatch.collidableComponent, NewCollidable(
			CollisionGroup.Projectile,
			utils.BuildTag(
				CollisionGroup.Agent,
				CollisionGroup.Obstacle,
				CollisionGroup.Projectile,
			),
		).SetCollisionScriptFunc(projectileCollisionScript))
}

func projectileCollisionScript(game *DeathmatchGame, entityID ecs.EntityID, otherEntityID ecs.EntityID, collidableAspect *Collidable, otherCollidableAspectB *Collidable, point vector.Vector2) {
	entityResult := game.getEntity(entityID, game.physicalBodyComponent, game.lifecycleComponent)
	if entityResult == nil {
		return
	}

	physicalAspect := entityResult.Components[game.physicalBodyComponent].(*PhysicalBody)
	lifecycleAspect := entityResult.Components[game.lifecycleComponent].(*Lifecycle)

	physicalAspect.
		SetVelocity(vector.MakeNullVector2()).
		SetPosition(point)

	lifecycleAspect.SetDeath(game.ticknum) // dead in this tick
}
