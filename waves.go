package main

import (
	"log"
	"time"

	"github.com/anonutopia/gowaves"
)

type WavesMonitor struct {
	StartedTime int64
}

func (wm *WavesMonitor) start() {
	wm.StartedTime = time.Now().Unix() * 1000
	for {
		// todo - make sure that everything is ok with 100 here
		pages, err := gowaves.WNC.TransactionsAddressLimit(conf.Address, 100)
		if err != nil {
			log.Println(err)
		}

		if len(pages) > 0 {
			for _, t := range pages[0] {
				wm.checkTransaction(&t)
			}
		}

		time.Sleep(time.Second * WavesMonitorTick)
	}
}

func (wm *WavesMonitor) checkTransaction(talr *gowaves.TransactionsAddressLimitResponse) {
	tr := &Transaction{TxID: talr.ID}
	db.FirstOrCreate(tr, tr)
	if !tr.Processed {
		wm.processTransaction(tr, talr)
	}
}

func (wm *WavesMonitor) processTransaction(tr *Transaction, talr *gowaves.TransactionsAddressLimitResponse) {
	tr.Processed = true
	db.Save(tr)
}

func initWavesMonitor() {
	wm = &WavesMonitor{}
	wm.start()
}
