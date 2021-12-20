package main

import (
	"encoding/json"
	"sort"
	"strconv"

	"ppt-go-bot/db"

	arangoDrv "github.com/arangodb/go-driver"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotCmd struct {
	BotCommand api.BotCommand
	fn         func(upd *api.Update) error
	order      int8
}

var BotCmds = map[string]BotCmd{
	"menu": {
		api.BotCommand{Command: "menu",
			Description: "главное меню"},
		menuCmdHandler,
		3,
	},
	"help": {
		api.BotCommand{Command: "help",
			Description: "справка"},
		helpCmdHandler,
		2,
	},
	"test": {
		api.BotCommand{Command: "test",
			Description: "test only"},
		testCmdHandler,
		0,
	},
}

func getBotCommandList() []api.BotCommand {
	botCmdList := make([]BotCmd, 0, len(BotCmds))
	for _, botCmd := range BotCmds {
		if botCmd.order != 0 {
			botCmdList = append(botCmdList, botCmd)
		}
	}
	sort.Slice(botCmdList, func(i, j int) bool {
		return botCmdList[i].order < botCmdList[j].order
	})
	BotCommandList := make([]api.BotCommand, 0, len(botCmdList))
	for _, botCmd := range botCmdList {
		BotCommandList = append(BotCommandList, botCmd.BotCommand)
	}
	return BotCommandList
}

func SetBotCommands() error {
	BotCommandList := getBotCommandList()
	config := api.NewSetMyCommands(BotCommandList...)
	_, err := bot.Request(config)
	if err != nil {
		return err
	}
	return nil
}

// command handlers
func startCmdHandler(upd *api.Update) error {
	var userKey = strconv.FormatInt(upd.Message.From.ID, 10)
	var user = db.UserDoc{
		Key:       userKey,
		ID:        upd.Message.From.ID,
		ChatID:    upd.Message.Chat.ID,
		FirstName: upd.Message.From.FirstName,
		LastName:  upd.Message.From.LastName,
		UserName:  upd.Message.From.UserName,
		StartDate: int64(upd.Message.Date),
		IsBot:     upd.Message.From.IsBot,
	}
	_, err := db.ColUsers.CreateDocument(nil, &user)
	if arangoDrv.IsArangoErrorWithCode(err, 409) { // conflict unique
		user.StartDate = 0 // to omitempty
		user.RestartDate = int64(upd.Message.Date)
		_, err = db.ColUsers.UpdateDocument(nil, userKey, user)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	kbMsg := api.NewMessage(upd.Message.Chat.ID, "Выберите свою роль:")

	setRoleIkb := api.NewInlineKeyboardMarkup(
		api.NewInlineKeyboardRow(
			api.NewInlineKeyboardButtonData("я - пассажир / отправлю груз 🙋‍♀🙋‍♂📦", "set_role:P"),
		),
		api.NewInlineKeyboardRow(
			api.NewInlineKeyboardButtonData("я - водитель, возьму пассаж-в/груз 🚙🚛", "set_role:D"),
		),
	)
	kbMsg.ReplyMarkup = setRoleIkb
	_, err = bot.Send(kbMsg)
	if err != nil {
		return err
	}
	return nil
}

func menuCmdHandler(upd *api.Update) error {
	txtmsg := api.NewMessage(upd.Message.Chat.ID, "menu cmd")
	_, err := bot.Send(txtmsg)
	if err != nil {
		return err
	}
	return nil
}

func helpCmdHandler(upd *api.Update) error {
	txtmsg := api.NewMessage(upd.Message.Chat.ID, "help cmd")
	_, err := bot.Send(txtmsg)
	if err != nil {
		return err
	}
	return nil
}

func testCmdHandler(upd *api.Update) error {
	_, err := json.Marshal(upd)
	if err != nil {
		panic(err)
	}
	// time.Sleep(time.Second * 15)
	txtmsg := api.NewMessage(upd.Message.Chat.ID, "test cmd")
	_, err = bot.Send(txtmsg)
	if err != nil {
		return err
	}
	return nil
}
