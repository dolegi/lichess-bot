package main

import (
	"encoding/json"
	"github.com/dolegi/uci"
	"log"
	"time"
)

type gameState struct {
	Type    string
	Id      string
	Rated   bool
	Variant struct {
		Key string
	}
	Clock struct {
		Initial   int
		Increment int
	}
	Speed      string
	InitialFen string
	State      struct {
		Type  string
		Moves string
		Wtime int
		Btime int
		Winc  int
		Binc  int
	}
	White struct {
		Id string
	}
	Black struct {
		Id string
	}
	Moves string
	Wtime int
	Btime int
	Winc  int
	Binc  int
}

func streamGame(gameId string) {
	resp := request("GET", "bot/game/stream/"+gameId)
	dec := json.NewDecoder(resp.Body)

	for dec.More() {
		var gS gameState
		err := dec.Decode(&gS)
		if err != nil {
			log.Println(err)
		}
		log.Printf("%v\n", gS)

		if gS.Type == "gameState" {
			eng.Position(gS.Moves)

			goResp := eng.Go(uci.GoOpts{
				Wtime: gS.Wtime,
				Btime: gS.Btime,
				Winc:  gS.Winc,
				Binc:  gS.Binc,
			})

			makeMove(gameId, goResp.Bestmove)
		}

		if gS.Type == "gameFull" && gS.White.Id == conf.Botname {
			goResp := eng.Go(uci.GoOpts{MoveTime: 100})

			makeMove(gameId, goResp.Bestmove)
		}
	}
}

func makeMove(gameId, move string) {
	resp := request("POST", "bot/game/"+gameId+"/move/"+move)
	log.Println(resp)
}
