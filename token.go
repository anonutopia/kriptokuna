package main

import "log"

const priceFactorLimit = uint64(0.0001 * float64(satInBtc))

// Token represents issuing token
type Token struct {
	Price           uint64
	PriceFactor     uint64
	TierPrice       uint64
	TierPriceFactor uint64
}

func (t *Token) issueAmount(investment int, assetID string) int {
	amount := int(0)

	p, err := pc.DoRequest()
	if err != nil {
		log.Printf("[token.issueAmount] error pc.DoRequest: %s", err)
		return 0
	}

	var cryptoPrice float64
	var investmentEur float64

	if len(assetID) == 0 {
		cryptoPrice = p.WAVES
	} else {
		return amount
	}

	priceChanged := false
	investmentEur = float64(investment) / float64(satInBtc) / cryptoPrice
	sendGroupsMessageInvestment(investmentEur)
	log.Println(investmentEur)

	for investment > 10 {
		log.Printf("token: %d %d %d %d", t.Price, t.PriceFactor, t.TierPrice, t.TierPriceFactor)
		log.Printf("cryptoPrice: %f", cryptoPrice)
		log.Printf("investment: %d", investment)

		tierAmount := uint64(float64(investment) / cryptoPrice / float64(t.Price) * float64(satInBtc))

		log.Printf("tierAmount: %d", tierAmount)

		if tierAmount > t.TierPrice {
			tierAmount = t.TierPrice
		}

		log.Printf("tierAmount: %d", tierAmount)

		tierInvestment := int(float64(tierAmount) * float64(t.Price) * cryptoPrice / float64(satInBtc))

		log.Printf("tierInvestment: %d", tierInvestment)

		amount = amount + int(tierAmount)

		log.Printf("amount: %d", amount)

		investment = investment - tierInvestment

		log.Printf("investment: %d", investment)

		t.TierPrice = t.TierPrice - tierAmount
		t.TierPriceFactor = t.TierPriceFactor - tierAmount

		log.Printf("token: %d %d %d %d", t.Price, t.PriceFactor, t.TierPrice, t.TierPriceFactor)

		if t.TierPrice == 0 {
			t.TierPrice = 1000 * satInBtc
			t.Price = t.Price + t.PriceFactor
			priceChanged = true
		}

		if t.TierPriceFactor == 0 {
			t.TierPriceFactor = 1000000 * satInBtc
			if t.PriceFactor > priceFactorLimit {
				t.PriceFactor = t.PriceFactor - priceFactorLimit
			}
		}

		t.saveState()

		log.Printf("token: %d %d %d %d", t.Price, t.PriceFactor, t.TierPrice, t.TierPriceFactor)
	}

	if priceChanged {
		newPrice := float64(t.Price) / float64(satInBtc)
		sendGroupsMessagePrice(newPrice)
		log.Println(newPrice)
	}

	return amount
}

func (t *Token) saveState() {
	ksip := &KeyValue{Key: "tokenPrice"}
	db.FirstOrCreate(ksip, ksip)
	ksip.ValueInt = t.Price
	db.Save(ksip)

	ksipf := &KeyValue{Key: "tokenPriceFactor"}
	db.FirstOrCreate(ksipf, ksipf)
	ksipf.ValueInt = t.PriceFactor
	db.Save(ksipf)

	ksitp := &KeyValue{Key: "tokenTierPrice"}
	db.FirstOrCreate(ksitp, ksitp)
	ksitp.ValueInt = t.TierPrice
	db.Save(ksitp)

	ksitpf := &KeyValue{Key: "tokenTierPriceFactor"}
	db.FirstOrCreate(ksitpf, ksitpf)
	ksitpf.ValueInt = t.TierPriceFactor
	db.Save(ksitpf)
}

func (t *Token) loadState() {
	ksip := &KeyValue{Key: "tokenPrice"}
	db.FirstOrCreate(ksip, ksip)

	if ksip.ValueInt > 0 {
		t.Price = ksip.ValueInt
	} else {
		ksip.ValueInt = t.Price
		db.Save(ksip)
	}

	ksipf := &KeyValue{Key: "tokenPriceFactor"}
	db.FirstOrCreate(ksipf, ksipf)

	if ksipf.ValueInt > 0 {
		t.PriceFactor = ksipf.ValueInt
	} else {
		ksipf.ValueInt = t.PriceFactor
		db.Save(ksipf)
	}

	ksitp := &KeyValue{Key: "tokenTierPrice"}
	db.FirstOrCreate(ksitp, ksitp)

	if ksitp.ValueInt > 0 {
		t.TierPrice = ksitp.ValueInt
	} else {
		ksitp.ValueInt = t.TierPrice
		db.Save(ksitp)
	}

	ksitpf := &KeyValue{Key: "tokenTierPriceFactor"}
	db.FirstOrCreate(ksitpf, ksitpf)

	if ksitpf.ValueInt > 0 {
		t.TierPriceFactor = ksitpf.ValueInt
	} else {
		ksitpf.ValueInt = t.TierPriceFactor
		db.Save(ksitpf)
	}
}

func initToken() *Token {
	token := &Token{
		Price:           uint64(0.01 * float64(satInBtc)),
		PriceFactor:     uint64(0.0021 * float64(satInBtc)),
		TierPrice:       1000 * satInBtc,
		TierPriceFactor: 1000000 * satInBtc}

	token.loadState()

	return token
}
