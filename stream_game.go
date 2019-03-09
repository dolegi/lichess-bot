package main

import (
	"encoding/json"
	"github.com/dolegi/uci"
	"log"
	"net/http"
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

func streamGame(gameId string, eng *uci.Engine) {
	req, err := http.NewRequest("GET", lichessUrl+"bot/game/stream/"+gameId, nil)
	req.Header.Add("Authorization", "Bearer "+apiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
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

		if gS.Type == "gameFull" && gS.White.Id == botName {
			goResp := eng.Go(uci.GoOpts{MoveTime: 100})

			makeMove(gameId, goResp.Bestmove)
		}
	}
}

func makeMove(gameId, move string) {
	url := lichessUrl + "bot/game/" + gameId + "/move/" + move
	log.Println(url)
	req, err := http.NewRequest("POST", lichessUrl+"bot/game/"+gameId+"/move/"+move, nil)
	req.Header.Add("Authorization", "Bearer "+apiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	log.Println(resp)
}
