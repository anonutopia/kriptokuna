package main

import (
	"fmt"
	"log"
	"time"

	"github.com/anonutopia/gowaves"
	"github.com/wavesplatform/gowaves/pkg/crypto"
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

		wm.checkPayouts()

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
	attachment := ""

	if len(talr.Attachment) > 0 {
		attachment = string(crypto.MustBytesFromBase58(talr.Attachment))
	}

	if talr.Type == 4 &&
		// talr.Timestamp >= wm.StartedTime &&
		talr.Sender != AHRKAddress &&
		talr.Recipient == AHRKAddress &&
		len(talr.AssetID) == 0 &&
		len(attachment) == 0 {

		wm.purchaseAsset(talr)
	} else if talr.Type == 4 &&
		talr.Timestamp >= wm.StartedTime &&
		talr.Sender != AHRKAddress &&
		talr.Recipient == AHRKAddress &&
		talr.AssetID == AHRKId &&
		len(attachment) == 0 {

		wm.sellAsset(talr)
	} else if talr.Type == 4 &&
		talr.Timestamp >= wm.StartedTime &&
		talr.Sender != AHRKAddress &&
		talr.Recipient == AHRKAddress &&
		talr.AssetID == AHRKId &&
		attachment == "collect" &&
		talr.Amount == 950000 {

		wm.collectInterest(talr)
	}

	tr.Processed = true
	db.Save(tr)
}

func (wm *WavesMonitor) purchaseAsset(talr *gowaves.TransactionsAddressLimitResponse) {
	// Take fee
	amountWaves := talr.Amount - WavesFee
	if amountWaves > 0 {
		amount := uint64((float64(amountWaves) / float64(SatInBTC)) * pc.Prices.USD * HRKUSD * float64(AHRKDec))
		sendAsset(amount, AHRKId, talr.Sender)
		messageTelegram(fmt.Sprintf("Promjena u kriptokunu: %.6f AHRK", float64(amount)/float64(AHRKDec)), TelAnonTeam)
	}
}

func (wm *WavesMonitor) sellAsset(talr *gowaves.TransactionsAddressLimitResponse) {
	// Take fee
	amountHRK := talr.Amount - AHRKFee
	if amountHRK > 0 {
		amount := uint64((float64(amountHRK) / float64(AHRKDec)) / pc.Prices.USD * HRKUSD * float64(SatInBTC))
		sendAsset(amount, "", talr.Sender)
		messageTelegram(fmt.Sprintf("Promjena iz kriptokune: %.8f WAVES\nAdresa: %s", float64(amount)/float64(SatInBTC), talr.Sender), TelAnonTeam)
	}
}

func (wm *WavesMonitor) collectInterest(talr *gowaves.TransactionsAddressLimitResponse) {
	u := &User{Address: talr.Sender}
	db.First(u, u)
	if u.ID != 0 {
		sendAsset(uint64(u.AmountAhrk), AHRKId, talr.Sender)
		u.AmountAhrk = 0
		db.Save(u)
	}
}

func (wm *WavesMonitor) checkPayouts() {
	ks := &KeyValue{Key: "lastPayoutDay"}
	db.FirstOrCreate(ks, ks)

	if ks.ValueInt != uint64(time.Now().Day()) {
		newValue := 0
		ns, err := gowaves.WNC.NodeStatus()
		if err != nil {
			log.Println(err)
			return
		}

		t, err := total(0, ns.BlockchainHeight-1, "")
		if err != nil {
			log.Println(err)
			return
		}

		talr, err := gowaves.WNC.TransactionsAddressLimit(AHRKAddress, 100)
		if err != nil {
			log.Println(err)
			return
		}

		for _, t := range talr[0] {
			tm := time.Unix(t.Timestamp/1000, 0)
			if t.Type == 11 && tm.Day() == time.Now().Day() {
				newValue = t.Transfers[0].Amount
				break
			}
		}

		newValueHRK := int((float64(newValue) / (float64(pc.Prices.JPY / pc.Prices.HRK))))
		newValueRatio := float64(newValueHRK) / float64(t)
		var extraValue int

		if newValueRatio > getDailyRatio(1.1) {
			extraValue := newValueHRK
			newValueHRK = int(float64(newValueHRK) * getDailyRatio(1.1))
			extraValue -= newValueHRK
		}

		if newValueHRK > 0 {
			err = wm.doPayouts(ns.BlockchainHeight-1, "", t, newValueHRK)
			if err != nil {
				log.Println(err)
			} else {
				ks.ValueInt = uint64(time.Now().Day())
				db.Save(ks)
			}
		}

		if extraValue > 0 {
			log.Println("There's extra value.")
		}
	}
}

func (wm *WavesMonitor) doPayouts(height int, after string, total int, newValueHRK int) error {
	abdr, err := gowaves.WNC.AssetsBalanceDistribution(AHRKId, height, 100, after)
	if err != nil {
		return err
	}

	for a, v := range abdr.Items {
		if !exclude(conf.Exclude, a) {
			ratio := float64(v) / float64(total)
			amount := int(float64(newValueHRK) * ratio)

			if amount > 0 {
				u := &User{Address: a}
				db.FirstOrCreate(u, u)
				u.AmountAhrk += uint(amount)
				db.Save(u)
				log.Printf("Added interest: %s - %.6f", u.Address, float64(amount)/float64(AHRKDec))
			}
		}
	}

	if abdr.HasNext {
		return wm.doPayouts(height, abdr.LastItem, total, newValueHRK)
	}

	return nil
}

func initWavesMonitor() {
	wm = &WavesMonitor{}
	wm.start()
}
