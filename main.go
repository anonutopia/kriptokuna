package main

import (
	"log"
	"net/http"

	"github.com/go-macaron/binding"
	"github.com/jinzhu/gorm"
	"gopkg.in/macaron.v1"
)

var m *macaron.Macaron

var conf *Config

var db *gorm.DB

func main() {
	conf = initConfig()

	m = initMacaron()

	db = initDb()

	m.Get("/", newPageData, homeView)
	m.Get("/o-kriptokuni/", newPageData, kriptokunaView)
	m.Get("/blokirani/", newPageData, manifestView)
	m.Get("/pitanja/", newPageData, faqView)
	m.Get("/volontiraj/", newPageData, volontirajView)
	m.Get("/plan/", newPageData, planView)
	m.Get("/novcanik/", newPageData, novcanikView)
	m.Get("/anote/", newPageData, anoteView)
	m.Get("/transparentnost/", newPageData, transparentnostView)
	m.Get("/kontakt/", newPageData, kontaktView)

	// m.Post("/", binding.Bind(SignupForm{}), newPageData, signupView)
	m.Post("/volontiraj/", binding.Bind(HackerSignupForm{}), newPageData, volontirajPostView)
	m.Post("/kontakt/", binding.Bind(ContactForm{}), newPageData, kontaktViewPost)

	m.NotFound(view404)

	// m.Run()
	log.Println("Server is running...")
	http.ListenAndServe("0.0.0.0:4001", m)
}
