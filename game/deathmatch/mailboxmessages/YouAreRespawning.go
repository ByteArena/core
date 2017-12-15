package mailboxmessages

type YouAreRespawning struct {
	RespawningIn int `json:"respawningin"`
}

func (msg YouAreRespawning) Subject() string {
	return "respawning"
}
