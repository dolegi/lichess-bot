package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type user struct {
	Id		string `json:"id"`
	Name		string `json:"name"`
	Title		string `json:"title"`
	Online		bool   `json:"online"`
	Playing		bool   `json:"playing"`
	Streaming	bool   `json:"streaming"`
	Patron		bool   `json:"patron"`
}

func getUsersStatus(ids string) []user {
	resp := request("GET", "users/status?ids="+ids)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var users []user
	err = json.Unmarshal(body, &users)
	if err != nil {
		log.Fatal(err)
	}
	return users
}
