package main

import (
	"encoding/json"
	"log"
	"net/http"
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

func streamEvent(ch chan event) {
	req, err := http.NewRequest("GET", lichessUrl+"stream/event", nil)
	req.Header.Add("Authorization", "Bearer "+apiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	dec := json.NewDecoder(resp.Body)

	for dec.More() {
		var e event
		err := dec.Decode(&e)
		if err != nil {
			log.Println(err)
		}

		ch <- e
	}
}
