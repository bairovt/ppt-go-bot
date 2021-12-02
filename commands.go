package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	arangoDrv "github.com/arangodb/go-driver"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotCmd struct {
	BotCommand api.BotCommand
	fn         func(u *api.Update) error
	set        bool
}

var BotCmds = map[string]BotCmd{
	"start": {
		api.BotCommand{Command: "start",
			Description: "start"},
		startCmd,
		false,
	},
	"menu": {
		api.BotCommand{Command: "menu",
			Description: "главное меню"},
		menuCmd,
		true,
	},
	"help": {
		api.BotCommand{Command: "help",
			Description: "справка"},
		helpCmd,
		true,
	},		
	"test": {
		api.BotCommand{Command: "test",
			Description: "test only"},
		testCmd,
		false,
	},
}

func handleCommand(command string, upd *api.Update) {
	if cmd, ok := BotCmds[command]; ok {
		cmd.fn(upd)
	} else {
		msg := api.NewMessage(upd.Message.Chat.ID, "нет такой команды")
		bot.Send(msg)
	}
}

func SetBotCommands() error {
	commandsToSet := make([]api.BotCommand, 0)
	for _, botCmd := range BotCmds {
		if botCmd.set {
			commandsToSet = append(commandsToSet, botCmd.BotCommand)
		}
	}
	config := api.NewSetMyCommands(commandsToSet...)
	_, err := bot.Request(config)
	if err != nil {
		return err
	}
	return nil
}

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
	kbMsg := api.NewMessage(u.Message.Chat.ID, "выберите Вашу роль:")
	// btn1 := api.NewKeyboardButton("btn1")
	// btn1 := api.InlineKeyboardButton{
	// 	Text: "btn1",
	// 	CallbackData: "btn1_cb_data",
	// }
	setRoleIkb := api.NewInlineKeyboardMarkup(
		api.NewInlineKeyboardRow(
			api.NewInlineKeyboardButtonData("я - пассажир (отправлю груз)", "set_role:P"),		
		),
		api.NewInlineKeyboardRow(			
			api.NewInlineKeyboardButtonData("я - водитель", "set_role:D"),
		),
	)
	kbMsg.ReplyMarkup = setRoleIkb

	// metaStr, err := json.MarshalIndent(meta, "", " ")
	_, err = bot.Send(kbMsg)
	// if err != nil {
	// 	return err
	// }
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
	b, _ := json.Marshal(u)
	fmt.Println(string(b))
	txtmsg := api.NewMessage(u.Message.Chat.ID, "test cmd")
	_, err := bot.Send(txtmsg)
	if err != nil {
		return err
	}
	return nil
}
