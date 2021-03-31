package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/anonutopia/gowaves"
	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

func sendAsset(amount uint64, assetId string, recipient string) error {
	if conf.Dev {
		return errors.New(fmt.Sprintf("Not sending (dev): %d - %s - %s", amount, assetId, recipient))
	}

	var assetBytes []byte

	// Create sender's public key from BASE58 string
	sender, err := crypto.NewPublicKeyFromBase58(conf.PublicKey)
	if err != nil {
		log.Println(err)
		return err
	}

	// Create sender's private key from BASE58 string
	sk, err := crypto.NewSecretKeyFromBase58(conf.PrivateKey)
	if err != nil {
		log.Println(err)
		return err
	}

	// Current time in milliseconds
	ts := time.Now().Unix() * 1000

	if len(assetId) > 0 {
		assetBytes = crypto.MustBytesFromBase58(assetId)
	} else {
		assetBytes = []byte{}
	}

	asset, err := proto.NewOptionalAssetFromBytes(assetBytes)
	if err != nil {
		log.Println(err)
		return err
	}

	rec, err := proto.NewAddressFromString(recipient)
	if err != nil {
		log.Println(err)
		return err
	}

	tr := proto.NewUnsignedTransferWithSig(sender, *asset, *asset, uint64(ts), amount, 100000, proto.Recipient{Address: &rec}, nil)

	err = tr.Sign(proto.MainNetScheme, sk)
	if err != nil {
		log.Println(err)
		return err
	}

	// Create new HTTP client to send the transaction to public TestNet nodes
	client, err := client.NewClient(client.Options{BaseUrl: WavesNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		return err
	}

	// Context to cancel the request execution on timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// // Send the transaction to the network
	_, err = client.Transactions.Broadcast(ctx, tr)
	if err != nil {
		log.Println(err)
		return err
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

func total(t int, height int, after string) (int, error) {
	abdr, err := gowaves.WNC.AssetsBalanceDistribution(AHRKId, height, 100, after)
	if err != nil {
		return 0, err
	}

	for a, v := range abdr.Items {
		if !exclude(conf.Exclude, a) {
			t = t + v
		}
	}

	if abdr.HasNext {
		return total(t, height, abdr.LastItem)
	}

	return t, nil
}
