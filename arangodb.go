package main

import (
	arangoDrv "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

var conn arangoDrv.Connection
var adb arangoDrv.Database
var arangoClient arangoDrv.Client
var colRecs arangoDrv.Collection
var colUsers arangoDrv.Collection

func initArangoDb() error {
	var err error
	conn, err = http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{Config.Adb.Url},
	})
	if err != nil {
		return err
	}
	arangoClient, err = arangoDrv.NewClient(arangoDrv.ClientConfig{
		Connection:     conn,
		Authentication: arangoDrv.BasicAuthentication(Config.Adb.Username, Config.Adb.Password),
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
