package main

import (
	"github.com/caddyserver/certmagic"
	macaron "gopkg.in/macaron.v1"
)

func initMacaron() *macaron.Macaron {
	m := macaron.New()

	m.Use(macaron.Renderer())

	if !conf.Dev {
		certmagic.DefaultACME.Email = conf.EmailAddress
		go certmagic.HTTPS([]string{conf.Hostname}, m)
	} else {
		go m.Run("0.0.0.0", 5000)
	}

	return m
}
