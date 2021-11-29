package main

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Configuration struct {
	Bot struct {
		Token string `json:"token"`
	} `json:"bot"`
	Adb struct {
		Url      string `json:"url"`
		Database string `json:"database"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"adb"`
}

var Config Configuration

func readConf() error {
	err := cleanenv.ReadConfig("./config/development.json", &Config)
	if err != nil {
		return err
	}
	return nil
}

// func GetConfig() Configuration {
// 	var config Configuration
// 	var configFile, err = os.Open("./config/development.json")
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	decoder := json.NewDecoder(configFile)
// 	decoder.Decode(&config)

// 	return config
// }
