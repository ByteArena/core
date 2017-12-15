package deathmatch

import (
	"github.com/bytearena/core/common/utils/trigo"
)

func systemShooting(deathmatch *DeathmatchGame) {

	for _, entityresult := range deathmatch.shootingView.Get() {

		shootingAspect := entityresult.Components[deathmatch.shootingComponent].(*Shooting)
		physicalAspect := entityresult.Components[deathmatch.physicalBodyComponent].(*PhysicalBody)

		shootingAspect.ShootEnergy += shootingAspect.ShootRecoveryRate
		if shootingAspect.ShootEnergy > shootingAspect.MaxShootEnergy {
			shootingAspect.ShootEnergy = 0
		}

		shots := shootingAspect.PopPendingShots()
		if len(shots) == 0 {
			continue
		}

		// //
		// // Levels consumption
		// //

		if deathmatch.ticknum-shootingAspect.LastShot <= shootingAspect.ShootCooldown {
			// invalid shot, cooldown not over
			continue
		}

		if shootingAspect.ShootEnergy < shootingAspect.ShootCost {
			// invalid shot, not enough energy
			continue
		}

		aiming := shots[0]
		entity := entityresult.Entity
		if aiming.IsNull() {
			// 0-mag aiming vector disabled (no mines !)
			continue
		}

		shootingAspect.LastShot = deathmatch.ticknum
		shootingAspect.ShootEnergy -= shootingAspect.ShootCost

		///////////////////////////////////////////////////////////////////////////
		///////////////////////////////////////////////////////////////////////////
		// Make physical body for projectile
		///////////////////////////////////////////////////////////////////////////
		///////////////////////////////////////////////////////////////////////////

		orientation := physicalAspect.GetOrientation()

		// on passe le vecteur de visée d'un angle relatif à un angle absolu
		velocity := trigo.
			LocalAngleToAbsoluteAngleVec(orientation, aiming, nil). // TODO: replace nil here by an actual angle constraint
			SetMag(1)                                               // Unit vector for aiming

		deathmatch.NewEntityBallisticProjectile(entity.GetID(), physicalAspect.GetPosition(), velocity)
	}
}
