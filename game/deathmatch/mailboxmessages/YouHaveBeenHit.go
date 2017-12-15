package mailboxmessages

type YouHaveBeenHit struct {
	Kind       string  `json:"kind"`
	ComingFrom float64 `json:"comingfrom"` // absolute azimuth
	Damage     float64 `json:"damage"`     // life consumed upon impact
}

func (msg YouHaveBeenHit) Subject() string {
	return "beenhit"
}
