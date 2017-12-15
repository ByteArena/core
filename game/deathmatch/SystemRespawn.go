package deathmatch

func systemRespawn(deathmatch *DeathmatchGame) {

	for _, entityresult := range deathmatch.respawnView.Get() {
		respawnAspect := entityresult.Components[deathmatch.respawnComponent].(*Respawn)
		if respawnAspect.isRespawning {
			respawnAspect.respawningCountdown--
			if respawnAspect.respawningCountdown <= 0 {
				respawnAspect.isRespawning = false
				respawnAspect.respawnCount++
				respawnAspect.onRespawn()
			}
		}
	}

}
