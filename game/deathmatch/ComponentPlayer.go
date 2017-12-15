package deathmatch

import "github.com/bytearena/core/common/types"

type stats struct {

	// Distance travelled by the agent in meters since the beginning of the game
	distanceTravelled float64

	nbBeenFragged uint
	nbHasFragged  uint

	nbBeenHit uint
	nbHasHit  uint
}

type Player struct {
	// Populated by systemScore
	Score int

	Stats stats

	Agent types.Agent
}
