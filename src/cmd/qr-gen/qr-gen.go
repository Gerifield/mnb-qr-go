package main

import (
	"fmt"
	"os"
	"os/exec"
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
	png, err := code.GeneratePNG()
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
