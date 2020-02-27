package main

import (
	"fmt"
	"log"
	"time"

	"github.com/anonutopia/gowaves"
)

const satInBtc = uint64(100000000)

// WavesMonitor represents waves monitoring object
type WavesMonitor struct {
	StartedTime int64
}

func (wm *WavesMonitor) start() {
	wm.StartedTime = time.Now().Unix() * 1000
	for {
		// todo - make sure that everything is ok with 100 here
		pages, err := wnc.TransactionsAddressLimit(conf.NodeAddress, 100)
		if err != nil {
			log.Println(err)
		}

		if len(pages) > 0 {
			for _, t := range pages[0] {
				wm.checkTransaction(&t)
			}
		}

		time.Sleep(time.Second)
	}
}

func (wm *WavesMonitor) checkTransaction(t *gowaves.TransactionsAddressLimitResponse) {
	tr := Transaction{TxID: t.ID}
	db.FirstOrCreate(&tr, &tr)
	if tr.Processed != true {
		wm.processTransaction(&tr, t)
	}
}

func (wm *WavesMonitor) processTransaction(tr *Transaction, t *gowaves.TransactionsAddressLimitResponse) {
	if t.Type == 4 &&
		t.Timestamp >= wm.StartedTime &&
		t.Sender != conf.NodeAddress &&
		t.Recipient == conf.NodeAddress &&
		len(t.AssetID) == 0 {

		wm.purchaseAsset(t)
	} else if t.Type == 4 &&
		t.Timestamp >= wm.StartedTime &&
		t.Sender != conf.NodeAddress &&
		t.Recipient == conf.NodeAddress &&
		t.AssetID == conf.TokenID {

		wm.sellAsset(t)
	}

	tr.Processed = true
	db.Save(tr)
}

func (wm *WavesMonitor) purchaseAsset(t *gowaves.TransactionsAddressLimitResponse) {
	prices, err := pc.DoRequest()
	if err != nil {
		logTelegram(fmt.Sprintf("pc.DoRequest error: %e", err))
		return
	}

	amount := int((float64(t.Amount) / float64(satInBtc)) / prices.WAVES * 100)

	if amount > 5 {
		amount = amount - 5
	}

	atr := &gowaves.AssetsTransferRequest{
		Amount:     amount,
		AssetID:    conf.TokenID,
		FeeAssetID: conf.TokenID,
		Fee:        5,
		Recipient:  t.Sender,
		Sender:     conf.NodeAddress,
	}

	_, err = wnc.AssetsTransfer(atr)
	if err != nil {
		log.Printf("[purchaseAsset] error assets transfer: %s", err)
		logTelegram(fmt.Sprintf("[purchaseAsset] error assets transfer: %s", err))
	} else {
		log.Printf("Sent token: %s => %d", t.Sender, amount)
	}
}

func (wm *WavesMonitor) sellAsset(t *gowaves.TransactionsAddressLimitResponse) {
	prices, err := pc.DoRequest()
	if err != nil {
		logTelegram(fmt.Sprintf("pc.DoRequest error: %e", err))
		return
	}

	amount := int((float64(t.Amount) / float64(100)) * prices.WAVES * float64(satInBtc))

	atr := &gowaves.AssetsTransferRequest{
		Amount:     amount,
		AssetID:    "",
		FeeAssetID: conf.TokenID,
		Fee:        5,
		Recipient:  t.Sender,
		Sender:     conf.NodeAddress,
	}

	_, err = wnc.AssetsTransfer(atr)
	if err != nil {
		log.Printf("[sellAsset] error assets transfer: %s", err)
		logTelegram(fmt.Sprintf("[sellAsset] error assets transfer: %s", err))
	} else {
		log.Printf("Sent WAVES: %s => %d", t.Sender, amount)
	}
}

func initMonitor() {
	wm := &WavesMonitor{}
	wm.start()
}
