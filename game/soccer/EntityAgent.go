package soccer

import (
	"github.com/bytearena/box2d"
	"github.com/bytearena/ecs"

	"github.com/bytearena/core/common/types"
	"github.com/bytearena/core/common/utils"
	"github.com/bytearena/core/common/utils/number"
	"github.com/bytearena/core/common/utils/vector"
)

func (game *SoccerGame) NewEntityAgent(
	agent *types.Agent,
	spawnPosition vector.Vector2, // spawnPosition in physical space; TODO: fix this, should be in agent space
) ecs.EntityID {
	agentEntity := game.manager.NewEntity()

	///////////////////////////////////////////////////////////////////////////
	// Définition de ses caractéristiques physiques de l'agent (spécifications)
	///////////////////////////////////////////////////////////////////////////

	// Linear unit expressed in agent space units (meters) per tick
	// Angular unit expressed in radians per tick

	bodyRadius := 0.5
	maxSpeed := 1.25
	maxSteering := 10000.0
	dragForce := 0.015
	maxAngularVelocity := number.DegreeToRadian(15.0)

	///////////////////////////////////////////////////////////////////////////
	// Création du corps physique de l'agent (Box2D)
	///////////////////////////////////////////////////////////////////////////

	bodydef := box2d.MakeB2BodyDef()
	bodydef.Type = box2d.B2BodyType.B2_dynamicBody
	bodydef.AllowSleep = false
	bodydef.FixedRotation = true

	body := game.PhysicalWorld.CreateBody(&bodydef)

	shape := box2d.MakeB2CircleShape()
	shape.SetRadius(bodyRadius * game.agentSpaceToPhysicalSpaceScale)

	fixturedef := box2d.MakeB2FixtureDef()
	fixturedef.Shape = &shape
	fixturedef.Density = 20.0
	body.CreateFixtureFromDef(&fixturedef)
	body.SetUserData(types.MakePhysicalBodyDescriptor(
		types.PhysicalBodyDescriptorType.Agent,
		agentEntity.GetID(),
	))
	body.SetBullet(false)

	///////////////////////////////////////////////////////////////////////////
	// Composition de l'agent dans l'ECS
	///////////////////////////////////////////////////////////////////////////
	tps := game.TPS

	agentEntity.
		AddComponent(game.physicalBodyComponent, (&PhysicalBody{
			body:               body,
			maxSpeed:           maxSpeed,
			maxAngularVelocity: maxAngularVelocity,
			dragForce:          dragForce,

			pointTransformIn:  game.physicalToAgentSpaceInverseTransform,
			pointTransformOut: game.physicalToAgentSpaceTransform,

			distanceScaleIn:  game.physicalToAgentSpaceInverseScale, // same as transform matrix, but scale only (for 1D transforms of length)
			distanceScaleOut: game.physicalToAgentSpaceScale,        // same as transform matrix, but scale only (for 1D transforms of length)

			timeScaleIn:  float64(tps),       // m/tick to m/s; => ticksPerSecond
			timeScaleOut: 1.0 / float64(tps), // m/s to m/tick; => 1 / ticksPerSecond
		}).SetPositionInPhysicalScale(
			//vector.MakeVector2(spawnPosition.GetX(), -1*spawnPosition.GetY()), // TODO(jerome): invert axes in transform, not here
			vector.MakeVector2(58.5, -1*39.75), // TODO: SOCCER: this is the center of the field
		)).
		AddComponent(game.perceptionComponent, &Perception{
			perception: newEmptyAgentPerception(),
		}).
		AddComponent(game.playerComponent, &Player{
			Agent: agent,
		}).
		AddComponent(game.steeringComponent, NewSteering(
			maxSteering, // MaxSteering
		)).
		AddComponent(game.renderComponent, &Render{
			type_:  "agent",
			static: false,
		}).
		AddComponent(game.collidableComponent, &Collidable{
			collisiongroup: CollisionGroup.Agent,
			collideswith: utils.BuildTag(
				CollisionGroup.Agent,
				CollisionGroup.Obstacle,
				CollisionGroup.Ball,
			),
			collisionScriptFunc: agentCollisionScript,
		})

	return agentEntity.GetID()
}

func (game *SoccerGame) RemoveEntityAgent(agent *types.Agent) {

	qr := game.getEntity(agent.EntityID)
	game.manager.DisposeEntity(qr)
}

func agentCollisionScript(game *SoccerGame, entityID ecs.EntityID, otherEntityID ecs.EntityID, collidableAspect *Collidable, otherCollidableAspectB *Collidable, point vector.Vector2) {
	entityResult := game.getEntity(entityID, game.physicalBodyComponent)
	if entityResult == nil {
		return
	}

	physicalAspect := entityResult.Components[game.physicalBodyComponent].(*PhysicalBody)
	physicalAspect.SetVelocity(vector.MakeNullVector2())
}
