package types

import (
	"github.com/bytearena/core/common/types/mapcontainer"
)

type GameDescriptionInterface interface {
	GetId() string
	GetName() string
	GetTps() int
	GetRunStatus() int
	GetLaunchedAt() string
	GetEndedAt() string
	GetContestants() []Contestant
	GetMapContainer() *mapcontainer.MapContainer
}
