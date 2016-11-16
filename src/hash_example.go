package main

import (
	"fmt"
	"strconv"
)

func main() {
	var h Header

	/*Preencher estrutura*/
	h.ZeroCount = 20
	h.Date = "031116"
	h.Resource = "amadilsons"
	h.RandString = Encode64([]byte(RndStr(6)))
	h.Counter = RndInt(1000)

	digest := HashDigest(h.HeaderToString())
	for !CheckZeroBits(digest, h.ZeroCount) {
		h.IncCounter()
		digest = HashDigest(h.HeaderToString())
	}
	fmt.Println(h.HeaderToString())

	var h1 Header
	finalHeader := h.HeaderToString()
	(&h1).StringToHeader(finalHeader)
	fmt.Println(h1, "  ", Encode64([]byte(strconv.Itoa(h.Counter))))
}
