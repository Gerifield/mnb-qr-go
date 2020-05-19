package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gerifield/mnb-qr-go/src/qr"
)

func main() {

	code, err := qr.NewPaymentSend("GIBAHUHBXXX", "Test User", "HU00123456789012345678901234")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_ = code.HUFAmount(10)
	code.ValidUntil(time.Now().Add(2 * time.Hour))

	fmt.Println(code.String())
}
