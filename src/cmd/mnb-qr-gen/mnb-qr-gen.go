package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gerifield/mnb-qr-go/src/qr"
)

func main() {
	qrType := flag.String("type", "RTP", "QR code type (RTP/HCT)")
	bic := flag.String("bic", "", "BIC code")
	name := flag.String("name", "", "Name")
	iban := flag.String("iban", "", "IBAN number")
	amount := flag.Int("amount", 0, "Amount to request (in HUF)")
	message := flag.String("message", "", "Message in the QR code")
	flag.Parse()

	qrt := strings.ToUpper(*qrType)
	if qrt != "RTP" && qrt != "HCT" {
		fmt.Println("Invalid QR code type (shoulb be RTP or HCT)")
		os.Exit(1)
	}

	var err error
	var code *qr.Code
	if qrt == "HCT" {
		code, err = qr.NewPaymentSend(*bic, *name, *iban)
	} else {
		code, err = qr.NewPaymentRequest(*bic, *name, *iban)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = code.HUFAmount(*amount)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = code.Message(*message)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_ = code.ValidUntil(time.Now().Add(2 * time.Hour))

	fmt.Println(code.String())
	png, err := code.GeneratePNG(256)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	f, err := os.Create("out.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	_, _ = f.Write(png)

	// Open the image
	c := exec.Command("open", "out.png")
	err = c.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
