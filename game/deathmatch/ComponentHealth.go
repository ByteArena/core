package deathmatch

type Health struct {
	maxLife float64 // Const
	life    float64 // Current life level
}

func (health *Health) Restore() *Health {
	health.life = health.maxLife
	return health
}

func (health Health) GetMaxLife() float64 {
	return health.maxLife
}

func (health Health) GetLife() float64 {
	return health.life
}

func (health *Health) SetLife(life float64) {
	if life < 0 {
		life = 0
	}

	if life > health.maxLife {
		life = health.maxLife
	}

	health.life = life
}

func (health *Health) AddLife(life float64) {
	health.SetLife(life + health.GetLife())
}
