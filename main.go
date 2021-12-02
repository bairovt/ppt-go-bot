package main

import (
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
	ctx, err := getCtx(&u)
	if err != nil {
		log.Panic(err)
	}
	if u.Message != nil {
		if u.Message.IsCommand() {
			err := commandHandler(ctx, &u)
			if err != nil {
				log.Panic(err)
			}
		} else {
			err := messageHandler(ctx, &u)
			if err != nil {
				log.Panic(err)
			}
		}
	} else if u.CallbackQuery != nil {
		err := callbackQueryHandler(ctx, &u)
		if err != nil {
			log.Panic(err)
		}
					
	} else if u.MyChatMember != nil {
		// switch u.MyChatMember.NewChatMember.Status {
		// case "kicked":
		// case "member":
		// }
		fmt.Printf("New status: %#v\n", *&u.MyChatMember.NewChatMember.Status)
	}
}

// func getCtxAndHandler(u *api.Update) (*Ctx, func(*Ctx, *api.Update) error, error) {
// 	var userKey string
// 	var user User
// 	if u.Message != nil {
// 			userKey = strconv.FormatInt(u.Message.From.ID, 10)
// 	} else if u.CallbackQuery != nil {
// 			userKey = strconv.FormatInt(u.CallbackQuery.From.ID, 10)
// 	} else if u.MyChatMember != nil {	// bot was blocked/unblocked by user
// 		userKey = strconv.FormatInt(u.MyChatMember.From.ID, 10)		
// 	} else {
// 		fmt.Printf("unknown update:\n%#v", u)
// 	}
// 	if userKey != "" {
// 		_, err := colUsers.ReadDocument(nil, userKey, &user)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	ctx := &Ctx{user}
// 	return ctx, nil
// }

func getCtx(u *api.Update) (*Ctx, error) {
	var userKey string
	var user User
	if u.Message != nil {
			userKey = strconv.FormatInt(u.Message.From.ID, 10)
	} else if u.CallbackQuery != nil {
			userKey = strconv.FormatInt(u.CallbackQuery.From.ID, 10)
	} else if u.MyChatMember != nil {	// bot was blocked/unblocked by user
		userKey = strconv.FormatInt(u.MyChatMember.From.ID, 10)		
	} else {
		fmt.Printf("unknown update:\n%#v", u)
	}
	if userKey != "" {
		_, err := colUsers.ReadDocument(nil, userKey, &user)
		if err != nil {
			return nil, err
		}
	}
	ctx := &Ctx{user}
	return ctx, nil
}


