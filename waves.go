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
	if talr.Type == 4 &&
		talr.Timestamp >= wm.StartedTime &&
		talr.Sender != AHRKAddress &&
		talr.Recipient == AHRKAddress &&
		len(talr.AssetID) == 0 {
			wm.purchaseAsset(talr)
	} else if talr.Type == 4 &&
		talr.Timestamp >= wm.StartedTime &&
		talr.Sender != AHRKAddress &&
		talr.Recipient == AHRKAddress &&
		talr.AssetID == AHRKId {
			wm.sellAsset(talr)
	}

	tr.Processed = true
	db.Save(tr)
}

func (wm *WavesMonitor) purchaseAsset(talr *gowaves.TransactionsAddressLimitResponse)  {
	// Take fee
	amountWaves := talr.Amount - 100000
	if (amountWaves > 0) {
		amount := uint64((float64(amountWaves) / float64(SatInBTC)) * pc.Prices.HRK * float64(AHRKDec))
		sendAsset(amount, AHRKId, talr.Sender)
		log.Println("done")
	}
}

func (wm *WavesMonitor) sellAsset(talr *gowaves.TransactionsAddressLimitResponse)  {
	log.Println(talr)
}

func initWavesMonitor() {
	wm = &WavesMonitor{}
	wm.start()
}
