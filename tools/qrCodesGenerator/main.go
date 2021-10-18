package main

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/skip2/go-qrcode"
)

type qrCodeJson struct {
	Uuid string `json:"uuid"`
	To   string `json:"to"`
}

type qrCodesJson []qrCodeJson

func main() {
	qrCodesBytes, err := os.ReadFile("../../data/qr.json")
	if err != nil {
		panic(err)
	}

	qrCodes := qrCodesJson{}
	err = json.Unmarshal(qrCodesBytes, &qrCodes)
	if err != nil {
		panic(err)
	}

	err = os.RemoveAll("./codes")
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll("./codes", os.ModePerm)
	if err != nil {
		panic(err)
	}

	for i, qrCode := range qrCodes {
		image, err := encode("https://api.terrainstory.com/qr/", qrCode.Uuid)
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("./codes/"+strconv.Itoa(i+1)+".png", image, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

func encode(base, uuid string) ([]byte, error) {
	return qrcode.Encode(base+uuid, qrcode.Highest, 3000)
}
