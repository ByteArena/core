package mailboxmessages

type YouHaveRespawned struct{}

func (msg YouHaveRespawned) Subject() string {
	return "respawned"
}
