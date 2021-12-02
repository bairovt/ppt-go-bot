package main

import (
	"log"
	"strconv"
	"strings"

	// "encoding/json"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *api.BotAPI

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

	bot, err = api.NewBotAPI(Config.Bot.Token)
	if err != nil {
		log.Panic(err)
	}
	// bot.Debug = true

	webHookInfo, err := bot.GetWebhookInfo()
	if err != nil {
		log.Panic(err)
	}
	if webHookInfo.IsSet() {
		_, err := bot.Request(api.DeleteWebhookConfig{DropPendingUpdates: true})
		if err != nil {
			log.Panic(err)
		}
	}

	err = SetBotCommands()
	if err != nil {
		log.Panic("erron SetBotCommands", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := api.NewUpdate(0)
	updateConfig.Timeout = 60

	updatesChan := bot.GetUpdatesChan(updateConfig)

	for u := range updatesChan {
		ctx, err := getCtx(&u)
		if err != nil {
			log.Panic(err)
		}
		if u.Message != nil {
			if u.Message.IsCommand() {
				go handleCommand(strings.ToLower(u.Message.Command()), &u)
			} else {
				go handleMessage(&u)
			}
		} else if u.CallbackQuery != nil {
			subs := strings.Split(u.CallbackQuery.Data, ":")			
			switch subs[0] {
			case "set_role": 	
				go setRoleCb(ctx, &u, subs[1])
			default: go handleCallbackQuery(&u)
			}			
		}
	}
}

func getCtx(u *api.Update) (*Ctx, error) {
	var userKey string
	var user User
	if u.Message != nil {
			userKey = strconv.FormatInt(u.Message.From.ID, 10)
	} else if u.CallbackQuery != nil {
			userKey = strconv.FormatInt(u.CallbackQuery.From.ID, 10)
	}
	_, err := colUsers.ReadDocument(nil, userKey, &user)
	if err != nil {
		return nil, err
	}
	ctx := &Ctx{user}
	return ctx, nil
}

func handleMessage(u *api.Update) {
	msg := api.NewMessage(u.Message.Chat.ID, "Message.Text")

	msg.ReplyToMessageID = u.Message.MessageID
	var rec Rec

	key := u.Message.Text
	_, err := colRecs.ReadDocument(nil, key, &rec)
	if err != nil {
		log.Printf("err read doc %s: %v", key, err)
		msg.Text = err.Error()
		bot.Send(msg)
		return
	}
	msg.Text = rec.Body

	bot.Send(msg)
}
