package deathmatch

import (
	"github.com/bytearena/box2d"

	commontypes "github.com/bytearena/core/common/types"
)

///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////
// Collision Handling
///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

type collisionFilter struct { /* implements box2d.B2World.B2ContactFilterInterface */
	game *DeathmatchGame
}

func (filter *collisionFilter) ShouldCollide(fixtureA *box2d.B2Fixture, fixtureB *box2d.B2Fixture) bool {

	descriptorA, ok := fixtureA.GetBody().GetUserData().(commontypes.PhysicalBodyDescriptor)
	if !ok {
		return false
	}

	descriptorB, ok := fixtureB.GetBody().GetUserData().(commontypes.PhysicalBodyDescriptor)
	if !ok {
		return false
	}

	game := filter.game

	entityResultA := game.getEntity(descriptorA.ID, game.collidableComponent)
	entityResultB := game.getEntity(descriptorB.ID, game.collidableComponent)

	if entityResultA == nil || entityResultB == nil {
		return false
	}

	collidableAspectA := entityResultA.Components[game.collidableComponent].(*Collidable)
	collidableAspectB := entityResultB.Components[game.collidableComponent].(*Collidable)

	mayGroupsCollide := collidableAspectA.MayCollideWith(collidableAspectB) || collidableAspectB.MayCollideWith(collidableAspectA)
	if !mayGroupsCollide {
		// groups cannot collide
		return false
	}

	// groups can collide; still have to check if there's an owner/owned relationship between the two (owned cannot collide with owner)
	// filtering here because unfiltered collisions do have an impact on Box2D bodies movements

	entityResultOwnedA := game.getEntity(descriptorA.ID, game.ownedComponent)
	entityResultOwnedB := game.getEntity(descriptorB.ID, game.ownedComponent)

	if entityResultOwnedA != nil {
		ownedAspect := entityResultOwnedA.Components[game.ownedComponent].(*Owned)
		if ownedAspect.GetOwner() == descriptorB.ID {
			return false
		}
	}

	if entityResultOwnedB != nil {
		ownedAspect := entityResultOwnedB.Components[game.ownedComponent].(*Owned)
		if ownedAspect.GetOwner() == descriptorA.ID {
			return false
		}
	}

	return true
}

func newCollisionFilter(game *DeathmatchGame) *collisionFilter {
	return &collisionFilter{
		game: game,
	}
}

type collisionListener struct { /* implements box2d.B2World.B2ContactListenerInterface */
	game            *DeathmatchGame
	collisionbuffer []box2d.B2ContactInterface
}

func (listener *collisionListener) PopCollisions() []box2d.B2ContactInterface {
	defer func() { listener.collisionbuffer = make([]box2d.B2ContactInterface, 0) }()
	return listener.collisionbuffer
}

/// Called when two fixtures begin to touch.
func (listener *collisionListener) BeginContact(contact box2d.B2ContactInterface) { // contact has to be backed by a pointer
	listener.collisionbuffer = append(listener.collisionbuffer, contact)
}

/// Called when two fixtures cease to touch.
func (listener *collisionListener) EndContact(contact box2d.B2ContactInterface) { // contact has to be backed by a pointer
	//log.Println("END:COLLISION !!!!!!!!!!!!!!")
}

/// This is called after a contact is updated. This allows you to inspect a
/// contact before it goes to the solver. If you are careful, you can modify the
/// contact manifold (e.g. disable contact).
/// A copy of the old manifold is provided so that you can detect changes.
/// Note: this is called only for awake bodies.
/// Note: this is called even when the number of contact points is zero.
/// Note: this is not called for sensors.
/// Note: if you set the number of contact points to zero, you will not
/// get an EndContact callback. However, you may get a BeginContact callback
/// the next step.
func (listener *collisionListener) PreSolve(contact box2d.B2ContactInterface, oldManifold box2d.B2Manifold) { // contact has to be backed by a pointer
	//log.Println("PRESOLVE !!!!!!!!!!!!!!")
}

/// This lets you inspect a contact after the solver is finished. This is useful
/// for inspecting impulses.
/// Note: the contact manifold does not include time of impact impulses, which can be
/// arbitrarily large if the sub-step is small. Hence the impulse is provided explicitly
/// in a separate data structure.
/// Note: this is only called for contacts that are touching, solid, and awake.
func (listener *collisionListener) PostSolve(contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) { // contact has to be backed by a pointer
	//log.Println("POSTSOLVE !!!!!!!!!!!!!!")
}

func newCollisionListener(game *DeathmatchGame) *collisionListener {
	return &collisionListener{
		game: game,
	}
}

///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////
