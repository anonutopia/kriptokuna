package main

import (
	"fmt"

	macaron "gopkg.in/macaron.v1"
)

func accumulatedInterest(ctx *macaron.Context) string {
	ctx.Resp.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	address := ctx.Params("address")

	u := &User{Address: address}
	db.FirstOrCreate(u, u)

	interest := float64(u.Accumulation) / float64(AHRKDec)
	response := fmt.Sprintf("document.onload(document.getElementById('accumulatedInerest').innerHTML = '%.6f');", interest)
	return response
}
