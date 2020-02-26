package main

import (
	macaron "gopkg.in/macaron.v1"
)

func pageView(ctx *macaron.Context, tu TelegramUpdate) string {
	executeBotCommand(tu)

	return "OK"
}
