package main

import (
	"encoding/json"
	"github.com/dolegi/uci"
	"io/ioutil"
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
	Id         string `json:"id"`
	Status     string `json:"status"`
	Challenger struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Title       string `json:"title"`
		Rating      int    `json:"rating"`
		Provisional bool   `json:"provisional"`
		Online      bool   `json:"online"`
		Lag         int    `json:"lag"`
	}
	DestUser struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Title       string `json:"title"`
		Rating      int    `json:"rating"`
		Provisional bool   `json:"provisional"`
		Online      bool   `json:"online"`
		Lag         int    `json:"lag"`
	}
	Variant struct {
		Key string `json:"key"`
	}
	Rated bool   `json:"rated"`
	Speed string `json:"speed"`
	Color string `json:"color"`
}

func streamEvent(eng *uci.Engine) {
	resp := request("GET", "stream/event")
	dec := json.NewDecoder(resp.Body)

	for dec.More() {
		var e event
		err := dec.Decode(&e)
		if err != nil {
			log.Println(err)
		} else {
			handleEvent(&e, eng)
		}
	}
}

func handleEvent(e *event, eng *uci.Engine) {
	switch e.Type {
	case "challenge":
		handleChallengeEvent(e)
	case "gameStart":
		streamGame(e.Game.Id, eng)
	default:
		log.Printf("Unhandled Event %v\n", e)
	}
}

func handleChallengeEvent(e *event) {
	challengeId := e.Challenge.Id
	if validChallenge(&e.Challenge) && !gameInProgress() {
		log.Println("Accepting challenge", e.Challenge)
		resp := request("POST", "challenge/"+challengeId+"/accept")
		resp.Body.Close()
	} else {
		log.Println("Declining challenge", e.Challenge)
		resp := request("POST", "challenge/"+challengeId+"/decline")
		resp.Body.Close()
	}
}

func validChallenge(c *challenge) bool {
	return c.Status == "created" &&
		c.Challenger.Online == true &&
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

func gameInProgress() bool {
	type status struct {
		Playing bool
	}
	s := []status{}
	resp := request("GET", "users/status?ids="+conf.Botname)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("gameInProgress failed to read response body")
		return false
	}

	json.Unmarshal(body, &s)

	return s[0].Playing
}
