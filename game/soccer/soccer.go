package soccer

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"

	"github.com/bytearena/box2d"
	commontypes "github.com/bytearena/core/common/types"
	"github.com/bytearena/core/common/utils/vector"
	"github.com/bytearena/ecs"
)

type SoccerGame struct {
	ticknum     int
	cbkGameOver func()
	manager     *ecs.Manager

	physicalBodyComponent *ecs.Component

	PhysicalWorld *box2d.B2World
}

func NewSoccerGame() *SoccerGame {

	manager := ecs.NewManager()
	game := &SoccerGame{
		physicalBodyComponent: manager.NewComponent(),
		manager:               manager,
	}

	gravity := box2d.MakeB2Vec2(0.0, 0.0) // gravity 0: the simulation is seen from the top
	world := box2d.MakeB2World(gravity)
	game.PhysicalWorld = &world

	initPhysicalWorld(game)

	game.physicalBodyComponent.SetDestructor(func(entity *ecs.Entity, data interface{}) {
		// physicalAspect := data.(*PhysicalBody)
		// game.PhysicalWorld.DestroyBody(physicalAspect.GetBody())
	})

	return game
}

func initPhysicalWorld(game *SoccerGame) {

	// Static physical objects
	// TODO: define outer perimeters (walls) and goals

	points := [][2]float64{
		[2]float64{0.0, 0.0},
		[2]float64{0.0, 1050.0},
		[2]float64{1050.0, 680.0},
		[2]float64{0, 680.0},
		[2]float64{0, 0.0},
	}

	// north terrain boundary
	bodydef := box2d.MakeB2BodyDef()
	bodydef.Type = box2d.B2BodyType.B2_staticBody

	body := game.PhysicalWorld.CreateBody(&bodydef)
	vertices := make([]box2d.B2Vec2, len(points))

	for i := 0; i < len(points); i++ {
		vertices[i].Set(points[i][0], points[i][1]*-1) // TODO(jerome): invert axes in transform, not here
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("\n\nERROR - Terrain outer boundaries are not valid; perhaps some vertices are duplicated?\n\n")
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

	spew.Dump(body)

}

func (game SoccerGame) getEntity(id ecs.EntityID, tagelements ...interface{}) *ecs.QueryResult {
	return game.manager.GetEntityByID(id, tagelements...)
}

func (game *SoccerGame) ImplementsGameInterface() {

}

func (game *SoccerGame) Initialize(cbkGameOver func()) {
	game.cbkGameOver = cbkGameOver
}

func (game *SoccerGame) Step(ticknum int, dt float64, mutations []commontypes.AgentMutationBatch) {
	game.ticknum = ticknum
}

func (game *SoccerGame) NewEntityAgent(
	agent *commontypes.Agent,
	spawnPosition vector.Vector2, // spawnPosition in physical space; TODO: fix this, should be in agent space
) ecs.EntityID {
	agentEntity := game.manager.NewEntity()
	return agentEntity.GetID()
}

func (game *SoccerGame) RemoveEntityAgent(agent *commontypes.Agent) {
	qr := game.getEntity(agent.EntityID)
	game.manager.DisposeEntity(qr)
}

func (game *SoccerGame) GetAgentPerception(entityid ecs.EntityID) []byte {
	return []byte("{}")
}

func (game *SoccerGame) GetAgentWelcome(entityid ecs.EntityID) []byte {
	return []byte("{}")
}

func (game *SoccerGame) GetVizInitJson() []byte {
	return []byte(
		`{
			"field": {
				"width": 1050,
				"height": 680,
				"padding": 60,
				"goallength": 70
			}
		}`,
	)
}

func (game *SoccerGame) GetVizFrameJson() []byte {
	return []byte("{}")
}
