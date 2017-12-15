package mailboxmessages

type Stats struct {
	DistanceTravelled float64 `json:"distanceTravelled"`
}

func (msg Stats) Subject() string {
	return "stats"
}
