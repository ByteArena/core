package types

type ArenaType struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Kind           string `json:"kind"`
	MaxContestants int    `json:"maxContestants"`
}
