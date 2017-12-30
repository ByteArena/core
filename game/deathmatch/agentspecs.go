package deathmatch

import "github.com/bytearena/core/common/types"

type agentSpecs struct {
	// Movements
	MaxSpeed           float64     `json:"maxspeed"`         // max distance covered per turn
	MaxSteeringForce   float64     `json:"maxsteeringforce"` // max force applied when steering (max length from tip of current velocity vector to tip of next velocity vector)
	MaxAngularVelocity float64     `json:"maxangularvelocity"`
	VisionRadius       float64     `json:"visionradius"`
	VisionAngle        types.Angle `json:"visionangle"`

	// Body
	BodyRadius float64 `json:"bodyradius"`

	// Shoot
	MaxShootEnergy    float64 `json:"maxshootenergy"`
	ShootRecoveryRate float64 `json:"shootrecoveryrate"`

	Gear map[string]agentGearSpecs `json:"gear"`
}

type agentGearSpecs struct {
	Genre string      `json:"genre"` // Gun
	Kind  string      `json:"kind"`
	Specs interface{} `json:"specs"`
}

type gunSpecs struct {
	ShootCost        float64 `json:"shootcost"`        // energy cost of 1 projectile
	ShootCooldown    int     `json:"shootcooldown"`    // time to wait between shots (in ticks)
	ProjectileSpeed  float64 `json:"projectilespeed"`  // projectile speed (in m/tick)
	ProjectileDamage float64 `json:"projectiledamage"` // damage inflicted when projectile hits
	ProjectileRange  float64 `json:"projectilerange"`  // range of projectile, in m
}
