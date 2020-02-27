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

		wm.checkPayouts()

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

func (wm *WavesMonitor) checkPayouts() {
	ks := &KeyValue{Key: "lastPayoutDay"}
	db.FirstOrCreate(ks, ks)

	if ks.ValueInt != uint64(time.Now().Day()) {
		ns, err := wnc.NodeStatus()
		if err != nil {
			logTelegram(fmt.Sprintf("wnc.NodeStatus error: %e", err))
			return
		}

		t, err := total(ns.BlockchainHeight-1, "")
		if err != nil {
			logTelegram(fmt.Sprintf("total error: %e", err))
			return
		}

		new := 0

		if new > 0 {
			err = wm.doPayouts(ns.BlockchainHeight-1, "", t, new)
			if err != nil {
				logTelegram(fmt.Sprintf("wm.doPayouts error: %e", err))
			}
		}

		ks.ValueInt = uint64(time.Now().Day())
		db.Save(ks)
	}
}

func (wm *WavesMonitor) doPayouts(height int, after string, total int, new int) error {
	abdr, err := wnc.AssetsBalanceDistribution(conf.TokenID, height, 100, after)
	if err != nil {
		return err
	}

	for a, v := range abdr.Items {
		if !exclude(conf.Exclude, a) {
			ratio := float64(v) / float64(total)
			amount := int(float64(new) * ratio)

			if amount > 0 {
				atr := &gowaves.AssetsTransferRequest{
					Amount:     amount,
					AssetID:    conf.TokenID,
					FeeAssetID: conf.TokenID,
					Fee:        5,
					Recipient:  a,
					Sender:     conf.NodeAddress,
				}

				_, err = wnc.AssetsTransfer(atr)
				if err != nil {
					log.Printf("[doPayouts] error assets transfer: %s", err)
					logTelegram(fmt.Sprintf("[doPayouts] error assets transfer: %s", err))
				} else {
					log.Printf("Sent token: %s => %d", a, amount)
				}
			}
		}
	}

	if abdr.HasNext {
		return wm.doPayouts(height, abdr.LastItem, total, new)
	}

	return nil
}

func exclude(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func total(height int, after string) (int, error) {
	t := 0

	abdr, err := wnc.AssetsBalanceDistribution(conf.TokenID, height, 100, after)
	if err != nil {
		return 0, err
	}

	for a, v := range abdr.Items {
		if !exclude(conf.Exclude, a) {
			t = t + v
		}
	}

	if abdr.HasNext {
		return total(height, abdr.LastItem)
	}

	return t, nil
}

func initMonitor() {
	wm := &WavesMonitor{}
	wm.start()
}
