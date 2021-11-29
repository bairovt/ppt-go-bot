package main

import (
	"encoding/json"
	"fmt"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotCmd struct {
	BotCommand api.BotCommand
	fn   func(u *api.Update) error
	set bool
}

var BotCmds = map[string]BotCmd{
	"start": {
		api.BotCommand{Command: "start",
			Description: "start"},
		func(u *api.Update) error {
			// msg := api.NewDeleteMessage(u.Message.Chat.ID,  u.Message.MessageID)
			txtmsg := api.NewMessage(u.Message.Chat.ID, "start cmd")
			_, err := Bot.Send(txtmsg)
			if err != nil {
				return err
			}
			return nil
		},
		false,
	},
	"menu": {
		api.BotCommand{Command: "menu",
			Description: "главное меню"},
		func(u *api.Update) error {
			txtmsg := api.NewMessage(u.Message.Chat.ID, "menu cmd")
			_, err := Bot.Send(txtmsg)
			if err != nil {
				return err
			}
			return nil
		},
		true,
	},
	"help": {
		api.BotCommand{Command: "help",
			Description: "справка"},

		func(u *api.Update) error {
			txtmsg := api.NewMessage(u.Message.Chat.ID, "help cmd")
			_, err := Bot.Send(txtmsg)
			if err != nil {
				return err
			}
			return nil
		},
		true,
	},
	"test": {
		api.BotCommand{Command: "test",
			Description: "test only"},

		func(u *api.Update) error {
			b, _ := json.Marshal(u)
			fmt.Println(string(b))
			txtmsg := api.NewMessage(u.Message.Chat.ID, "test cmd")
			_, err := Bot.Send(txtmsg)
			if err != nil {
				return err
			}
			return nil
		},
		false,
	},
}

func RunCommand(command string, upd *api.Update) {
	if cmd, ok := BotCmds[command]; ok {
		cmd.fn(upd)
	} else {
		msg := api.NewMessage(upd.Message.Chat.ID, "нет такой команды")
		Bot.Send(msg)
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
	_, err := Bot.Request(config)
	if err != nil {		
		return err
	}
	return nil
}
