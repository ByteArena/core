package types

// TODO: merge / disambiguify this with Contestant
type ContestantType struct {
	Id              string     `json:"id"`
	Agent           *AgentType `json:"agent"`
	Game            *GameType  `json:"game"`
	EnrolledAt      string     `json:"enrolledAt"`
	ContestantLogId string     `json:"contestantLogId"`
}
