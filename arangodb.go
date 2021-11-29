package main

import (
	adbDrv "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

var conn adbDrv.Connection
var adb adbDrv.Database
var arangoClient adbDrv.Client

func initArangoDb() error {
	var err error
	conn, err = http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{Config.Adb.Url},
	})
	if err != nil {
		return err
	}
	arangoClient, err = adbDrv.NewClient(adbDrv.ClientConfig{
		Connection:     conn,
		Authentication: adbDrv.BasicAuthentication(Config.Adb.Username, Config.Adb.Password),
	})
	if err != nil {
		return err
	}
	adb, err = arangoClient.Database(nil, Config.Adb.Database)
	if err != nil {
		return err
	}
	return nil
}
