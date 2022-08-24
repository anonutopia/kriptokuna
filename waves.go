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
			logTelegram(err.Error())
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
		talr.Timestamp >= wm.StartedTime &&
		talr.Sender != conf.Address &&
		talr.Recipient == conf.Address &&
		len(talr.AssetID) == 0 &&
		len(attachment) == 0 {

		wm.purchaseAsset(talr)
	} else if talr.Type == 4 &&
		talr.Timestamp >= wm.StartedTime &&
		talr.Sender != conf.Address &&
		talr.Recipient == conf.Address &&
		talr.AssetID == tokenID() &&
		len(attachment) == 0 {

		wm.sellAsset(talr)
	} else if talr.Type == 4 &&
		talr.Timestamp >= wm.StartedTime &&
		talr.Sender != conf.Address &&
		talr.Recipient == conf.Address &&
		talr.AssetID == tokenID() &&
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
		if conf.Address == AHRKAddress {
			amount := uint64((float64(amountWaves) / float64(SatInBTC)) * pc.Prices.EUR * 7.5 * float64(AHRKDec))
			sendAsset(amount, tokenID(), talr.Sender)
			messageTelegram(fmt.Sprintf("Promjena u kriptokunu: %.6f AHRK", float64(amount)/float64(AHRKDec)), TelAnonTeam)
			messageTelegram(fmt.Sprintf("Promjena u kriptokunu: %.6f AHRK", float64(amount)/float64(AHRKDec)), TelKriptokuna)
		} else if conf.Address == AEURAddress {
			amount := uint64((float64(amountWaves) / float64(SatInBTC)) * pc.Prices.USD * float64(AHRKDec))
			sendAsset(amount, tokenID(), talr.Sender)
			messageTelegram(fmt.Sprintf("Exchanged to AnonEuro: %.6f AEUR", float64(amount)/float64(AHRKDec)), TelAnonTeam)
			messageTelegram(fmt.Sprintf("Exchanged to AnonEuro: %.6f AEUR", float64(amount)/float64(AHRKDec)), TelAnonEuro)
		}
	}
}

func (wm *WavesMonitor) sellAsset(talr *gowaves.TransactionsAddressLimitResponse) {
	// Take fee
	amountHRK := talr.Amount - AHRKFee
	if amountHRK > 0 {
		if conf.Address == AHRKAddress {
			amount := uint64((float64(amountHRK) / float64(AHRKDec)) / (pc.Prices.EUR * 7.5 * float64(SatInBTC)))
			sendAsset(amount, "", talr.Sender)
			messageTelegram(fmt.Sprintf("Promjena iz kriptokune: %.8f WAVES\nAdresa: %s", float64(amount)/float64(SatInBTC), talr.Sender), TelAnonTeam)
		} else if conf.Address == AEURAddress {
			amount := uint64((float64(amountHRK) / float64(AHRKDec)) / (pc.Prices.USD) * float64(SatInBTC))
			sendAsset(amount, "", talr.Sender)
			messageTelegram(fmt.Sprintf("Exchanged from AnonEuro: %.8f WAVES\nAddress: %s", float64(amount)/float64(SatInBTC), talr.Sender), TelAnonTeam)
			messageTelegram(fmt.Sprintf("Exchanged from AnonEuro: %.8f WAVES\nAddress: %s", float64(amount)/float64(SatInBTC), talr.Sender), TelAnonOps)
		}
	}
}

func (wm *WavesMonitor) collectInterest(talr *gowaves.TransactionsAddressLimitResponse) {
	u := &User{Address: talr.Sender}
	db.First(u, u)
	if u.ID != 0 {
		var amount uint
		if conf.Address == AHRKAddress {
			amount = u.AmountAhrk
		} else if conf.Address == AEURAddress {
			amount = u.AmountAeur
		}
		sendAsset(uint64(amount), tokenID(), talr.Sender)
		u.AmountAhrk = 0
		u.AmountAeur = 0
		db.Save(u)
	}
}

func (wm *WavesMonitor) checkPayouts() {
	ks := &KeyValue{Key: "lastPayoutDay"}
	db.FirstOrCreate(ks, ks)

	if ks.ValueInt != uint64(time.Now().Day()) {
		ns, err := gowaves.WNC.NodeStatus()
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
			return
		}

		t, err := total(0, ns.BlockchainHeight-1, "")
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
			return
		}

		err = wm.doPayouts(ns.BlockchainHeight-1, "", t)
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		} else {
			ks.ValueInt = uint64(time.Now().Day())
			db.Save(ks)
		}
	}
}

func (wm *WavesMonitor) doPayouts(height int, after string, total int) error {
	abdr, err := gowaves.WNC.AssetsBalanceDistribution(tokenID(), height, 100, after)
	if err != nil {
		return err
	}

	for a, v := range abdr.Items {
		if !exclude(conf.Exclude, a) {
			ratio := getDailyRatio()
			amount := int(float64(v) * ratio)

			if amount > 0 {
				u := &User{Address: a}
				db.FirstOrCreate(u, u)
				if conf.Address == AHRKAddress {
					u.AmountAhrk += uint(amount)
				} else if conf.Address == AEURAddress {
					u.AmountAeur += uint(amount)
				}
				db.Save(u)
				log.Printf("Added interest: %s - %.6f", u.Address, float64(amount)/float64(AHRKDec))

				r := &User{}
				if u.ReferralID != 0 {
					db.First(r, u.ReferralID)
					ramount := int(float64(amount) * 0.2)
					if conf.Address == AHRKAddress {
						r.AmountAhrk += uint(ramount)
					} else if conf.Address == AEURAddress {
						r.AmountAeur += uint(ramount)
					}
					db.Save(r)
					log.Printf("Added referral interest: %s - %.6f", r.Address, float64(ramount)/float64(AHRKDec))
				}
			}
		}
	}

	if abdr.HasNext {
		return wm.doPayouts(height, abdr.LastItem, total)
	}

	return nil
}

func initWavesMonitor() {
	wm = &WavesMonitor{}
	wm.start()
}
