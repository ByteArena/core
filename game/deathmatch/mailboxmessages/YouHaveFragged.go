package mailboxmessages

type YouHaveFragged struct {
	Who string `json:"who"`
}

func (msg YouHaveFragged) Subject() string {
	return "havefragged"
}
