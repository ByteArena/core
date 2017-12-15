package deathmatch

import (
	"sync"

	"github.com/bytearena/core/common/utils/vector"
)

type Shooting struct {
	pendingShots []vector.Vector2
	lock         *sync.RWMutex

	MaxShootEnergy    float64 // Const; When shooting, energy decreases
	ShootEnergy       float64 // Current energy level
	ShootRecoveryRate float64 // Const; Energy regained every tick
	ShootCost         float64 // Const; Energy consumed by a shot
	ShootCooldown     int     // Const; number of ticks to wait between every shot
	LastShot          int     // Number of ticks since last shot

	ProjectileSpeed  float64 // Const; expressed in m/tick
	ProjectileDamage float64 // Const; amount of life consumed on impact
	ProjectileRange  float64 // Const, in meter
}

func BuildShooting(shooting *Shooting) *Shooting {
	shooting.lock = &sync.RWMutex{}
	return shooting
}

func (shooting *Shooting) PushShot(aiming vector.Vector2) {
	shooting.lock.Lock()
	shooting.pendingShots = append(shooting.pendingShots, aiming)
	shooting.lock.Unlock()
}

func (shooting *Shooting) PopPendingShots() []vector.Vector2 {
	shooting.lock.RLock()
	res := shooting.pendingShots
	shooting.pendingShots = make([]vector.Vector2, 0)
	shooting.lock.RUnlock()

	return res
}
