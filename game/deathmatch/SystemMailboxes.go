package deathmatch

import (
	"github.com/bytearena/ecs"

	"github.com/bytearena/core/game/deathmatch/mailboxmessages"
)

func systemMailboxes(deathmatch *DeathmatchGame) map[ecs.EntityID]([]mailboxmessages.MailboxMessageInterface) {

	mailboxes := make(map[ecs.EntityID]([]mailboxmessages.MailboxMessageInterface))

	for _, entityresult := range deathmatch.mailboxView.Get() {
		mailboxAspect := entityresult.Components[deathmatch.mailboxComponent].(*Mailbox)

		messages := mailboxAspect.PopMessages()
		if len(messages) > 0 {
			mailboxes[entityresult.Entity.GetID()] = messages
		}
	}

	return mailboxes
}
