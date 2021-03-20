package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Prices struct {
	BTC float64 `json:"BTC"`
	ETH float64 `json:"ETH"`
	LTC float64 `json:"LTC"`
	HRK float64 `json:"HRK"`
	USD float64 `json:"USD"`
	EUR float64 `json:"EUR"`
}

type PriceClient struct {
	Prices *Prices
}

func (pc *PriceClient) doRequest() (*Prices, error) {
	p := &Prices{}
	cl := http.Client{}

	var req *http.Request
	var err error

	req, err = http.NewRequest(http.MethodGet, PricesURL, nil)

	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}

	res, err := cl.Do(req)

	if err == nil {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		if res.StatusCode != 200 {
			log.Println(string(body))
			err := errors.New(res.Status)
			return nil, err
		}
		json.Unmarshal(body, p)
	} else {
		return nil, err
	}

	return p, nil
}

func (pc *PriceClient) start() {
	go func() {
		for {
			if p, err := pc.doRequest(); err != nil {
				log.Println(err.Error())
			} else {
				pc.Prices = p
			}

			if conf.Debug {
				log.Printf("%#v\n", pc.Prices)
			}

			time.Sleep(time.Minute * 15)
		}
	}()
}

func initPriceClient() *PriceClient {
	pc := &PriceClient{}
	pc.start()
	return pc
}
