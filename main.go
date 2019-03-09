package main

import (
	"fmt"
	"net/http"
	"os"
)

const lichessUrl = "https://lichess.org/api/"
const botName = "dolegibot"

var apiKey string
var client = &http.Client{}

func main() {
	apiKey = os.Getenv("API_KEY")
	if len(apiKey) == 0 {
		fmt.Println("No API key found")
		return
	}

	ch := make(chan event)
	go streamEvent(ch)

	for e := range ch {
		fmt.Printf("%v\n", e)
		handleEvent(&e)
	}
}
