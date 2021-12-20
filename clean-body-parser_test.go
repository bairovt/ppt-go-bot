package main

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

type Msg struct {
	Body        string   `json:"body"`
	CleanedBody string   `json:"cleanedBody"`
	Route       []string `json:"route"`
}

func TestCleanBodyParser(t *testing.T) {
	f, err := os.Open("./tests/msgs.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var msgs []Msg
	recsDecoder := json.NewDecoder(f)
	err = recsDecoder.Decode(&msgs)
	if err != nil {
		log.Panic(err)
	}
	for _, msg := range msgs {
		cleanedBody, err := CleanBodyParser(msg.Body)
		if err != nil {
			t.Errorf("err CleanBodyParser on body: %s", msg.Body)
		}
		if msg.CleanedBody != "" && cleanedBody != msg.CleanedBody {
			t.Errorf("%s==%s", cleanedBody, msg.CleanedBody)
		}
	}
}
