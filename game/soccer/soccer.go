package soccer

import (
	"github.com/bytearena/box2d"
	"github.com/bytearena/core/common/types"
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
}

func (game SoccerGame) getEntity(id ecs.EntityID, tagelements ...interface{}) *ecs.QueryResult {
	return game.manager.GetEntityByID(id, tagelements...)
}

func (game *SoccerGame) ImplementsGameInterface() {

}

func (game *SoccerGame) Initialize(cbkGameOver func()) {
	game.cbkGameOver = cbkGameOver
}

func (game *SoccerGame) Step(ticknum int, dt float64, mutations []types.AgentMutationBatch) {
	game.ticknum = ticknum
}

func (game *SoccerGame) NewEntityAgent(
	agent *types.Agent,
	spawnPosition vector.Vector2, // spawnPosition in physical space; TODO: fix this, should be in agent space
) ecs.EntityID {
	agentEntity := game.manager.NewEntity()
	return agentEntity.GetID()
}

func (game *SoccerGame) RemoveEntityAgent(agent *types.Agent) {
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
	return []byte("{}")
}

func (game *SoccerGame) GetVizFrameJson() []byte {
	return []byte("{}")
}
