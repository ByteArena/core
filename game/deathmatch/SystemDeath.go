package deathmatch

import "github.com/bytearena/ecs"

func systemDeath(deathmatch *DeathmatchGame, filter ecs.Tag) {

	for _, entityresult := range deathmatch.lifecycleView.Get() {

		if !entityresult.Entity.Matches(filter) {
			continue
		}

		lifecycleAspect := entityresult.Components[deathmatch.lifecycleComponent].(*Lifecycle)
		if lifecycleAspect.tickDeath > 0 && !lifecycleAspect.deathProcessed {
			if lifecycleAspect.onDeath != nil {
				lifecycleAspect.onDeath()
			} else {
				//entitiesToRemove = append(entitiesToRemove, entityresult.Entity)
				lifecycleAspect.delete = true
			}
			lifecycleAspect.deathProcessed = true
		}
	}
}
