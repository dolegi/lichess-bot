package main

import (
	"encoding/json"
	"github.com/dolegi/uci"
	"log"
	"time"
)

type event struct {
	Type      string    `json:"type"`
	Challenge challenge `json:"challenge"`
	Game      game      `json:"game"`
}

type game struct {
	Id string `json:"id"`
}

type challenge struct {
	Id      string `json:"id"`
	Status  string `json:"status"`
	Variant struct {
		Key string `json:"key"`
	}
	Rated bool   `json:"rated"`
	Speed string `json:"speed"`
	Color string `json:"color"`
}

func streamEvent() {
	resp := request("GET", "stream/event")
	dec := json.NewDecoder(resp.Body)

	for dec.More() {
		var e event
		err := dec.Decode(&e)
		if err != nil {
			log.Println(err)
		} else {
			handleEvent(&e)
		}
	}
}

func handleEvent(e *event) {
	switch e.Type {
	case "challenge":
		if validChallenge(&e.Challenge) {
			acceptChallenge(&e.Challenge)
		} else {
			log.Println("Invalid challenge", e.Challenge)
		}
	case "gameStart":
		go startGame(e.Game.Id)
	default:
		log.Printf("Unhandled Event %v\n", e)
	}
}

func validChallenge(c *challenge) bool {
	return c.Status == "created" &&
		c.Variant.Key == "standard"
	// e.Challenge.Rated == true &&
	// e.Challenge.Speed == "blitz"
}

func acceptChallenge(c *challenge) {
	request("POST", "challenge/"+c.Id+"/accept")
}

func startGame(gameId string) {
	retries := 0
	for _ = range time.Tick(50 * time.Millisecond) {
		if eng.IsReady() {
			break
		}
		if retries > 10 {
			break
		}
		retries++
	}
	if retries > 10 {
		log.Println("startGame failed to find ready engine. Max retries exceeded")
		return
	}

	streamGame(gameId)
}
