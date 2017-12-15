package deathmatch

import (
	"sync"

	"github.com/bytearena/core/common/utils/vector"
)

type Steering struct {
	pendingSteers []vector.Vector2
	lock          *sync.RWMutex

	maxSteeringForce float64 // expressed in m/tick
}

func NewSteering(maxSteeringForce float64) *Steering {
	return &Steering{
		lock:             &sync.RWMutex{},
		maxSteeringForce: maxSteeringForce,
	}
}

func (steering Steering) GetMaxSteeringForce() float64 {
	return steering.maxSteeringForce
}

func (steering *Steering) SetMaxSteeringForce(maxSteeringForce float64) *Steering {
	steering.maxSteeringForce = maxSteeringForce
	return steering
}

func (steering *Steering) PushSteer(movement vector.Vector2) {
	steering.lock.Lock()
	steering.pendingSteers = append(steering.pendingSteers, movement)
	steering.lock.Unlock()
}

func (steering *Steering) PopPendingSteers() []vector.Vector2 {
	steering.lock.RLock()
	res := steering.pendingSteers
	steering.pendingSteers = make([]vector.Vector2, 0)
	steering.lock.RUnlock()

	return res
}
