package main

import (
	"fmt"
	"strings"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func commandHandler(ctx *Ctx, upd *api.Update) error {
	command := strings.ToLower(upd.Message.Command())
	if cmd, ok := BotCmds[command]; ok {
		err := cmd.fn(upd)
		if err != nil {
			return err
		}
	} else {
		msg := api.NewMessage(upd.Message.Chat.ID, "нет такой команды")
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func callbackQueryHandler(ctx *Ctx, upd *api.Update) error {
	subs := strings.Split(upd.CallbackQuery.Data, ":")
	switch subs[0] {
	case "set_role":
		err := setRoleCb(ctx, upd, subs[1])
		if err != nil {
			return err
		}
	default:
		err := handleCallbackQuery(upd)
		if err != nil {
			return err
		}
	}
	return nil
}

func messageHandler(ctx *Ctx, upd *api.Update) error {
	msg := api.NewMessage(upd.Message.Chat.ID, "Message.Text")

	msg.ReplyToMessageID = upd.Message.MessageID
	// var rec RecDoc

	// key := upd.Message.Text
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

func myChatMemberHandler(ctx *Ctx, upd *api.Update) error {
	// switch upd.MyChatMember.NewChatMember.Status {
	// case "kicked":
	// case "member":
	// }
	fmt.Printf("New status: %#v\n", *&upd.MyChatMember.NewChatMember.Status)
	return nil
}
