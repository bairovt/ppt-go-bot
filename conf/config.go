package conf

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

func init() {
	err := cleanenv.ReadConfig("./conf/development.json", &Config)
	if err != nil {
		panic(err)
	}	
}
