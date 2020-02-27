package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

var client = &http.Client{}

func chat(gameId, room, text string) *http.Response {
	method := "POST"
	path := "bot/game/"+gameId+"/chat"
	line, _ := json.Marshal(map[string]string{"room": room, "text": text})
	req, err := http.NewRequest(method, conf.Url+path, bytes.NewBuffer(line))
	req.Header.Add("Authorization", "Bearer "+conf.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed request "+method+" "+path, err)
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Response %d %s %s %s", resp.StatusCode, http.StatusText(resp.StatusCode), method, path)
	}
	return resp
}

func request(method, path string) *http.Response {
	req, err := http.NewRequest(method, conf.Url+path, nil)
	req.Header.Add("Authorization", "Bearer "+conf.Token)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed request "+method+" "+path, err)
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Response %d %s %s %s", resp.StatusCode, http.StatusText(resp.StatusCode), method, path)
	}
	return resp
}
