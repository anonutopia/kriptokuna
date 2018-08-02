package main

import (
	"fmt"
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
	m.Get("/airdrop/", newPageData, airdropView)
	m.Get("/o-kriptokuni/", newPageData, kriptokunaView)
	m.Get("/blokirani/", newPageData, manifestView)
	m.Get("/pitanja/", newPageData, faqView)
	m.Get("/pridruzi-se/", newPageData, volontirajView)
	m.Get("/plan/", newPageData, planView)
	m.Get("/novcanik/", newPageData, novcanikView)
	m.Get("/anote/", newPageData, anoteView)
	m.Get("/transparentnost/", newPageData, transparentnostView)
	m.Get("/kontakt/", newPageData, kontaktView)
	m.Get("/zaradi/", newPageData, zaradiView)
	m.Get("/projekti/", newPageData, projektiView)

	// m.Post("/", binding.Bind(SignupForm{}), newPageData, signupView)
	m.Post("/pridruzi-se/", binding.Bind(HackerSignupForm{}), newPageData, volontirajPostView)
	m.Post("/kontakt/", binding.Bind(ContactForm{}), newPageData, kontaktViewPost)

	m.NotFound(view404)

	// m.Run()
	log.Println("Server is running...")
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", conf.Port), m)
}
