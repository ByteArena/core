package mailboxmessages

type YouHaveHit struct {
	Who string
}

func (msg YouHaveHit) Subject() string {
	return "havehit"
}
