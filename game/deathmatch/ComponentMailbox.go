package deathmatch

import (
	"github.com/bytearena/core/game/deathmatch/mailboxmessages"
)

type Mailbox struct {
	messages []mailboxmessages.MailboxMessageInterface
}

func (m *Mailbox) PushMessage(msg mailboxmessages.MailboxMessageInterface) *Mailbox {
	m.messages = append(m.messages, msg)
	return m
}

func (m *Mailbox) PopMessages() []mailboxmessages.MailboxMessageInterface {
	defer func() { m.messages = make([]mailboxmessages.MailboxMessageInterface, 0) }()
	return m.messages
}
