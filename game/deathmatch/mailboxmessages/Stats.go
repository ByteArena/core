package mailboxmessages

type Stats struct {
	DistanceTravelled float64 `json:"distancetravelled"`
}

func (msg Stats) Subject() string {
	return "stats"
}
