package deathmatch

import "github.com/bytearena/core/common/utils/vector"

type agentPerception struct {
	Score int `json:"score"`

	Energy        float64                           `json:"energy"`   // niveau en millièmes; reconstitution automatique ?
	Velocity      vector.Vector2                    `json:"velocity"` // vecteur de force (direction, magnitude)
	Azimuth       float64                           `json:"azimuth"`  // azimuth en degrés par rapport au "Nord" de l'arène
	Vision        []agentPerceptionVisionItem       `json:"vision"`
	ShootEnergy   float64                           `json:"shootenergy"`
	ShootCooldown int                               `json:"shootcooldown"`
	Messages      []mailboxMessagePerceptionWrapper `json:"messages"`
}

var agentPerceptionVisionItemTag = struct {
	Agent      string
	Obstacle   string
	Projectile string
}{
	Agent:      "agent",
	Obstacle:   "obstacle",
	Projectile: "projectile",
}

type agentPerceptionVisionItem struct {
	Tag      string         `json:"tag"`
	NearEdge vector.Vector2 `json:"nearedge"`
	Center   vector.Vector2 `json:"center"`
	FarEdge  vector.Vector2 `json:"faredge"`
	Velocity vector.Vector2 `json:"velocity"`
}

type mailboxMessagePerceptionWrapper struct {
	Subject string      `json:"subject"`
	Body    interface{} `json:"body"`
}
