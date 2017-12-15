package types

type GameType struct {
	Id              string           `json:"id"`
	Tps             int              `json:"tps"`
	LaunchedAt      string           `json:"launchedAt"`
	EndedAt         string           `json:"endedAt"`
	Arena           *ArenaType       `json:"arena"`
	Contestants     []ContestantType `json:"contestants"`
	ArenaServerUUID string           `json:"arenaServerUUID"`
	RunStatus       int              `json:"runStatus"`
	RunError        string           `json:"runError"`
}

var GameRunStatus = struct {
	Pending  int
	Running  int
	Finished int
}{
	Pending:  0,
	Running:  1,
	Finished: 2,
}
