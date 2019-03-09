package main

import (
	"github.com/dolegi/uci"
	"log"
	"net/http"
	"time"
)

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
	req, err := http.NewRequest("POST", lichessUrl+"challenge/"+c.Id+"/accept", nil)
	if err != nil {
		log.Println("Failed to accept challenge", err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+apiKey)
	_, err = client.Do(req)
	if err != nil {
		log.Println("Failed to accept challenge", err)
	}
}

func startGame(gameId string) {
	eng, err := uci.NewEngine("../stockfish")
	if err != nil {
		log.Println("startGame failed to create new engine", err)
		return
	}
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

	streamGame(gameId, eng)
}
