package soccer

import (
	"encoding/json"
	"fmt"

	"github.com/bytearena/box2d"
	commontypes "github.com/bytearena/core/common/types"
	"github.com/bytearena/ecs"
	"github.com/go-gl/mathgl/mgl64"
)

type SoccerGame struct {
	ticknum     int
	cbkGameOver func()
	manager     *ecs.Manager

	perceptionComponent   *ecs.Component
	physicalBodyComponent *ecs.Component
	playerComponent       *ecs.Component
	steeringComponent     *ecs.Component
	collidableComponent   *ecs.Component
	renderComponent       *ecs.Component

	agentsView     *ecs.View
	renderableView *ecs.View
	physicalView   *ecs.View
	steeringView   *ecs.View
	playerView     *ecs.View
	perceptorsView *ecs.View

	PhysicalWorld *box2d.B2World

	agentSpaceToPhysicalSpaceScale float64

	physicalToAgentSpaceTransform   *mgl64.Mat4
	physicalToAgentSpaceTranslation [3]float64
	physicalToAgentSpaceRotation    [3]float64
	physicalToAgentSpaceScale       float64

	physicalToAgentSpaceInverseTransform   *mgl64.Mat4
	physicalToAgentSpaceInverseTranslation [3]float64
	physicalToAgentSpaceInverseRotation    [3]float64
	physicalToAgentSpaceInverseScale       float64

	TPS int

	collisionListener *collisionListener

	vizframe []byte
}

func NewSoccerGame() *SoccerGame {

	transform := mgl64.Ident4()
	inverseTransform := mgl64.Ident4()

	manager := ecs.NewManager()
	game := &SoccerGame{
		manager: manager,
		agentSpaceToPhysicalSpaceScale: 1,
		TPS: 20,
		physicalToAgentSpaceTransform:        &transform,
		physicalToAgentSpaceInverseTransform: &inverseTransform,

		physicalBodyComponent: manager.NewComponent(),

		playerComponent:     manager.NewComponent(),
		renderComponent:     manager.NewComponent(),
		perceptionComponent: manager.NewComponent(),
		steeringComponent:   manager.NewComponent(),
		collidableComponent: manager.NewComponent(),
	}

	gravity := box2d.MakeB2Vec2(0.0, 0.0) // gravity 0: the simulation is seen from the top
	world := box2d.MakeB2World(gravity)
	game.PhysicalWorld = &world

	initPhysicalWorld(game)

	game.physicalView = manager.CreateView(game.physicalBodyComponent)

	game.perceptorsView = manager.CreateView(game.perceptionComponent)

	game.agentsView = manager.CreateView(
		game.playerComponent,
		game.physicalBodyComponent,
	)

	game.renderableView = manager.CreateView(
		game.renderComponent,
		game.physicalBodyComponent,
	)

	game.steeringView = manager.CreateView(
		game.steeringComponent,
		game.physicalBodyComponent,
	)

	game.playerView = manager.CreateView(
		game.playerComponent,
	)

	game.physicalBodyComponent.SetDestructor(func(entity *ecs.Entity, data interface{}) {
		physicalAspect := data.(*PhysicalBody)
		game.PhysicalWorld.DestroyBody(physicalAspect.GetBody())
	})

	game.collisionListener = newCollisionListener(game)
	game.PhysicalWorld.SetContactListener(game.collisionListener)
	game.PhysicalWorld.SetContactFilter(newCollisionFilter(game))

	return game
}

func initPhysicalWorld(game *SoccerGame) {

	// Static physical objects
	// TODO: define outer perimeters (walls) and goals

	points := [][2]float64{
		[2]float64{0.0, 0.0},
		[2]float64{0.0, 117.0},
		[2]float64{117.0, 79.5},
		[2]float64{79.5, 0.0},
		[2]float64{0.0, 0.0},
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

	///////////////////////////////////////////////////////////////////////////
	// On traite les mutations
	///////////////////////////////////////////////////////////////////////////
	//watch.Start("systemMutations")
	systemMutations(game, mutations)
	//watch.Stop("systemMutations")

	///////////////////////////////////////////////////////////////////////////
	// On traite les déplacements
	///////////////////////////////////////////////////////////////////////////
	//watch.Start("systemSteering")
	systemSteering(game)
	//watch.Stop("systemSteering")

	///////////////////////////////////////////////////////////////////////////
	// On met l'état des objets physiques à jour
	///////////////////////////////////////////////////////////////////////////
	//watch.Start("systemPhysics")
	systemPhysics(game, dt)
	//watch.Stop("systemPhysics")

	///////////////////////////////////////////////////////////////////////////
	// On identifie les collisions
	///////////////////////////////////////////////////////////////////////////
	//watch.Start("systemCollisions")
	_ = systemCollisions(game)
	//watch.Stop("systemCollisions")

	// log.Println(collisions)

	///////////////////////////////////////////////////////////////////////////
	// On réagit aux collisions
	///////////////////////////////////////////////////////////////////////////

	// TODO

	///////////////////////////////////////////////////////////////////////////
	// On construit les perceptions
	///////////////////////////////////////////////////////////////////////////
	//watch.Start("systemPerception")
	systemPerception(game)
	//watch.Stop("systemPerception")

	game.ComputeVizFrame()
}

func (game *SoccerGame) GetAgentPerception(entityid ecs.EntityID) []byte {
	entityResult := game.getEntity(entityid, game.perceptionComponent)

	if entityResult == nil {
		return []byte(`{"case": "a"}`)
	}

	if perceptionAspect, ok := entityResult.Components[game.perceptionComponent].(*Perception); ok {
		bytes, _ := json.Marshal(perceptionAspect.GetPerception())
		return bytes
	}

	return []byte(`{"case": "b"}`)
}

func (game *SoccerGame) GetAgentWelcome(entityid ecs.EntityID) []byte {
	return []byte(
		`{ "id": "1", "team": "red" }`,
	)
}

func (game *SoccerGame) GetVizInitJson() []byte {
	return []byte(
		`{
			"type": "init",
			"field": {
				"width": 117,
				"height": 79.5,
				"padding": 6,
				"goallength": 7.3
			}
		}`,
	)
}

func (game *SoccerGame) GetVizFrameJson() []byte {
	return game.vizframe
}

func (game *SoccerGame) ComputeVizFrame() {

	msg := commontypes.VizMessage{
		GameID:  "1",
		Objects: []commontypes.VizMessageObject{},
	}

	for _, entityresult := range game.renderableView.Get() {

		renderAspect := entityresult.Components[game.renderComponent].(*Render)
		physicalBodyAspect := entityresult.Components[game.physicalBodyComponent].(*PhysicalBody)

		obj := commontypes.VizMessageObject{
			Id:   entityresult.Entity.GetID().String(),
			Type: renderAspect.GetType(),

			// Here, viz coord space and physical world coord space match
			// No transform is therefore needed
			Position:    physicalBodyAspect.GetPhysicalReferentialPosition(),
			Velocity:    physicalBodyAspect.GetPhysicalReferentialVelocity(),
			Radius:      physicalBodyAspect.GetPhysicalReferentialRadius(),
			Orientation: physicalBodyAspect.GetPhysicalReferentialOrientation(),

			PlayerInfo: nil,
		}

		entityResultPlayer := game.getEntity(entityresult.Entity.ID, game.playerComponent)

		if entityResultPlayer != nil {
			playerAspect := entityResultPlayer.Components[game.playerComponent].(*Player)

			obj.PlayerInfo = &commontypes.PlayerInfo{
				PlayerName: playerAspect.Agent.Manifest.Name,
				PlayerId:   entityresult.Entity.GetID().String(),
			}
		}

		msg.Objects = append(msg.Objects, obj)
	}

	game.vizframe, _ = msg.MarshalJSON()
}
