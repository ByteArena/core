package types

type AgentType struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	Title       string           `json:"title"`
	Image       *DockerImageType `json:"image"`
	Contestants []ContestantType `json:"contestants"`
}
