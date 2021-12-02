package main

import (
	"encoding/json"
	"sort"
	"strconv"

	arangoDrv "github.com/arangodb/go-driver"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotCmd struct {
	BotCommand api.BotCommand
	fn         func(u *api.Update) error
	order        int8
}

var BotCmds = map[string]BotCmd{
	"start": {
		api.BotCommand{Command: "start",
			Description: "start"},
		startCmd,
		0,
	},
	"menu": {
		api.BotCommand{Command: "menu",
			Description: "главное меню"},
		menuCmd,
		3,
	},
	"help": {
		api.BotCommand{Command: "help",
			Description: "справка"},
		helpCmd,
		2,
	},		
	"test": {
		api.BotCommand{Command: "test",
			Description: "test only"},
		testCmd,
		1,
	},
}

func SetBotCommands() error {
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
	
	config := api.NewSetMyCommands(BotCommandList...)
	_, err := bot.Request(config)
	if err != nil {
		return err
	}
	return nil
}

// command handlers
func startCmd(u *api.Update) error {
	var userKey = strconv.FormatInt(u.Message.From.ID, 10)
	var user = User{
		Key:       userKey,
		ID:        u.Message.From.ID,
		ChatID:    u.Message.Chat.ID,
		FirstName: u.Message.From.FirstName,
		LastName:  u.Message.From.LastName,
		UserName:  u.Message.From.UserName,
		StartDate: int64(u.Message.Date),
		IsBot:     u.Message.From.IsBot,
	}
	_, err := colUsers.CreateDocument(nil, &user)
	if arangoDrv.IsArangoErrorWithCode(err, 409) { // conflict unique
		user.StartDate = 0 // to omitempty
		user.RestartDate = int64(u.Message.Date)
		_, err = colUsers.UpdateDocument(nil, userKey, user)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	kbMsg := api.NewMessage(u.Message.Chat.ID, "Выберите свою роль:")

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

func menuCmd(u *api.Update) error {
	txtmsg := api.NewMessage(u.Message.Chat.ID, "menu cmd")
	_, err := bot.Send(txtmsg)
	if err != nil {
		return err
	}
	return nil
}

func helpCmd(u *api.Update) error {
	txtmsg := api.NewMessage(u.Message.Chat.ID, "help cmd")
	_, err := bot.Send(txtmsg)
	if err != nil {
		return err
	}
	return nil
}

func testCmd(u *api.Update) error {
	_, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	// time.Sleep(time.Second * 15)
	txtmsg := api.NewMessage(u.Message.Chat.ID, "test cmd")
	_, err = bot.Send(txtmsg)
	if err != nil {
		return err
	}
	return nil
}
