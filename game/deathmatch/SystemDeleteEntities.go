package deathmatch

import "github.com/bytearena/ecs"

func systemDeleteEntities(game *DeathmatchGame) {
	entitiesToRemove := make([]*ecs.Entity, 0)

	for _, entityResult := range game.lifecycleView.Get() {
		lifecycleAspect := entityResult.Components[game.lifecycleComponent].(*Lifecycle)
		if lifecycleAspect.delete {
			entitiesToRemove = append(entitiesToRemove, entityResult.Entity)
		}
	}

	if len(entitiesToRemove) > 0 {
		game.manager.DisposeEntities(entitiesToRemove...)
	}
}
