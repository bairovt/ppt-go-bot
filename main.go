package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"

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

	for update := range updatesChan {
		go updateHandler(update)
	}
}

func updateHandler(u api.Update) {
	ctx, handler, err := getCtxAndHandler(&u)
	if err != nil {
		log.Panic(err)
	}
	err = handler(ctx, &u)
	if err != nil {
		log.Panic(err)
	}
}

func getCtxAndHandler(u *api.Update) (ctx *Ctx, handler func(*Ctx, *api.Update) error, err error) {
	var userKey string
	var user User
	if u.Message != nil {
		userKey = strconv.FormatInt(u.Message.From.ID, 10)
		if u.Message.IsCommand() {
			handler = commandHandler
		} else {
			handler = messageHandler
		}
	} else if u.CallbackQuery != nil {
		userKey = strconv.FormatInt(u.CallbackQuery.From.ID, 10)
		handler = callbackQueryHandler
	} else if u.MyChatMember != nil {
		userKey = strconv.FormatInt(u.MyChatMember.From.ID, 10)
		handler = myChatMemberHandler
	} else {
		// todo refactor not to panic
		err = errors.New(fmt.Sprintf("unknown update:\n%#v", u))
		return nil, nil, err
	}
	if userKey != "" {
		_, err := colUsers.ReadDocument(nil, userKey, &user)
		if err != nil {
			return nil, nil, err
		}
	} else {
		err = errors.New("empty userKey")
		return nil, nil, err
	}
	ctx = &Ctx{user}
	return ctx, handler, nil
}
