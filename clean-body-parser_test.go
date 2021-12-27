package main

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

type MsgCase struct {
	Body        string   `json:"body"`
	CleanedBody string   `json:"cleanedBody"`
	Route       []string `json:"route"`
}

var msgCases []MsgCase

func TestMain(m *testing.M){	
	f, err := os.Open("./test/msg-cases.json")
	if err != nil {
		log.Fatal(err)
	}	
	recsDecoder := json.NewDecoder(f)
	err = recsDecoder.Decode(&msgCases)
	if err != nil {
		log.Panic(err)
	}

	code := m.Run()	
	f.Close()
	os.Exit(code)
}

func TestCleanBodyParser(t *testing.T) {	
	for _, testCase := range msgCases {
		cleanedBody, err := CleanBodyParser(testCase.Body)
		if err != nil {
			t.Errorf("err CleanBodyParser on body: %s", testCase.Body)
		}
		if testCase.CleanedBody != "" && cleanedBody != testCase.CleanedBody {
			t.Errorf("%s===%s", cleanedBody, testCase.CleanedBody)
		}
	}
}
