package handler

import (
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/bytearena/core/common/mappack"
	commontypes "github.com/bytearena/core/common/types"
	"github.com/bytearena/core/common/visualization/types"
)

func Game(fetchVizGames func() ([]*types.VizGame, error), mappack *mappack.MappackInMemoryArchive) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		vizgames, err := fetchVizGames()
		if err != nil {
			w.Write([]byte("ERROR: Could not fetch viz games"))
			return
		}

		var gameDescription commontypes.GameDescriptionInterface
		foundgame := false

		for _, vizgame := range vizgames {
			if vizgame.GetGame().GetId() == vars["id"] {
				gameDescription = vizgame.GetGame()
				foundgame = true
				break
			}
		}

		if !foundgame {
			w.Write([]byte("GAME NOT FOUND !"))
			return
		}

		vizhtml, err := mappack.Open("index.html")
		if err != nil {
			w.Write([]byte("ERROR: could not render game"))
			return
		}

		protocol := "ws"

		if os.Getenv("ENV") == "prod" {
			protocol = "wss"
		}

		var vizhtmlTemplate = template.Must(template.New("").Parse(string(vizhtml)))
		err = vizhtmlTemplate.Execute(w, struct {
			WsURL      string
			CDNBaseURL string
			Rand       int64
			Tps        int
			Mappack    string
		}{
			WsURL:   protocol + "://" + r.Host + "/arena/" + gameDescription.GetId() + "/ws",
			Rand:    time.Now().Unix(),
			Tps:     gameDescription.GetTps(),
			Mappack: "/mappack/",
		})

		if err != nil {
			panic("err")
		}
	}
}
