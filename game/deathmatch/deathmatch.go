package deathmatch

import (
	"encoding/json"
	"log"
	"strconv"

	ebus "github.com/asaskevich/EventBus"
	"github.com/go-gl/mathgl/mgl64"

	"github.com/bytearena/box2d"
	"github.com/bytearena/ecs"

	"github.com/bytearena/core/common/types"
	commontypes "github.com/bytearena/core/common/types"
	"github.com/bytearena/core/common/utils"
	"github.com/bytearena/core/game/deathmatch/events"
	"github.com/bytearena/core/game/deathmatch/mailboxmessages"
)

type DeathmatchGame struct {
	ticknum int

	gameDescription commontypes.GameDescriptionInterface
	manager         *ecs.Manager

	bus ebus.Bus

	physicalToAgentSpaceTransform   *mgl64.Mat4
	physicalToAgentSpaceTranslation [3]float64
	physicalToAgentSpaceRotation    [3]float64
	physicalToAgentSpaceScale       float64

	physicalToAgentSpaceInverseTransform   *mgl64.Mat4
	physicalToAgentSpaceInverseTranslation [3]float64
	physicalToAgentSpaceInverseRotation    [3]float64
	physicalToAgentSpaceInverseScale       float64

	physicalBodyComponent *ecs.Component
	healthComponent       *ecs.Component
	playerComponent       *ecs.Component
	renderComponent       *ecs.Component
	scriptComponent       *ecs.Component
	perceptionComponent   *ecs.Component
	ownedComponent        *ecs.Component
	steeringComponent     *ecs.Component
	shootingComponent     *ecs.Component
	impactorComponent     *ecs.Component
	collidableComponent   *ecs.Component
	lifecycleComponent    *ecs.Component
	respawnComponent      *ecs.Component
	mailboxComponent      *ecs.Component
	sensorComponent       *ecs.Component

	agentsView      *ecs.View
	renderableView  *ecs.View
	physicalView    *ecs.View
	perceptorsView  *ecs.View
	shootingView    *ecs.View
	steeringView    *ecs.View
	impactorView    *ecs.View
	lifecycleView   *ecs.View
	respawnView     *ecs.View
	mailboxView     *ecs.View
	playerView      *ecs.View
	playerStatsView *ecs.View

	PhysicalWorld     *box2d.B2World
	collisionListener *collisionListener

	vizframe []byte

	variant     string
	cbkGameOver func()
}

func NewDeathmatchGame(gameDescription commontypes.GameDescriptionInterface) *DeathmatchGame {
	manager := ecs.NewManager()

	transform := mgl64.Ident4()
	inverseTransform := mgl64.Ident4()

	game := &DeathmatchGame{
		gameDescription: gameDescription,
		manager:         manager,

		// Variant: empty, or maze
		variant: gameDescription.GetMapContainer().Meta.Variant,

		bus: ebus.New(),

		physicalToAgentSpaceTransform:        &transform,
		physicalToAgentSpaceInverseTransform: &inverseTransform,

		physicalBodyComponent: manager.NewComponent(),
		healthComponent:       manager.NewComponent(),
		playerComponent:       manager.NewComponent(),
		renderComponent:       manager.NewComponent(),
		scriptComponent:       manager.NewComponent(),
		perceptionComponent:   manager.NewComponent(),
		ownedComponent:        manager.NewComponent(),
		steeringComponent:     manager.NewComponent(),
		shootingComponent:     manager.NewComponent(),
		impactorComponent:     manager.NewComponent(),
		collidableComponent:   manager.NewComponent(),
		lifecycleComponent:    manager.NewComponent(),
		respawnComponent:      manager.NewComponent(),
		mailboxComponent:      manager.NewComponent(),
		sensorComponent:       manager.NewComponent(),
	}

	game.setPhysicalToAgentSpaceTransform(
		100.0,               // scale
		[3]float64{0, 0, 0}, // translation
		[3]float64{0, 0, 0}, // rotation
	)

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

	game.shootingView = manager.CreateView(
		game.shootingComponent,
		game.physicalBodyComponent,
	)

	game.steeringView = manager.CreateView(
		game.steeringComponent,
		game.physicalBodyComponent,
		game.lifecycleComponent,
	)

	game.impactorView = manager.CreateView(
		game.impactorComponent,
		game.physicalBodyComponent,
	)

	game.lifecycleView = manager.CreateView(
		game.lifecycleComponent,
	)

	game.respawnView = manager.CreateView(
		game.respawnComponent,
	)

	game.mailboxView = manager.CreateView(
		game.mailboxComponent,
	)

	game.playerView = manager.CreateView(
		game.mailboxComponent,
		game.playerComponent,
	)

	game.playerStatsView = manager.CreateView(
		game.mailboxComponent,
		game.physicalBodyComponent,
		game.playerComponent,
	)

	game.physicalBodyComponent.SetDestructor(func(entity *ecs.Entity, data interface{}) {
		physicalAspect := data.(*PhysicalBody)
		game.PhysicalWorld.DestroyBody(physicalAspect.GetBody())
	})

	game.collisionListener = newCollisionListener(game)
	game.PhysicalWorld.SetContactListener(game.collisionListener)
	game.PhysicalWorld.SetContactFilter(newCollisionFilter(game))

	// Subscribing to gameplay events
	game.BusSubscribe(events.EntityFragged{}, func(e events.EntityFragged) {
		// We have to determine the identity of the fragger
		// Since it's a projectile that actually made the frag, not an agent

		fraggerEntityID := e.FraggedBy
		ownedQuery := game.getEntity(e.FraggedBy, game.ownedComponent)
		if ownedQuery != nil {
			ownedAspect := ownedQuery.Components[game.ownedComponent].(*Owned)
			fraggerEntityID = ownedAspect.GetOwner()
		}

		game.onEntityFraggedUpdateMailbox(e, fraggerEntityID)
		game.onEntityFraggedUpdateScore(e, fraggerEntityID)
	})

	game.BusSubscribe(events.EntityHit{}, game.onEntityHit)
	game.BusSubscribe(events.EntityRespawning{}, game.onEntityRespawning)
	game.BusSubscribe(events.EntityRespawned{}, game.onEntityRespawned)

	if game.variant == "maze" {
		game.BusSubscribe(events.EntityExitedMaze{}, game.onEntityExitedMaze)
	}

	return game
}

func (deathmatch *DeathmatchGame) Initialize(cbkGameOver func()) {
	deathmatch.cbkGameOver = cbkGameOver
}

func (deathmatch *DeathmatchGame) setPhysicalToAgentSpaceTransform(scale float64, translation, rotation [3]float64) *DeathmatchGame {

	deathmatch.physicalToAgentSpaceScale = scale
	deathmatch.physicalToAgentSpaceTranslation = translation
	deathmatch.physicalToAgentSpaceRotation = rotation

	rotxM := mgl64.HomogRotate3DX(mgl64.DegToRad(deathmatch.physicalToAgentSpaceRotation[0]))
	rotyM := mgl64.HomogRotate3DY(mgl64.DegToRad(deathmatch.physicalToAgentSpaceRotation[1]))
	rotzM := mgl64.HomogRotate3DZ(mgl64.DegToRad(deathmatch.physicalToAgentSpaceRotation[2]))
	transM := mgl64.Translate3D(deathmatch.physicalToAgentSpaceTranslation[0], deathmatch.physicalToAgentSpaceTranslation[1], deathmatch.physicalToAgentSpaceTranslation[2])
	scaleM := mgl64.Scale3D(deathmatch.physicalToAgentSpaceScale, deathmatch.physicalToAgentSpaceScale, deathmatch.physicalToAgentSpaceScale)

	transform := mgl64.Ident4().
		Mul4(transM).
		Mul4(rotzM).
		Mul4(rotyM).
		Mul4(rotxM).
		Mul4(scaleM)

	deathmatch.physicalToAgentSpaceTransform = &transform

	deathmatch.physicalToAgentSpaceInverseScale = 1.0 / scale
	deathmatch.physicalToAgentSpaceInverseTranslation = [3]float64{translation[0] * -1, translation[1] * -1, translation[2] * -1}
	deathmatch.physicalToAgentSpaceInverseRotation = [3]float64{rotation[0] * -1, rotation[1] * -1, rotation[2] * -1}

	inv := deathmatch.physicalToAgentSpaceTransform.Inv()
	deathmatch.physicalToAgentSpaceInverseTransform = &inv

	return deathmatch
}

func (deathmatch DeathmatchGame) getEntity(id ecs.EntityID, tagelements ...interface{}) *ecs.QueryResult {
	return deathmatch.manager.GetEntityByID(id, tagelements...)
}

// <GameInterface>

func (deathmatch *DeathmatchGame) ImplementsGameInterface() {}

func (deathmatch *DeathmatchGame) Step(ticknum int, dt float64, mutations []types.AgentMutationBatch) {

	//watch := utils.MakeStopwatch("deathmatch::Step()")
	//watch.Start("Step")

	deathmatch.ticknum = ticknum
	respawnersTag := ecs.BuildTag(deathmatch.respawnComponent)

	///////////////////////////////////////////////////////////////////////////
	// On fait mourir les non respawners début du tour (donc après le tour
	// précédent et la construction du message de visualisation du tour précédent).
	// Cela permet de conserver la vision des projectiles à l'endroit de leur disparition pendant 1 tick
	// Pour une meilleure précision de la position de collision dans la visualisation
	///////////////////////////////////////////////////////////////////////////

	//watch.Start("systemDeath")
	systemDeath(deathmatch, respawnersTag.Inverse())
	//watch.Stop("systemDeath")

	///////////////////////////////////////////////////////////////////////////
	// On traite les mutations
	///////////////////////////////////////////////////////////////////////////
	//watch.Start("systemMutations")
	systemMutations(deathmatch, mutations)
	//watch.Stop("systemMutations")

	///////////////////////////////////////////////////////////////////////////
	// On traite les tirs
	///////////////////////////////////////////////////////////////////////////
	//watch.Start("systemShooting")
	systemShooting(deathmatch)
	//watch.Stop("systemShooting")

	///////////////////////////////////////////////////////////////////////////
	// On traite les déplacements
	///////////////////////////////////////////////////////////////////////////
	//watch.Start("systemSteering")
	systemSteering(deathmatch)
	//watch.Stop("systemSteering")

	///////////////////////////////////////////////////////////////////////////
	// On met l'état des objets physiques à jour
	///////////////////////////////////////////////////////////////////////////
	//watch.Start("systemPhysics")
	systemPhysics(deathmatch, dt)
	//watch.Stop("systemPhysics")

	///////////////////////////////////////////////////////////////////////////
	// On identifie les collisions
	///////////////////////////////////////////////////////////////////////////
	//watch.Start("systemCollisions")
	collisions := systemCollisions(deathmatch)
	//watch.Stop("systemCollisions")

	///////////////////////////////////////////////////////////////////////////
	// On réagit aux collisions
	///////////////////////////////////////////////////////////////////////////
	//watch.Start("systemHealth")
	systemHealth(deathmatch, collisions)
	//watch.Stop("systemHealth")

	///////////////////////////////////////////////////////////////////////////
	// On fait vivre les entités
	///////////////////////////////////////////////////////////////////////////
	//watch.Start("systemLifecycle")
	systemLifecycle(deathmatch)
	//watch.Stop("systemLifecycle")

	///////////////////////////////////////////////////////////////////////////
	// On fait mourir les respawners tués au cours du tour
	///////////////////////////////////////////////////////////////////////////
	//watch.Start("systemDeath")
	systemDeath(deathmatch, respawnersTag)
	//watch.Stop("systemDeath")

	///////////////////////////////////////////////////////////////////////////
	// On ressuscite les entités qui peuvent l'être
	///////////////////////////////////////////////////////////////////////////
	//watch.Start("systemRespawn")
	systemRespawn(deathmatch)
	//watch.Stop("systemRespawn")

	///////////////////////////////////////////////////////////////////////////
	// On calcule les stats des agents
	///////////////////////////////////////////////////////////////////////////
	systemPlayerStats(deathmatch)

	///////////////////////////////////////////////////////////////////////////
	// On calcule le score des agents
	///////////////////////////////////////////////////////////////////////////
	systemScore(deathmatch)

	///////////////////////////////////////////////////////////////////////////
	// Fetching and emptying mailboxes for entities mailed in this tick
	///////////////////////////////////////////////////////////////////////////
	mailboxes := systemMailboxes(deathmatch)

	///////////////////////////////////////////////////////////////////////////
	// On construit les perceptions
	///////////////////////////////////////////////////////////////////////////
	//watch.Start("systemPerception")
	systemPerception(deathmatch, mailboxes)
	//watch.Stop("systemPerception")

	///////////////////////////////////////////////////////////////////////////
	// On supprime les entités marquées comme à supprimer
	// à la fin du tour pour éviter que box2D ne nile pas les références lors du disposeEntities
	///////////////////////////////////////////////////////////////////////////
	//watch.Start("systemDeleteEntities")
	systemDeleteEntities(deathmatch)
	//watch.Stop("systemDeleteEntities")

	//watch.Stop("Step")
	//fmt.Println(watch.String())

	deathmatch.ComputeVizFrame(mailboxes)
}

func (deathmatch *DeathmatchGame) GetAgentPerception(entityid ecs.EntityID) []byte {
	entityResult := deathmatch.getEntity(entityid, deathmatch.perceptionComponent)

	if entityResult == nil {
		return []byte{}
	}

	if perceptionAspect, ok := entityResult.Components[deathmatch.perceptionComponent].(*Perception); ok {
		bytes, _ := perceptionAspect.GetPerception().MarshalJSON()
		return bytes
	}

	return []byte{}
}

func (deathmatch *DeathmatchGame) GetAgentWelcome(entityid ecs.EntityID) []byte {

	entityresult := deathmatch.getEntity(entityid,
		deathmatch.physicalBodyComponent,
		deathmatch.steeringComponent,
		deathmatch.shootingComponent,
		deathmatch.perceptionComponent,
	)

	if entityresult == nil {
		return []byte{}
	}

	physicalAspect := entityresult.Components[deathmatch.physicalBodyComponent].(*PhysicalBody)
	steeringAspect := entityresult.Components[deathmatch.steeringComponent].(*Steering)
	shootingAspect := entityresult.Components[deathmatch.shootingComponent].(*Shooting)
	perceptionAspect := entityresult.Components[deathmatch.perceptionComponent].(*Perception)

	p := agentSpecs{
		// Movement
		MaxSpeed:           physicalAspect.GetMaxSpeed(),
		MaxAngularVelocity: physicalAspect.GetMaxAngularVelocity(),
		MaxSteeringForce:   steeringAspect.GetMaxSteeringForce(),
		VisionRadius:       perceptionAspect.GetVisionRadius(),
		VisionAngle:        commontypes.Angle(perceptionAspect.GetVisionAngle()),

		// Body
		BodyRadius: physicalAspect.GetRadius(),

		// Shoot
		MaxShootEnergy:    shootingAspect.MaxShootEnergy,
		ShootRecoveryRate: shootingAspect.ShootRecoveryRate,

		// DefaultWeapon: "gun",

		Gear: map[string]agentGearSpecs{
			"gun": agentGearSpecs{
				Genre: "weapon",
				Kind:  "gun",
				Specs: gunSpecs{
					ShootCost:        shootingAspect.ShootCost,
					ShootCooldown:    shootingAspect.ShootCooldown,
					ProjectileSpeed:  shootingAspect.ProjectileSpeed,
					ProjectileDamage: shootingAspect.ProjectileDamage,
					ProjectileRange:  shootingAspect.ProjectileRange,
				},
			},
		},
	}

	res, _ := p.MarshalJSON()
	return res
}

func (deathmatch *DeathmatchGame) GetVizInitJson() []byte {

	type vizMsgInit struct {
		//Map *mapcontainer.MapContainer `json:"map"`
		MapName string         `json:"mapname"`
		Tps     int            `json:"tps"`
		Agents  []*types.Agent `json:"agents"`
	}

	type VizMessage struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}

	initMsg := VizMessage{
		Type: "init",
		Data: vizMsgInit{
			MapName: deathmatch.gameDescription.GetName(),
			Tps:     deathmatch.gameDescription.GetTps(),
			Agents:  deathmatch.gameDescription.GetAgents(),
		},
	}

	res, _ := json.Marshal(initMsg)
	return res
}

func (deathmatch *DeathmatchGame) GetVizFrameJson() []byte {
	return deathmatch.vizframe
}

// </GameInterface>

func (deathmatch *DeathmatchGame) ComputeVizFrame(mailboxes map[ecs.EntityID]([]mailboxmessages.MailboxMessageInterface)) {

	msg := commontypes.VizMessage{
		GameID:        deathmatch.gameDescription.GetId(),
		Objects:       []commontypes.VizMessageObject{},
		DebugPoints:   make([][2]float64, 0),
		DebugSegments: make([][2][2]float64, 0),
		Events:        []commontypes.VizMessageEvent{},
	}

	for _, entityresult := range deathmatch.renderableView.Get() {

		renderAspect := entityresult.Components[deathmatch.renderComponent].(*Render)
		physicalBodyAspect := entityresult.Components[deathmatch.physicalBodyComponent].(*PhysicalBody)

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

		entityResultPlayer := deathmatch.getEntity(entityresult.Entity.ID, deathmatch.playerComponent)

		if entityResultPlayer != nil {
			playerAspect := entityResultPlayer.Components[deathmatch.playerComponent].(*Player)

			obj.PlayerInfo = &commontypes.PlayerInfo{
				PlayerName: playerAspect.Agent.Manifest.Name,
				PlayerId:   entityresult.Entity.GetID().String(),
				Score:      commontypes.VizMessagePlayerScore{playerAspect.Score},
				IsAlive:    true,
			}
		}

		entityResultRespawing := deathmatch.getEntity(entityresult.Entity.ID, deathmatch.respawnComponent)

		if entityResultRespawing != nil && obj.PlayerInfo != nil {
			respawnAspect := entityResultRespawing.Components[deathmatch.respawnComponent].(*Respawn)

			if respawnAspect.isRespawning {
				obj.PlayerInfo.IsAlive = false
			}
		}

		msg.Objects = append(msg.Objects, obj)

		// scaledDebugPoints := make([][2]float64, len(renderAspect.DebugPoints))
		// for i := 0; i < len(renderAspect.DebugPoints); i++ {
		// 	scaledDebugPoints[i] = vector.Vector2(renderAspect.DebugPoints[i]).
		// 		Transform(deathmatch.physicalToAgentSpaceInverseTransform).
		// 		ToFloatArray()
		// }
		// msg.DebugPoints = append(msg.DebugPoints, scaledDebugPoints...)

		// scaledDebugSegments := make([][2][2]float64, len(renderAspect.DebugSegments))
		// for i := 0; i < len(renderAspect.DebugSegments); i++ {
		// 	scaledDebugSegments[i] = [2][2]float64{
		// 		vector.Vector2(renderAspect.DebugSegments[i][0]).Transform(deathmatch.physicalToAgentSpaceInverseTransform).ToFloatArray(),
		// 		vector.Vector2(renderAspect.DebugSegments[i][1]).Transform(deathmatch.physicalToAgentSpaceInverseTransform).ToFloatArray(),
		// 	}
		// }
		// msg.DebugSegments = append(msg.DebugSegments, scaledDebugSegments...)
	}

	// Collecting events
	for entityid, mailbox := range mailboxes {

		for _, message := range mailbox {

			var payload interface{}
			subject := ""

			switch v := message.(type) {
			case mailboxmessages.YouHaveBeenFragged:
				subject = v.Subject()
				payload = map[string]string{
					"who": strconv.Itoa(int(entityid)),
					"by":  v.By,
				}
			case mailboxmessages.YouHaveRespawned:
				subject = v.Subject()
				payload = map[string]string{
					"who": strconv.Itoa(int(entityid)),
				}
			case mailboxmessages.YouHaveExitedTheMaze:
				subject = v.Subject()
				payload = map[string]string{
					"who": strconv.Itoa(int(entityid)),
				}
			}

			if subject != "" {
				msg.Events = append(msg.Events, commontypes.VizMessageEvent{
					Subject: subject,
					Payload: payload,
				})
			}
		}
	}

	deathmatch.vizframe, _ = msg.MarshalJSON()
}

func initPhysicalWorld(deathmatch *DeathmatchGame) {

	arenaMap := deathmatch.gameDescription.GetMapContainer()

	// Static obstacles formed by the grounds
	for _, ground := range arenaMap.Data.Grounds {
		deathmatch.NewEntityGround(ground.Polygon, ground.Name)
	}

	// Explicit obstacles
	for _, obstacle := range arenaMap.Data.Obstacles {
		polygon := obstacle.Polygon
		deathmatch.NewEntityObstacle(polygon, obstacle.Name)
	}

	if deathmatch.variant == "maze" {
		// Sensors

		for _, otherObject := range arenaMap.Data.OtherPolygonObjects {

			if utils.IsStringInArray(otherObject.Tags, "maze:exit") {
				polygon := otherObject.Polygon
				deathmatch.NewEntitySensor(
					polygon,
					otherObject.Name,
					func(entityid ecs.EntityID, sensorid ecs.EntityID) {
						deathmatch.BusPublish(events.EntityExitedMaze{
							Entity: entityid,
							Exit:   sensorid,
						})
					},
					utils.BuildTag(
						CollisionGroup.Agent,
						CollisionGroup.Projectile,
					),
				)
			}
		}
	}
}

func (deathmatch *DeathmatchGame) BusSubscribe(e events.EventInterface, cbk interface{}) {
	deathmatch.bus.Subscribe(e.Topic(), cbk)
}

func (deathmatch *DeathmatchGame) BusPublish(e events.EventInterface) {
	deathmatch.bus.Publish(e.Topic(), e)
}

func (game *DeathmatchGame) onEntityFraggedUpdateMailbox(e events.EntityFragged, fraggerEntityID ecs.EntityID) {

	///////////////////////////////////////////////////////////////////////
	// Notifying fraggee
	///////////////////////////////////////////////////////////////////////
	fraggeeMailboxQuery := game.getEntity(e.Entity, game.mailboxComponent)
	if fraggeeMailboxQuery == nil {
		return // should never happen
	}

	mailboxAspect := fraggeeMailboxQuery.Components[game.mailboxComponent].(*Mailbox)
	mailboxAspect.PushMessage(mailboxmessages.YouHaveBeenFragged{
		By: fraggerEntityID.String(),
	})

	///////////////////////////////////////////////////////////////////////
	// Notifying fragger
	///////////////////////////////////////////////////////////////////////

	fraggerMailboxQuery := game.getEntity(fraggerEntityID, game.mailboxComponent)
	if fraggerMailboxQuery == nil {
		// should never happen
		return
	}

	mailboxAspect = fraggerMailboxQuery.Components[game.mailboxComponent].(*Mailbox)
	mailboxAspect.PushMessage(mailboxmessages.YouHaveFragged{
		Who: e.Entity.String(),
	})
}

func (game *DeathmatchGame) onEntityFraggedUpdateScore(e events.EntityFragged, fraggerEntityID ecs.EntityID) {

	///////////////////////////////////////////////////////////////////////
	// Update score of the fraggee
	///////////////////////////////////////////////////////////////////////
	fraggeePlayerQuery := game.getEntity(e.Entity, game.playerComponent)
	if fraggeePlayerQuery == nil {
		return // should never happen
	}

	fraggeePlayerAspect := fraggeePlayerQuery.Components[game.playerComponent].(*Player)
	fraggeePlayerAspect.Stats.nbBeenFragged++

	///////////////////////////////////////////////////////////////////////
	// Update score of the fragger
	///////////////////////////////////////////////////////////////////////
	fraggerPlayerQuery := game.getEntity(fraggerEntityID, game.playerComponent)
	if fraggerPlayerQuery == nil {
		return // should never happen
	}

	fraggerPlayerAspect := fraggerPlayerQuery.Components[game.playerComponent].(*Player)
	fraggerPlayerAspect.Stats.nbHasFragged++
}

func (game *DeathmatchGame) onEntityHit(e events.EntityHit) {

	// We have to determine the identity of the hitter
	// Since it's a projectile that actually made the hit, not an agent

	hitterEntityID := e.HitBy
	ownedQuery := game.getEntity(e.HitBy, game.ownedComponent)
	if ownedQuery != nil {
		ownedAspect := ownedQuery.Components[game.ownedComponent].(*Owned)
		hitterEntityID = ownedAspect.GetOwner()
	}

	///////////////////////////////////////////////////////////////////////
	// Notifying hit entity
	///////////////////////////////////////////////////////////////////////

	query := game.getEntity(e.Entity, game.mailboxComponent)
	if query == nil {
		// should never happen
		return
	}

	mailboxAspect := query.Components[game.mailboxComponent].(*Mailbox)
	mailboxAspect.PushMessage(mailboxmessages.YouHaveBeenHit{
		Kind:       "projectile",
		ComingFrom: e.ComingFrom,
		Damage:     e.Damage,
	})

	///////////////////////////////////////////////////////////////////////
	// Notifying hitter entity
	///////////////////////////////////////////////////////////////////////

	query = game.getEntity(hitterEntityID, game.mailboxComponent)
	if query == nil {
		// should never happen
		return
	}

	mailboxAspect = query.Components[game.mailboxComponent].(*Mailbox)
	mailboxAspect.PushMessage(mailboxmessages.YouHaveHit{
		Who: string(e.Entity),
	})
}

func (game *DeathmatchGame) onEntityRespawning(e events.EntityRespawning) {
	query := game.getEntity(e.Entity, game.mailboxComponent)
	if query == nil {
		// should never happen
		return
	}

	mailboxAspect := query.Components[game.mailboxComponent].(*Mailbox)
	mailboxAspect.PushMessage(mailboxmessages.YouAreRespawning{
		RespawningIn: e.RespawnsIn,
	})
}

func (game *DeathmatchGame) onEntityRespawned(e events.EntityRespawned) {
	query := game.getEntity(e.Entity, game.mailboxComponent)
	if query == nil {
		// should never happen
		return
	}

	mailboxAspect := query.Components[game.mailboxComponent].(*Mailbox)
	mailboxAspect.PushMessage(mailboxmessages.YouHaveRespawned{})
}

func (game *DeathmatchGame) onEntityExitedMaze(e events.EntityExitedMaze) {
	query := game.getEntity(e.Entity, game.mailboxComponent)
	if query == nil {
		// should never happen
		return
	}

	mailboxAspect := query.Components[game.mailboxComponent].(*Mailbox)
	mailboxAspect.PushMessage(mailboxmessages.YouHaveExitedTheMaze{
		Entity: e.Entity,
	})

	log.Println("EXITED THE MAZE")
	game.cbkGameOver()
}
