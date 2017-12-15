package deathmatch

type Respawn struct {
	isRespawning        bool
	respawningCountdown int
	respawnCount        int
	onRespawn           func()
}
