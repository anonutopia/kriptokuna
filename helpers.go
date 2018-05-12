package main

import (
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

func newPageData(ctx *macaron.Context, sess session.Store) {
	ctx.Data["ProjectName"] = PROJECT_NAME
	ctx.Data["URI"] = ctx.Req.RequestURI

	var hacktivists []Hacktivist
	carsharing := 0
	blocked := 0

	db.Find(&hacktivists)

	for _, h := range hacktivists {
		if h.Type == TYPE_BLOCKED {
			blocked++
		} else if h.Type == TYPE_USER || h.Type == TYPE_DRIVER {
			carsharing++
		}
	}

	ctx.Data["Blocked"] = blocked + 13
	ctx.Data["Carsharing"] = carsharing + 9
	ctx.Data["Hacktivists"] = len(hacktivists) + 21
}

type JsonResponse struct {
	Success      bool   `json:"success"`
	Response     string `json:"response"`
	ErrorMessage string `json:"errorMessage"`
}
