package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
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

	go runHttpServer()

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

func updateHandler(upd api.Update) {
	if upd.Message != nil && upd.Message.Command() == "start" {
		err := startCmdHandler(&upd)
		if err != nil {
			log.Panic(err)
		}
		return
	}

	ctx, handler, err := getCtxAndHandler(&upd)
	if err != nil {
		log.Panic(err)
	}
	err = handler(ctx, &upd)
	if err != nil {
		log.Panic(err)
	}
}

func getCtxAndHandler(upd *api.Update) (ctx *Ctx, handler func(*Ctx, *api.Update) error, err error) {
	var userKey string
	var user UserDoc
	if upd.Message != nil {
		userKey = strconv.FormatInt(upd.Message.From.ID, 10)
		if upd.Message.IsCommand() {
			handler = commandHandler
		} else {
			handler = messageHandler
		}
	} else if upd.CallbackQuery != nil {
		userKey = strconv.FormatInt(upd.CallbackQuery.From.ID, 10)
		handler = callbackQueryHandler
	} else if upd.MyChatMember != nil {
		userKey = strconv.FormatInt(upd.MyChatMember.From.ID, 10)
		handler = myChatMemberHandler
	} else {
		// todo refactor not to panic
		err = errors.New(fmt.Sprintf("unknown update:\n%#v", upd))
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

func runHttpServer(){
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})
	err := http.ListenAndServe(":1818", nil)
	if err != nil {
		log.Fatal("can not start http sever")
	}
	log.Println("Http server is started")
}