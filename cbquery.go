package main

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleCallbackQuery(upd *api.Update) error {
	cb := api.NewCallbackWithAlert(upd.CallbackQuery.ID, upd.CallbackQuery.Data)
	if _, err := bot.Request(cb); err != nil {
		return err
	}
	return nil
}

func setRoleCb(ctx *Ctx, upd *api.Update, role string) error {
	type RoleUpd struct {
		PptRole string `json:"role"`
	}
	roleUpd := RoleUpd{role}
	_, err := colUsers.UpdateDocument(nil, ctx.user.Key, &roleUpd)
	if err != nil {
		return err
	}
	cb := api.NewCallback(upd.CallbackQuery.ID, "роль установлена")
	if _, err := bot.Request(cb); err != nil {
		return err
	}
	return nil
}
