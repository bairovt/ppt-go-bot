package db

import (
	"log"

	arangoDriver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"

	"ppt-go-bot/conf"
)

var Conn arangoDriver.Connection
var DB arangoDriver.Database
var ArangoClient arangoDriver.Client
var ColRecs arangoDriver.Collection
var ColUsers arangoDriver.Collection

func init() {
	var err error
	Conn, err = http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{conf.Config.Adb.Url},
	})
	if err != nil {
		log.Panic(err)
	}
	ArangoClient, err = arangoDriver.NewClient(arangoDriver.ClientConfig{
		Connection:     Conn,
		Authentication: arangoDriver.BasicAuthentication(conf.Config.Adb.Username, conf.Config.Adb.Password),
	})
	if err != nil {
		log.Panic(err)
	}
	DB, err = ArangoClient.Database(nil, conf.Config.Adb.Database)	
	if err != nil {
		log.Panic(err)
	}
	// init collections
	ColRecs, err = DB.Collection(nil, "Recs")
	if err != nil {
		log.Panic(err)
	}
	ColUsers, err = DB.Collection(nil, "Users")
	if err != nil {
		log.Panic(err)
	}	
}
