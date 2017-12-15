package handler

import (
	"fmt"
	"log"
	"net/http"

	notify "github.com/bitly/go-notify"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/bytearena/core/common/utils"
	"github.com/bytearena/core/common/visualization/types"
)

type wsincomingmessage struct {
	messageType int
	p           []byte
	err         error
}

func Websocket(fetchVizGames func() ([]*types.VizGame, error)) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		vizgames, err := fetchVizGames()
		if err != nil {
			w.Write([]byte("ERROR: Could not fetch viz games"))
			return
		}

		var vizgame *types.VizGame
		foundgame := false

		for _, vizgameit := range vizgames {
			if vizgameit.GetGame().GetId() == vars["id"] {
				vizgame = vizgameit
				foundgame = true
				break
			}
		}

		if !foundgame {
			w.Write([]byte("GAME NOT FOUND !"))
			return
		}

		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			EnableCompression: true,
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		c.EnableWriteCompression(true)

		watcher := types.NewWatcher(c)
		vizgame.SetWatcher(watcher)

		defer func(c *websocket.Conn) {
			vizgame.RemoveWatcher(watcher.GetId())
			c.Close()
		}(c)

		/////////////////////////////////////////////////////////////
		/////////////////////////////////////////////////////////////
		/////////////////////////////////////////////////////////////

		clientclosedsocket := make(chan bool)
		c.SetCloseHandler(func(code int, text string) error {
			clientclosedsocket <- true
			return nil
		})

		// Listen to messages incoming from viz; mandatory to notice when websocket is closed client side
		incomingmsg := make(chan wsincomingmessage)
		go func(client *websocket.Conn, ch chan wsincomingmessage) {
			messageType, p, err := client.ReadMessage()
			ch <- wsincomingmessage{messageType, p, err}
		}(c, incomingmsg)

		// Listen to viz messages coming from arenaserver
		vizmsgchan := make(chan interface{})

		notify.Start("viz:message:"+vizgame.GetGame().GetId(), vizmsgchan)

		for {
			select {
			case <-clientclosedsocket:
				{
					utils.Debug("ws", "disconnected")
					return
				}
			case vizmsg := <-vizmsgchan:
				{
					vizmsgString, ok := vizmsg.(string)
					utils.Assert(ok, "Failed to cast vizmessage into string")

					data := fmt.Sprintf("{\"type\":\"framebatch\", \"data\": %s}", vizmsgString)

					c.WriteMessage(websocket.TextMessage, []byte(data))
				}
			}
		}
	}
}
