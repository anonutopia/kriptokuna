package main

import (
	"reflect"

	"github.com/go-macaron/captcha"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

func homeView(ctx *macaron.Context) {
	ctx.Data["Title"] = ""

	sup := &SignupForm{Type: 4}
	ctx.Data["Form"] = sup

	ctx.HTML(200, "home")
}

func kriptokunaView(ctx *macaron.Context) {
	ctx.Data["Title"] = "Što je Kriptokuna? | "

	ctx.HTML(200, "kriptokuna")
}

func manifestView(ctx *macaron.Context) {
	ctx.Data["Title"] = "Manifest Blokiranima | "

	ctx.HTML(200, "blokirani")
}

func volontirajView(ctx *macaron.Context) {
	ctx.Data["Title"] = "Volontiraj već danas! | "

	hsup := &HackerSignupForm{Type: "avatari"}
	ctx.Data["Form"] = hsup

	ctx.HTML(200, "volontiraj")
}

func volontirajPostView(ctx *macaron.Context, hsup HackerSignupForm, f *session.Flash, cpt *captcha.Captcha) {
	ctx.Data["Title"] = ""

	s := reflect.ValueOf(ctx.Data["Errors"])

	if cpt.VerifyReq(ctx.Req) {
		if s.Len() == 0 {
			ha := &Hacker{Email: hsup.Email, Type: hsup.Type}
			db.FirstOrCreate(ha, ha)
			f.Success("Dodan/a u našu Core Team bazu, javit ćemo se uskoro.")
			f.Warning("Za Kriptokunu reci svim korisnicima Ubera (vozačima i putnicima) i svim blokiranim građanima koje poznaješ!")
			ctx.Redirect("/volontiraj/#hackersignup")
			return
		} else {
			f.Error("Email je obavezno polje.")
			ctx.Data["Flash"] = f
			ctx.Data["Form"] = hsup
		}
	} else {
		f.Error("Pogrešan captcha broj. Pokušajte ponovo!")
		ctx.Data["Flash"] = f
		ctx.Data["Form"] = hsup
	}

	ctx.HTML(200, "volontiraj")
}

func faqView(ctx *macaron.Context) {
	ctx.Data["Title"] = "Često postavljena pitanja | "

	ctx.HTML(200, "faq")
}

func signupView(ctx *macaron.Context, sup SignupForm, f *session.Flash, cpt *captcha.Captcha) {
	ctx.Data["Title"] = ""

	s := reflect.ValueOf(ctx.Data["Errors"])

	if cpt.VerifyReq(ctx.Req) {
		if s.Len() == 0 {
			ha := &Hacktivist{Email: sup.Email, Type: sup.Type}
			db.FirstOrCreate(ha, ha)
			f.Success("Uspješno si dodan/a u našu bazu hacktivista, javit ćemo se uskoro.")
			f.Warning("Za Kriptokunu reci svim korisnicima Ubera (vozačima i putnicima) i svim blokiranim građanima koje poznaješ!")
			ctx.Redirect("/#signup")
			return
		} else {
			f.Error("Email je obavezno polje.")
			ctx.Data["Flash"] = f
			ctx.Data["Form"] = sup
		}
	} else {
		f.Error("Pogrešan captcha broj. Pokušajte ponovo!")
		ctx.Data["Flash"] = f
		ctx.Data["Form"] = sup
	}

	ctx.HTML(200, "home")
}

func view404(ctx *macaron.Context) {
	ctx.Data["URI"] = "/not-found/"
	ctx.Data["Title"] = "404 | "

	ctx.HTML(200, "404")
}
