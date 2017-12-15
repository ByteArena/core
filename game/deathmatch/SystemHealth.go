package deathmatch

import (
	"github.com/bytearena/ecs"

	"github.com/bytearena/core/game/deathmatch/events"
)

type killedType struct {
	Entity   ecs.EntityID
	KilledBy ecs.EntityID
}

func systemHealth(deathmatch *DeathmatchGame, collisions []collision) {

	killed := make([]killedType, 0)

	for _, coll := range collisions {

		entityResultAImpactor := deathmatch.getEntity(coll.entityIDA, deathmatch.impactorComponent)
		entityResultAHealth := deathmatch.getEntity(coll.entityIDA, deathmatch.healthComponent)

		entityResultBImpactor := deathmatch.getEntity(coll.entityIDB, deathmatch.impactorComponent)
		entityResultBHealth := deathmatch.getEntity(coll.entityIDB, deathmatch.healthComponent)

		if entityResultAHealth != nil && entityResultBImpactor != nil {
			impactIfPossible(deathmatch, coll.entityIDA, entityResultAHealth, entityResultBImpactor, coll.collisionAngleA, &killed)
		}

		if entityResultBHealth != nil && entityResultAImpactor != nil {
			impactIfPossible(deathmatch, coll.entityIDB, entityResultBHealth, entityResultAImpactor, coll.collisionAngleB, &killed)
		}
	}

	for _, kill := range killed {
		lifecycleQr := deathmatch.getEntity(kill.Entity, deathmatch.lifecycleComponent)
		if lifecycleQr == nil {
			continue
		}

		lifecycleAspect := lifecycleQr.Components[deathmatch.lifecycleComponent].(*Lifecycle)
		lifecycleAspect.SetDeath(deathmatch.ticknum)

		// Publish Frag event
		deathmatch.BusPublish(events.EntityFragged{
			Entity:    kill.Entity,
			FraggedBy: kill.KilledBy,
		})
	}
}

func impactIfPossible(deathmatch *DeathmatchGame, impacteeID ecs.EntityID, impacteeHealth *ecs.QueryResult, impactor *ecs.QueryResult, collisionAngle float64, killed *[]killedType) {

	lifecycleQr := deathmatch.getEntity(impacteeID, deathmatch.lifecycleComponent)
	if lifecycleQr == nil {

		// no lifecycle on impactee; cannot be locked, impacting !
		impactWithDamage(deathmatch, impacteeHealth, impactor, collisionAngle, killed)
	} else {

		// There's a lifecycle on impactee; check if entity is locked
		lifecycleAspect := lifecycleQr.Components[deathmatch.lifecycleComponent].(*Lifecycle)

		if !lifecycleAspect.locked {
			// impactee not be locked, impacting !
			impactWithDamage(deathmatch, impacteeHealth, impactor, collisionAngle, killed)
		}
	}
}

func impactWithDamage(deathmatch *DeathmatchGame, qrHealth *ecs.QueryResult, qrImpactor *ecs.QueryResult, collisionAngle float64, killed *[]killedType) {

	impactedID := qrHealth.Entity.GetID()
	impactorID := qrImpactor.Entity.GetID()

	healthAspect := qrHealth.Components[deathmatch.healthComponent].(*Health)
	impactorAspect := qrImpactor.Components[deathmatch.impactorComponent].(*Impactor)

	// Publish Hit event
	deathmatch.BusPublish(events.EntityHit{
		Entity:     impactedID,
		HitBy:      impactorID,
		ComingFrom: collisionAngle,
		Damage:     impactorAspect.damage,
	})

	healthAspect.AddLife(-1 * impactorAspect.damage)
	if healthAspect.GetLife() <= 0 {
		healthAspect.SetLife(0)
		*killed = append(*killed, killedType{
			Entity:   impactedID,
			KilledBy: impactorID,
		})
	}
}
