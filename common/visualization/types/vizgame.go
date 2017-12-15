package types

import (
	"github.com/bytearena/core/common/types"
	"github.com/bytearena/core/common/utils"
)

type VizGame struct {
	gameDescription types.GameDescriptionInterface
	pool            *WatcherMap
}

func NewVizGame(gameDescription types.GameDescriptionInterface) *VizGame {
	return &VizGame{
		pool:            NewWatcherMap(),
		gameDescription: gameDescription,
	}
}

func (vizgame *VizGame) GetGame() types.GameDescriptionInterface {
	return vizgame.gameDescription
}

func (vizgame *VizGame) SetGame(game types.GameDescriptionInterface) {
	vizgame.gameDescription = game
}

func (vizgame *VizGame) GetTps() int {
	return vizgame.gameDescription.GetTps()
}

type vizMsgInit struct {
	//Map *mapcontainer.MapContainer `json:"map"`
	MapName string        `json:"mapname"`
	Tps     int           `json:"tps"`
	Agents  []types.Agent `json:"agents"`
}

type VizMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (vizgame *VizGame) SetWatcher(watcher *Watcher) {
	vizgame.pool.Set(watcher.GetId(), watcher)

	initMsg := VizMessage{
		Type: "init",
		Data: vizMsgInit{
			MapName: vizgame.gameDescription.GetName(),
			Tps:     vizgame.gameDescription.GetTps(),
			Agents:  vizgame.gameDescription.GetAgents(),
		},
	}

	err := watcher.conn.WriteJSON(initMsg)
	if err != nil {
		utils.Debug("viz-server", "Could not send VizInitMessage JSON;"+err.Error())
	}
}

func (vizgame *VizGame) RemoveWatcher(watcherid string) {
	vizgame.pool.Remove(watcherid)
}

func (vizgame *VizGame) GetNumberWatchers() int {
	return vizgame.pool.Size()
}
