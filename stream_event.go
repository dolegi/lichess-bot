package main

import (
	"encoding/json"
	"log"
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
			acceptChallenge(e.Challenge.Id)
		} else {
			log.Println("Invalid challenge", e.Challenge)
		}
	case "gameStart":
		go streamGame(e.Game.Id)
	default:
		log.Printf("Unhandled Event %v\n", e)
	}
}

func validChallenge(c *challenge) bool {
	return c.Status == "created" &&
		includes(conf.Challenge.Variants, c.Variant.Key) &&
		includes(conf.Challenge.Speeds, c.Speed) &&
		(c.Rated == true && includes(conf.Challenge.Modes, "rated") ||
			c.Rated == false && includes(conf.Challenge.Modes, "casual"))
}

func includes(arr []string, x string) bool {
	for _, a := range arr {
		if a == x {
			return true
		}
	}
	return false
}

func acceptChallenge(challengeId string) {
	request("POST", "challenge/"+challengeId+"/accept")
}
