package deathmatch

import (
	"math"

	"github.com/bytearena/core/game/deathmatch/mailboxmessages"
)

func systemPlayerStats(deathmatch *DeathmatchGame) {
	for _, result := range deathmatch.playerStatsView.Get() {
		playerAspect := result.Components[deathmatch.playerComponent].(*Player)
		physicalAspect := result.Components[deathmatch.physicalBodyComponent].(*PhysicalBody)

		playerAspect.Stats.distanceTravelled += math.Abs(physicalAspect.GetVelocity().Mag())

		mailboxAspect := result.Components[deathmatch.mailboxComponent].(*Mailbox)

		sendStatsToAgent(mailboxAspect, playerAspect)
	}
}

func sendStatsToAgent(mailbox *Mailbox, player *Player) {

	mailbox.PushMessage(mailboxmessages.Stats{
		DistanceTravelled: player.Stats.distanceTravelled,
	})
}
