package main

import (
	"github.com/anonutopia/gowaves"
)

func initWaves() *gowaves.WavesNodeClient {
	wnc := &gowaves.WavesNodeClient{
		Host:   conf.WavesNode,
		Port:   6869,
		ApiKey: conf.WavesNodeAPIKey,
	}

	return wnc
}
