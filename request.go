package main

import (
	"log"
	"net/http"
)

var client = &http.Client{}

func request(method, path string) *http.Response {
	req, err := http.NewRequest(method, conf.Url+path, nil)
	if err != nil {
		log.Println("Failed request "+method+" "+path, err)
		return nil
	}
	req.Header.Add("Authorization", "Bearer "+conf.Token)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed request "+method+" "+path, err)
		return nil
	}
	return resp
}
