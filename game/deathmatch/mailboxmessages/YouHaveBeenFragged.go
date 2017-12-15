package mailboxmessages

type YouHaveBeenFragged struct {
	By string `json:"by"`
}

func (msg YouHaveBeenFragged) Subject() string {
	return "beenfragged"
}
