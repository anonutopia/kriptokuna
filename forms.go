package main

import (
	// "log"

	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
)

type SignupForm struct {
	Email string `form:"email" binding:"Required"`
	Type  int    `form:"type"`
}

func (cf SignupForm) Error(ctx *macaron.Context, errs binding.Errors) {
	ctx.Data["Errors"] = errs
}

type HackerSignupForm struct {
	Email string `form:"email" binding:"Required"`
	Type  string `form:"type"`
}

func (cf HackerSignupForm) Error(ctx *macaron.Context, errs binding.Errors) {
	ctx.Data["Errors"] = errs
}

type ContactForm struct {
	Name    string `form:"name" binding:"Required"`
	Email   string `form:"email" binding:"Required"`
	Subject string `form:"subject"`
	Message string `form:"message" binding:"Required"`
	Captcha string `form:"captcha" binding:"Required"`
}

func (cf ContactForm) Error(ctx *macaron.Context, errs binding.Errors) {
	ctx.Data["Errors"] = errs
}
