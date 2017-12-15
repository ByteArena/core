package mailboxmessages

type Score struct {
	Value int `json:"value"`
}

func (msg Score) Subject() string {
	return "score"
}
