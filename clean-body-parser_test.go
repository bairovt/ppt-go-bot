package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

type TestRec struct {
	Body string `json:"body"`
	CleanedBody string `json:"cleanedBody"`
	Route []string `json:"route"`
}
func TestCleanBodyParser(t *testing.T) {
	f, err := os.Open("./tests/msgs.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	
	var testRecs []TestRec 
	recsDecoder := json.NewDecoder(f)
	err = recsDecoder.Decode(&testRecs)
	if err != nil {
		log.Panic(err)
	}
fmt.Print(testRecs)
}