package deathmatch

import (
	"github.com/bytearena/core/game/deathmatch/mailboxmessages"
)

func systemScore(deathmatch *DeathmatchGame) {

	for _, result := range deathmatch.playerView.Get() {
		playerAspect := result.Components[deathmatch.playerComponent].(*Player)

		oldScore := playerAspect.Score

		playerAspect.Score = calculatePlayerScore(playerAspect)

		if playerAspect.Score != oldScore {
			mailboxAspect := result.Components[deathmatch.mailboxComponent].(*Mailbox)

			sendScoreToAgent(mailboxAspect, playerAspect)
		}

	}
}

func calculatePlayerScore(p *Player) (score int) {
	score += int(p.Stats.nbHasFragged)
	score -= int(p.Stats.nbBeenFragged)

	return score
}

func sendScoreToAgent(mailbox *Mailbox, player *Player) {
	mailbox.PushMessage(mailboxmessages.Score{
		Value: player.Score,
	})
}
