package main

import (
	"fmt"
	"time"

	"github.com/gerifield/mnb-qr-go/src/qr"
)

func main() {

	code := qr.NewPaymentSend()
	code.Amount = qr.AmountHUF(10)
	code.ValidUntil = qr.Date(time.Now().Add(2 * time.Hour))

	fmt.Println(code)
}
