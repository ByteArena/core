package deathmatch

import (
	"math"
	"strconv"
	"sync"

	"github.com/bytearena/box2d"
	"github.com/bytearena/ecs"

	commontypes "github.com/bytearena/core/common/types"
	"github.com/bytearena/core/common/utils/trigo"
	"github.com/bytearena/core/common/utils/vector"
	"github.com/bytearena/core/common/visibility2d"
	"github.com/bytearena/core/game/deathmatch/mailboxmessages"

	"github.com/bytearena/core/common/types/mapcontainer"
)

var pi2 = math.Pi * 2
var halfpi = math.Pi / 2
var threepi2 = math.Pi + halfpi

// https://legends2k.github.io/2d-fov/design.html
// http://ncase.me/sight-and-light/

func systemPerception(deathmatch *DeathmatchGame, mailboxes map[ecs.EntityID]([]mailboxmessages.MailboxMessageInterface)) {
	entitiesWithPerception := deathmatch.perceptorsView.Get()
	wg := sync.WaitGroup{}
	wg.Add(len(entitiesWithPerception))

	for _, entityResult := range entitiesWithPerception {
		perceptionAspect := entityResult.Components[deathmatch.perceptionComponent].(*Perception)
		go func(perceptionAspect *Perception, entity *ecs.Entity, wg *sync.WaitGroup) {

			entityID := entity.GetID()

			messages, ok := mailboxes[entityID]
			if !ok {
				messages = nil
			}

			perceptionAspect.SetPerception(computeAgentPerception(
				deathmatch,
				deathmatch.gameDescription.GetMapContainer(),
				entity.GetID(),
				messages,
			))
			wg.Done()
		}(perceptionAspect, entityResult.Entity, &wg)
	}

	wg.Wait()
}

func computeAgentPerception(game *DeathmatchGame, arenaMap *mapcontainer.MapContainer, entityid ecs.EntityID, messages []mailboxmessages.MailboxMessageInterface) *agentPerception {
	//watch := utils.MakeStopwatch("computeAgentPerception()")
	//watch.Start("global")

	p := &agentPerception{}

	entityresult := game.getEntity(entityid,
		game.physicalBodyComponent,
		game.steeringComponent,
		game.perceptionComponent,
		game.playerComponent,
	)

	if entityresult == nil {
		return p
	}

	physicalAspect := entityresult.Components[game.physicalBodyComponent].(*PhysicalBody)
	perceptionAspect := entityresult.Components[game.perceptionComponent].(*Perception)
	playerAspect := entityresult.Components[game.playerComponent].(*Player)

	orientation := physicalAspect.GetOrientation()
	velocity := physicalAspect.GetVelocity()

	p.Velocity = velocity.Clone().SetAngle(velocity.Angle() - orientation)
	p.Azimuth = orientation // l'angle d'orientation de l'agent par rapport au "Nord" de l'arène

	//watch.Start("p.External.Vision =")
	p.Vision = computeAgentVision(game, entityresult.Entity, physicalAspect, perceptionAspect)
	//watch.Stop("p.External.Vision =")

	// watch.Stop("global")
	// fmt.Println(watch.String())

	p.Messages = make([]mailboxMessagePerceptionWrapper, 0)

	if messages != nil {
		for i := 0; i < len(messages); i++ {
			msg := messages[i]
			p.Messages = append(p.Messages, mailboxMessagePerceptionWrapper{
				Subject: msg.Subject(),
				Body:    msg,
			})
		}
	}

	p.Score = playerAspect.Score

	return p
}

func computeAgentVision(game *DeathmatchGame, entity *ecs.Entity, physicalAspect *PhysicalBody, perceptionAspect *Perception) []agentPerceptionVisionItem {

	//watch := utils.MakeStopwatch("viewEntities()")
	//watch.Start("global")

	vision := make([]agentPerceptionVisionItem, 0)

	// for _, entityresult := range game.physicalView.Get() {
	// 	physicalAspect := entityresult.Components[game.physicalBodyComponent].(*PhysicalBody)
	// 	if physicalAspect.GetVelocity().Mag() > 0.01 {
	// 		physicalAspect.SetOrientation(physicalAspect.GetVelocity().Angle())
	// 	}
	// }

	agentPosition := physicalAspect.GetPosition()
	agentOrientation := physicalAspect.GetOrientation()
	visionAngle := perceptionAspect.GetVisionAngle()
	visionRadius := perceptionAspect.GetVisionRadius()
	visionRadiusSq := visionRadius * visionRadius

	halfVisionAngle := visionAngle / 2
	leftVisionEdgeAngle := math.Mod(agentOrientation-halfVisionAngle, pi2)
	rightVisionEdgeAngle := math.Mod(agentOrientation+halfVisionAngle, pi2)
	leftVisionRelvec := vector.MakeVector2(1, 1).SetMag(visionRadius).SetAngle(leftVisionEdgeAngle)
	rightVisionRelvec := vector.MakeVector2(1, 1).SetMag(visionRadius).SetAngle(rightVisionEdgeAngle)

	// Determine View cone AABB

	notableVisionConePoints := make([]vector.Vector2, 0)
	notableVisionConePoints = append(notableVisionConePoints, agentPosition)                        // center
	notableVisionConePoints = append(notableVisionConePoints, leftVisionRelvec.Add(agentPosition))  // left radius
	notableVisionConePoints = append(notableVisionConePoints, rightVisionRelvec.Add(agentPosition)) // right radius

	minAngle := math.Min(leftVisionEdgeAngle, rightVisionEdgeAngle)
	maxAngle := math.Max(leftVisionEdgeAngle, rightVisionEdgeAngle)

	if minAngle <= 0 && maxAngle > 0 {
		// Determine north point on circle
		notableVisionConePoints = append(notableVisionConePoints,
			vector.MakeVector2(1, 1).SetMag(visionRadius).SetAngle(0).Add(agentPosition),
		)
	}

	if minAngle <= halfpi && maxAngle > halfpi {
		// Determine east point on circle
		notableVisionConePoints = append(notableVisionConePoints,
			vector.MakeVector2(1, 1).SetMag(visionRadius).SetAngle(halfpi).Add(agentPosition),
		)
	}

	if minAngle <= math.Pi && maxAngle > math.Pi {
		// Determine south point on circle
		notableVisionConePoints = append(notableVisionConePoints,
			vector.MakeVector2(1, 1).SetMag(visionRadius).SetAngle(math.Pi).Add(agentPosition),
		)
	}

	if minAngle <= (threepi2) && maxAngle > (threepi2) {
		// Determine west point on circle
		notableVisionConePoints = append(notableVisionConePoints,
			vector.MakeVector2(1, 1).SetMag(visionRadius).SetAngle(threepi2).Add(agentPosition),
		)
	}

	entityAABB := vector.GetAABBForPointList(notableVisionConePoints...)
	elementsInAABB := make(map[ecs.EntityID]commontypes.PhysicalBodyDescriptor)

	game.PhysicalWorld.QueryAABB(func(fixture *box2d.B2Fixture) bool {
		if descriptor, ok := fixture.GetBody().GetUserData().(commontypes.PhysicalBodyDescriptor); ok {
			//elementsInAABB = append(elementsInAABB, descriptor)
			if _, isInMap := elementsInAABB[descriptor.ID]; !isInMap {
				elementsInAABB[descriptor.ID] = descriptor
			}
		}
		return true // keep going to find all fixtures in the query area
	}, entityAABB.Transform(game.physicalToAgentSpaceInverseTransform).ToB2AABB())

	//log.Println("AABB:", len(elementsInAABB))

	for _, bodyDescriptor := range elementsInAABB {

		if bodyDescriptor.ID == entity.ID {
			// one does not see itself
			continue
		}

		if bodyDescriptor.Type == commontypes.PhysicalBodyDescriptorType.Agent || bodyDescriptor.Type == commontypes.PhysicalBodyDescriptorType.Projectile {

			visionType := agentPerceptionVisionItemTag.Obstacle
			switch bodyDescriptor.Type {

			case commontypes.PhysicalBodyDescriptorType.Agent:
				visionType = agentPerceptionVisionItemTag.Agent

			case commontypes.PhysicalBodyDescriptorType.Projectile:
				visionType = agentPerceptionVisionItemTag.Projectile

			case commontypes.PhysicalBodyDescriptorType.Obstacle:
				visionType = agentPerceptionVisionItemTag.Obstacle
			case commontypes.PhysicalBodyDescriptorType.Ground:
				visionType = agentPerceptionVisionItemTag.Obstacle

			default:
				continue
			}

			//log.Println("Circle", bodyDescriptor.Type)
			// view a circle

			if bodyDescriptor.Type == commontypes.PhysicalBodyDescriptorType.Projectile {
				ownedQr := game.getEntity(bodyDescriptor.ID, game.ownedComponent)
				if ownedQr != nil {
					ownedAspect := ownedQr.Components[game.ownedComponent].(*Owned)
					if ownedAspect.GetOwner() == entity.GetID() {
						// do not show projectiles to their sender
						continue
					}
				}
			}

			otherQr := game.getEntity(bodyDescriptor.ID, game.physicalBodyComponent)
			otherPhysicalAspect := otherQr.Components[game.physicalBodyComponent].(*PhysicalBody)

			otherPosition := otherPhysicalAspect.GetPosition()
			otherVelocity := otherPhysicalAspect.GetVelocity()
			otherRadius := otherPhysicalAspect.GetRadius()

			if otherPosition.Equals(agentPosition) {
				// bodies have the exact same position; should never happen
				continue
			}

			centervec := otherPosition.Sub(agentPosition)
			centersegment := vector.MakeSegment2(vector.MakeNullVector2(), centervec)
			agentdiameter := centersegment.OrthogonalToBCentered().SetLengthFromCenter(otherRadius * 2)

			nearEdge, farEdge := agentdiameter.Get()

			distsq := centervec.MagSq()
			if distsq <= visionRadiusSq-0.5 { // -0.5: agent sight is sqrt(0.5)m shorter than obstacle sight, to avoid edge cases when agents and obstacles are really close and on the edge of the vision circle

				// Check that the obstacle is in our field of view
				centervecAngle := centervec.Angle()
				if centervecAngle < leftVisionRelvec.Angle() || centervecAngle > rightVisionRelvec.Angle() {
					continue
				}

				// On aligne l'angle du vecteur sur le heading courant de l'agent
				centervec = centervec.SetAngle(centervec.Angle() - agentOrientation)

				visionitem := agentPerceptionVisionItem{
					NearEdge:   nearEdge.Clone().SetAngle(nearEdge.Angle() - agentOrientation), // perpendicular to relative position vector, left side
					Center:     centervec,
					FarEdge:    farEdge.Clone().SetAngle(farEdge.Angle() - agentOrientation), // perpendicular to relative position vector, right side
					Velocity:   otherVelocity.Clone().SetAngle(otherVelocity.Angle() - agentOrientation),
					Tag:        visionType,
					EntityID:   bodyDescriptor.ID,
					SegmentNum: 0, // only one segment for circular bodies (diameter perpendicular to the viewer)
				}

				vision = append(vision, visionitem)
			}
		} else {

			// Obstacle
			// view a polygon

			otherQr := game.getEntity(bodyDescriptor.ID, game.physicalBodyComponent)
			otherPhysicalAspect := otherQr.Components[game.physicalBodyComponent].(*PhysicalBody)

			segmentNumber := -1
			fixture := otherPhysicalAspect.body.GetFixtureList()
			for fixture != nil {

				// Iterating over each segment of the polygon shape
				segmentNumber++ // starts at 0

				b2edge := fixture.GetShape().(*box2d.B2EdgeShape)
				fixture = fixture.M_next

				pointA := vector.FromB2Vec2(b2edge.M_vertex1).Transform(game.physicalToAgentSpaceTransform)
				pointB := vector.FromB2Vec2(b2edge.M_vertex2).Transform(game.physicalToAgentSpaceTransform)

				if !vector.GetAABBForPointList(pointA, pointB).Overlaps(entityAABB) {
					continue
				}

				edges := make([]vector.Vector2, 0)

				relvecA := pointA.Sub(agentPosition)
				relvecB := pointB.Sub(agentPosition)

				distsqA := relvecA.MagSq()
				distsqB := relvecB.MagSq()

				// Comment déterminer si le vecteur entre dans le champ de vision ?
				// => Intersection entre vecteur et segment gauche, droite

				if distsqA <= visionRadiusSq {
					// in radius
					absAngleA := relvecA.Angle()
					relAngleA := absAngleA - agentOrientation

					// On passe de 0° / 360° à -180° / +180°
					relAngleA = trigo.FullCircleAngleToSignedHalfCircleAngle(relAngleA)

					if math.Abs(relAngleA) <= halfVisionAngle {
						// point dans le champ de vision !
						edges = append(edges, relvecA.Add(agentPosition))
					} else {
						//rejectededges = append(rejectededges, relvecA.Add(absoluteposition))
					}
				}

				if distsqB <= visionRadiusSq {
					absAngleB := relvecB.Angle()
					relAngleB := absAngleB - agentOrientation

					// On passe de 0° / 360° à -180° / +180°
					relAngleB = trigo.FullCircleAngleToSignedHalfCircleAngle(relAngleB)

					if math.Abs(relAngleB) <= halfVisionAngle {
						// point dans le champ de vision !
						edges = append(edges, relvecB.Add(agentPosition))
					} else {
						//rejectededges = append(rejectededges, relvecB.Add(absoluteposition))
					}
				}

				{
					// Sur les bords de la perception
					if point, intersects, colinear, _ := trigo.IntersectionWithLineSegment(vector.MakeNullVector2(), leftVisionRelvec, relvecA, relvecB); intersects && !colinear {
						// INTERSECT LEFT
						edges = append(edges, point.Add(agentPosition))
					}

					if point, intersects, colinear, _ := trigo.IntersectionWithLineSegment(vector.MakeNullVector2(), rightVisionRelvec, relvecA, relvecB); intersects && !colinear {
						// INTERSECT RIGHT
						edges = append(edges, point.Add(agentPosition))
					}
				}

				{
					// Sur l'horizon de perception (arc de cercle)
					intersections := trigo.LineCircleIntersectionPoints(
						relvecA,
						relvecB,
						vector.MakeNullVector2(),
						visionRadius,
					)

					for _, point := range intersections {
						// il faut vérifier que le point se trouve bien sur le segment
						// il faut vérifier que l'angle du point de collision se trouve bien dans le champ de vision de l'agent

						if trigo.PointOnLineSegment(point, relvecA, relvecB) {
							relvecangle := point.Angle() - agentOrientation

							// On passe de 0° / 360° à -180° / +180°
							relvecangle = trigo.FullCircleAngleToSignedHalfCircleAngle(relvecangle)

							if math.Abs(relvecangle) <= halfVisionAngle {
								edges = append(edges, point.Add(agentPosition))
							} else {
								//rejectededges = append(rejectededges, point.Add(absoluteposition))
							}
						} else {
							//rejectededges = append(rejectededges, point.Add(absoluteposition))
						}
					}
				}

				if len(edges) == 2 {
					edgeone := edges[0]
					edgetwo := edges[1]
					center := edgetwo.Add(edgeone).DivScalar(2)

					//visiblemag := edgetwo.Sub(edgeone).Mag()

					relCenter := center.Sub(agentPosition) // aligned on north
					relCenterAngle := relCenter.Angle()
					relCenterAgentAligned := relCenter.SetAngle(relCenterAngle - agentOrientation)

					relEdgeOne := edgeone.Sub(agentPosition)
					relEdgeTwo := edgetwo.Sub(agentPosition)

					relEdgeOneAgentAligned := relEdgeOne.SetAngle(relEdgeOne.Angle() - agentOrientation)
					relEdgeTwoAgentAligned := relEdgeTwo.SetAngle(relEdgeTwo.Angle() - agentOrientation)

					var nearEdge, farEdge vector.Vector2
					if relEdgeTwoAgentAligned.MagSq() > relEdgeOneAgentAligned.MagSq() {
						nearEdge = relEdgeOneAgentAligned
						farEdge = relEdgeTwoAgentAligned
					} else {
						nearEdge = relEdgeTwoAgentAligned
						farEdge = relEdgeOneAgentAligned
					}

					obstacleperception := agentPerceptionVisionItem{
						NearEdge:   nearEdge,
						Center:     relCenterAgentAligned,
						FarEdge:    farEdge,
						Velocity:   vector.MakeNullVector2(),
						Tag:        agentPerceptionVisionItemTag.Obstacle,
						EntityID:   bodyDescriptor.ID,
						SegmentNum: segmentNumber,
					}

					vision = append(vision, obstacleperception)

				} else if len(edges) > 0 {
					// problems with FOV > 180
					//log.Println("SOMETHING'S WRONG !!!!!!!!!!!!!!!!!!!", len(edges))
				}
			}
		}
	}

	vision = processOcclusions(vision, agentPosition)

	return vision
}

type occlusionItem struct {
	visionItem     agentPerceptionVisionItem
	angleRealFrom  float64
	angleRealTo    float64
	angleRatioFrom float64
	angleRatioTo   float64
	distanceSq     float64
}

type byAngleRatio []occlusionItem

func (a byAngleRatio) Len() int           { return len(a) }
func (a byAngleRatio) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byAngleRatio) Less(i, j int) bool { return a[i].angleRatioFrom < a[j].angleRatioFrom }

func processOcclusions(vision []agentPerceptionVisionItem, agentPosition vector.Vector2) []agentPerceptionVisionItem {
	//return vision

	// Breaking segments at intersections

	breakableSegments := make([]visibility2d.ObstacleSegment, len(vision))
	for i := 0; i < len(vision); i++ {
		v := vision[i]
		breakableSegments[i] = visibility2d.ObstacleSegment{
			Points: [2][2]float64{
				v.NearEdge,
				v.FarEdge,
			},
			UserData: v,
		}
	}

	brokenSegments := visibility2d.OnlyVisible(
		agentPosition,
		breakableSegments,
	)

	// lenbefore := len(brokenSegments)

	/*

		A-----B C-----D

		* SI: A==B => destruction AB; reloop
		* SINON:
			* SI: C==D => destruction CD; reloop
			* SINON:
				* SI A==C => destruction AB; destruction CD; création BD; reloop
				* SINON:
					* SI: A==D => destruction AB; destruction CD; création BC; reloop
					* SINON:
						* SI B==C => destruction AB; destruction CD; création AD; reloop
						* SINON:
							* SI B==D => destruction AB; destruction CD; création AC; reloop

	*/

	// Sorting the broken segments by entity+segmentnum
	sortedSegments := map[string]([]*visibility2d.ObstacleSegment){}
	for i, _ := range brokenSegments {

		brokensegment := &brokenSegments[i]

		visionItem := brokensegment.UserData.(agentPerceptionVisionItem)
		segmentHash := strconv.Itoa(int(visionItem.EntityID)) + ":" + strconv.Itoa(visionItem.SegmentNum)

		var collection []*visibility2d.ObstacleSegment
		var found bool
		if collection, found = sortedSegments[segmentHash]; !found {
			collection = make([]*visibility2d.ObstacleSegment, 0)
		}

		collection = append(collection, brokensegment)
		sortedSegments[segmentHash] = collection
	}

	precision := 0.001

	finalSegments := make([]*visibility2d.ObstacleSegment, 0)

	for _, collection := range sortedSegments {
		mustIterate := true

		for mustIterate {

			mustIterate = false
			for i, _ := range collection {

				abIndex := i
				cdIndex := i + 1

				ab := collection[abIndex]
				var cd *visibility2d.ObstacleSegment

				if cdIndex < len(collection) {
					cd = collection[cdIndex]
				}

				a := vector.Vector2(ab.Points[0])
				b := vector.Vector2(ab.Points[1])

				if a.EqualsWithPrecision(b, precision) {
					// SI: A==B => destruction AB; reloop
					collection = append(collection[:abIndex], collection[abIndex+1:]...)
					mustIterate = true
					break
				} else {
					if cd == nil {
						break
					}

					c := vector.Vector2(cd.Points[0])
					d := vector.Vector2(cd.Points[1])

					if c.EqualsWithPrecision(d, precision) {
						// destruction CD; reloop
						collection = append(collection[:cdIndex], collection[cdIndex+1:]...)
						mustIterate = true
						break
					} else {

						if a.EqualsWithPrecision(c, precision) {
							// destruction AB; destruction CD; création BD; reloop
							collection = append(collection[:cdIndex], collection[cdIndex+1:]...)
							collection = append(collection[:abIndex], collection[abIndex+1:]...)
							collection = append(collection, &visibility2d.ObstacleSegment{
								Points: [2][2]float64{
									b,
									d,
								},
								UserData: ab.UserData,
							})

							mustIterate = true
							break
						} else {
							if a.EqualsWithPrecision(d, precision) {
								// destruction AB; destruction CD; création BC; reloop

								collection = append(collection[:cdIndex], collection[cdIndex+1:]...)
								collection = append(collection[:abIndex], collection[abIndex+1:]...)
								collection = append(collection, &visibility2d.ObstacleSegment{
									Points: [2][2]float64{
										b,
										c,
									},
									UserData: ab.UserData,
								})

								mustIterate = true
								break
							} else {
								if b.EqualsWithPrecision(c, precision) {
									// destruction AB; destruction CD; création AD; reloop

									collection = append(collection[:cdIndex], collection[cdIndex+1:]...)
									collection = append(collection[:abIndex], collection[abIndex+1:]...)
									collection = append(collection, &visibility2d.ObstacleSegment{
										Points: [2][2]float64{
											a,
											d,
										},
										UserData: ab.UserData,
									})

									mustIterate = true
									break
								} else {
									if b.EqualsWithPrecision(d, precision) {
										// destruction AB; destruction CD; création AC; reloop

										collection = append(collection[:cdIndex], collection[cdIndex+1:]...)
										collection = append(collection[:abIndex], collection[abIndex+1:]...)
										collection = append(collection, &visibility2d.ObstacleSegment{
											Points: [2][2]float64{
												a,
												c,
											},
											UserData: ab.UserData,
										})

										mustIterate = true
										break
									}
								}
							}
						}
					}
				}
			}
		}

		finalSegments = append(finalSegments, collection...)
	}

	//fmt.Println("FROM", lenbefore, "TO", len(finalSegments))

	realVision := make([]agentPerceptionVisionItem, len(finalSegments))

	//fmt.Println("--------------------------------------------------")
	for i, brokenSegment := range finalSegments {

		obs := vector.MakeSegment2(brokenSegment.Points[0], brokenSegment.Points[1])

		a := vector.Vector2(brokenSegment.Points[0])
		b := vector.Vector2(brokenSegment.Points[1])

		var nearEdge, farEdge vector.Vector2

		if a.MagSq() <= b.MagSq() {
			nearEdge, farEdge = a, b
		} else {
			nearEdge, farEdge = b, a
		}

		data := brokenSegment.UserData.(agentPerceptionVisionItem)

		realVision[i] = agentPerceptionVisionItem{
			Tag:        data.Tag,
			NearEdge:   nearEdge,
			FarEdge:    farEdge,
			Center:     obs.Center(),
			Velocity:   data.Velocity,
			EntityID:   data.EntityID,
			SegmentNum: data.SegmentNum,
		}

		//fmt.Println(realVision[i].EntityID, realVision[i].SegmentNum, realVision[i].NearEdge.String(), realVision[i].FarEdge.String())
	}

	return realVision

}

func getCircleSegmentAABB(center vector.Vector2, radius float64, angleARad float64, angleBRad float64) (lowerBound vector.Vector2, upperBound vector.Vector2) {
	return vector.MakeVector2(0, 0), vector.MakeVector2(0, 0)
}
