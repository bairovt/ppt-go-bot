package main

import (
	arangoDriver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

var conn arangoDriver.Connection
var adb arangoDriver.Database
var arangoClient arangoDriver.Client
var colRecs arangoDriver.Collection
var colUsers arangoDriver.Collection

func initArangoDb() error {
	var err error
	conn, err = http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{Config.Adb.Url},
	})
	if err != nil {
		return err
	}
	arangoClient, err = arangoDriver.NewClient(arangoDriver.ClientConfig{
		Connection:     conn,
		Authentication: arangoDriver.BasicAuthentication(Config.Adb.Username, Config.Adb.Password),
	})
	if err != nil {
		return err
	}
	adb, err = arangoClient.Database(nil, Config.Adb.Database)
	if err != nil {
		return err
	}
	// init collections
	colRecs, err = adb.Collection(nil, "Recs")
	if err != nil {
		return err
	}
	colUsers, err = adb.Collection(nil, "Users")
	if err != nil {
		return err
	}
	return nil
}
