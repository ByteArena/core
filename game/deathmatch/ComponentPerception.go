package deathmatch

type Perception struct {
	visionAngle  float64 // expressed in rad
	visionRadius float64 // expressed in rad
	perception   *agentPerception
}

func (p Perception) GetVisionAngle() float64 {
	return p.visionAngle
}

func (p Perception) GetVisionRadius() float64 {
	return p.visionRadius
}

func (p *Perception) SetPerception(perception *agentPerception) *Perception {
	p.perception = perception
	return p
}

func (p Perception) GetPerception() *agentPerception {
	return p.perception
}
