package types

import (
	"github.com/bytearena/core/common/types"
	"github.com/bytearena/core/common/utils"
	commongame "github.com/bytearena/core/game/common"
	"github.com/gorilla/websocket"
)

type VizGame struct {
	game            commongame.GameInterface
	gameDescription types.GameDescriptionInterface
	pool            *WatcherMap
}

func NewVizGame(game commongame.GameInterface, gameDescription types.GameDescriptionInterface) *VizGame {
	return &VizGame{
		pool:            NewWatcherMap(),
		game:            game,
		gameDescription: gameDescription,
	}
}

func (vizgame *VizGame) GetGameDescription() types.GameDescriptionInterface {
	return vizgame.gameDescription
}

func (vizgame *VizGame) SetGameDescription(game types.GameDescriptionInterface) {
	vizgame.gameDescription = game
}

func (vizgame *VizGame) GetTps() int {
	return vizgame.gameDescription.GetTps()
}

func (vizgame *VizGame) SetWatcher(watcher *Watcher) {
	vizgame.pool.Set(watcher.GetId(), watcher)

	initMsg := vizgame.game.GetVizInitJson()

	//err := watcher.conn.WriteJSON(initMsg)
	err := watcher.conn.WriteMessage(websocket.TextMessage, initMsg)
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
