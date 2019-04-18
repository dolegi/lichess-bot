package main

import (
	"encoding/json"
	"github.com/dolegi/uci"
	"log"
	"strings"
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

			if white {
				gS.Wtime -= conf.Network.Latency
			} else {
				gS.Btime -= conf.Network.Latency
			}

			if (gS.Moves == "" || (strings.Count(gS.Moves, " ")%2 == 1)) != (white == whiteFirst) {
				continue
			}

			opts := uci.GoOpts{
				Wtime: gS.Wtime,
				Btime: gS.Btime,
				Winc:  gS.Winc,
				Binc:  gS.Binc,
			}
			if conf.Engine.Go.Nodes > 0 {
				opts.Nodes = conf.Engine.Go.Nodes
			}
			if conf.Engine.Go.Depth > 0 {
				opts.Depth = conf.Engine.Go.Depth
			}
			if conf.Engine.Go.Movetime > 0 {
				opts.MoveTime = conf.Engine.Go.Movetime
			}

			goResp := eng.Go(opts)

			makeMove(gameId, goResp.Bestmove)
		}

		if gS.Type == "gameFull" {
			eng.NewGame(uci.NewGameOpts{Variant: gS.Variant, InitialFen: gS.InitialFen, Moves: gS.State.Moves})

			white = (gS.White.Id == conf.Botname)
			whiteFirst = gS.InitialFen == "" || gS.InitialFen == "startpos" || strings.Contains(gS.InitialFen, "w")
			if (gS.State.Moves == "" || (strings.Count(gS.State.Moves, " ")%2 == 1)) != (white == whiteFirst) {
				continue
			}

			opts := uci.GoOpts{
				Wtime: gS.State.Wtime,
				Btime: gS.State.Btime,
				Winc:  gS.State.Winc,
				Binc:  gS.State.Binc,
			}
			if conf.Engine.Go.Nodes > 0 {
				opts.Nodes = conf.Engine.Go.Nodes
			}
			if conf.Engine.Go.Depth > 0 {
				opts.Depth = conf.Engine.Go.Depth
			}
			if conf.Engine.Go.Movetime > 0 {
				opts.MoveTime = conf.Engine.Go.Movetime
			}

			goResp := eng.Go(opts)

			makeMove(gameId, goResp.Bestmove)
		}
	}
}

func makeMove(gameId, move string) {
	log.Println("REQUEST", "bot/game/"+gameId+"/move/"+move)
	if move == "(none)" {
		return
	}
	request("POST", "bot/game/"+gameId+"/move/"+move)
}
