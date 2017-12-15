package deathmatch

type Impactor struct {
	damage float64
}

func (o Impactor) GetDamage() float64 {
	return o.damage
}
