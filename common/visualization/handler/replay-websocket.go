package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/bytearena/core/common/replay"

	"github.com/bytearena/core/common/recording"
	"github.com/bytearena/core/common/utils"
)

func ReplayWebsocket(recordStore recording.RecordStoreInterface, basepath string) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UUID := vars["recordId"]

		if !recordStore.RecordExists(UUID) {
			w.Write([]byte("Record not found"))
			return
		}

		recordFile := recordStore.GetFilePathForUUID(UUID)

		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		clientclosedsocket := make(chan bool)

		defer func(c *websocket.Conn) {
			c.Close()
			clientclosedsocket <- true
		}(c)

		/////////////////////////////////////////////////////////////
		/////////////////////////////////////////////////////////////
		/////////////////////////////////////////////////////////////

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

		debug := false

		replayer := replay.NewReplayer(recordFile, debug, UUID)

		vizmapmsgchan := replayer.ReadMap()
		vizmsgchan := replayer.Read()

		for {
			select {
			case <-clientclosedsocket:
				{
					utils.Debug("ws", "disconnected")
					replayer.Stop()
					break
				}
			case vizmsg := <-vizmsgchan:
				{
					// End of the record
					if vizmsg == nil {
						return
					}

					data := fmt.Sprintf("{\"type\":\"framebatch\", \"data\": %s}", vizmsg.Line)

					c.WriteMessage(websocket.TextMessage, []byte(data))
					<-time.NewTimer(1 * time.Second).C
				}
			case vizmap := <-vizmapmsgchan:
				{
					initMessage := "{\"type\":\"init\",\"data\": " + vizmap + "}"
					c.WriteMessage(websocket.TextMessage, []byte(initMessage))
				}
			}
		}
	}
}
