package soccer

import (
	"math"

	"github.com/bytearena/ecs"
)

var pi2 = math.Pi * 2
var halfpi = math.Pi / 2
var threepi2 = math.Pi + halfpi

// https://legends2k.github.io/2d-fov/design.html
// http://ncase.me/sight-and-light/

func systemPerception(game *SoccerGame) {
	entitiesWithPerception := game.perceptorsView.Get()

	for _, entityResult := range entitiesWithPerception {
		perceptionAspect := entityResult.Components[game.perceptionComponent].(*Perception)
		perceptionAspect.SetPerception(computeAgentPerception(
			game,
			entityResult.Entity.GetID(),
		))
	}
}

func computeAgentPerception(game *SoccerGame, entityid ecs.EntityID) *agentPerception {
	//watch := utils.MakeStopwatch("computeAgentPerception()")
	//watch.Start("global")

	p := &agentPerception{}

	entityresult := game.getEntity(entityid,
		game.physicalBodyComponent,
		game.perceptionComponent,
	)

	if entityresult == nil {
		return p
	}

	physicalAspect := entityresult.Components[game.physicalBodyComponent].(*PhysicalBody)
	perceptionAspect := entityresult.Components[game.perceptionComponent].(*Perception)

	velocity := physicalAspect.GetVelocity()

	p.Velocity = velocity.Clone()
	p.Vision = computeAgentVision(game, entityresult.Entity, physicalAspect, perceptionAspect)

	return p
}

func computeAgentVision(game *SoccerGame, entity *ecs.Entity, physicalAspect *PhysicalBody, perceptionAspect *Perception) []agentPerceptionVisionItem {

	vision := make([]agentPerceptionVisionItem, 0)

	bodies := game.renderableView.Get()

	for _, entityResult := range bodies {
		physicalAspect := entityResult.Components[game.physicalBodyComponent].(*PhysicalBody)
		renderAspect := entityResult.Components[game.renderComponent].(*Render)
		bodyPosition := physicalAspect.GetPosition()
		bodyVelocity := physicalAspect.GetVelocity()
		bodyRadius := physicalAspect.GetRadius()

		vision = append(vision, agentPerceptionVisionItem{
			Center:   bodyPosition,
			EntityID: entityResult.Entity.GetID(),
			Radius:   bodyRadius,
			Tag:      renderAspect.GetType(),
			Velocity: bodyVelocity,
		})
	}

	return vision
}
