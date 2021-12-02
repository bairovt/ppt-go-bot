package main

import (
	"fmt"
	"strings"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func commandHandler(ctx *Ctx, u *api.Update) error {
	command := strings.ToLower(u.Message.Command())
	if cmd, ok := BotCmds[command]; ok {
		err := cmd.fn(u)
		if err != nil {
			return err
		}
	} else {
		msg := api.NewMessage(u.Message.Chat.ID, "нет такой команды")
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func callbackQueryHandler(ctx *Ctx, u *api.Update) error {
	subs := strings.Split(u.CallbackQuery.Data, ":")			
	switch subs[0] {
	case "set_role": 	
		err := setRoleCb(ctx, u, subs[1])
		if err != nil {
			return err
		}
	default: 
		err := handleCallbackQuery(u)
		if err != nil {
			return err
		}
	}
	return nil
}

func messageHandler(ctx *Ctx, u *api.Update) error {
	msg := api.NewMessage(u.Message.Chat.ID, "Message.Text")

	msg.ReplyToMessageID = u.Message.MessageID
	// var rec Rec

	// key := u.Message.Text
	// _, err := colRecs.ReadDocument(nil, key, &rec)
	// if err != nil {
	// 	log.Printf("err read doc %s: %v", key, err)
	// 	msg.Text = err.Error()
	// 	bot.Send(msg)
	// 	return
	// }
	// msg.Text = rec.Body

	_, err := bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func myChatMemberHandler(ctx *Ctx, u *api.Update) error {
	// switch u.MyChatMember.NewChatMember.Status {
	// case "kicked":
	// case "member":
	// }
	fmt.Printf("New status: %#v\n", *&u.MyChatMember.NewChatMember.Status)
	return nil
}