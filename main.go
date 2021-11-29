package main

import (
	"log"
	"strings"

	// "encoding/json"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *api.BotAPI

func main() {	
	var err error
	err = readConf()
	if err != nil {
		log.Panic("erron readConf", err)
	}
	err = initArangoDb()
	if err != nil {
		log.Panic("erron initArangoDb: ", err)
	}

	Bot, err = api.NewBotAPI(Config.Bot.Token)
	if err != nil {
		log.Panic(err)
	}
	// Bot.Debug = true

	webHookInfo, err := Bot.GetWebhookInfo()
	if err != nil {
		log.Panic(err)
	}
	if webHookInfo.IsSet() {
		_, err := Bot.Request(api.DeleteWebhookConfig{DropPendingUpdates: true})
		if err != nil {
			log.Panic(err)
		}
	}

	err = SetBotCommands()
	if err != nil {
		log.Panic("erron SetBotCommands", err)
	}

	log.Printf("Authorized on account %s", Bot.Self.UserName)		

	updateConfig := api.NewUpdate(0)
	updateConfig.Timeout = 60

	updatesChan := Bot.GetUpdatesChan(updateConfig)

	for update := range updatesChan {
		if update.Message != nil {
			cmd := update.Message.Command()
			if cmd != "" {
				go RunCommand(strings.ToLower(cmd), &update)
				continue
			}
			go handleMessage(&update)
		}
	}
}

func handleMessage(u *api.Update) {
	msg := api.NewMessage(u.Message.Chat.ID, "Message.Text")

	msg.ReplyToMessageID = u.Message.MessageID
	var rec Rec
	colRecs, err := adb.Collection(nil, "Recs")
	if err != nil {
		log.Panic(err)
	}
	key := u.Message.Text
	_, err = colRecs.ReadDocument(nil, key, &rec)
	if err != nil {
		log.Printf("err read doc %s: %v", key, err)
		msg.Text = err.Error()
		Bot.Send(msg)
		return
	}
	msg.Text = rec.Body

	Bot.Send(msg)
}
