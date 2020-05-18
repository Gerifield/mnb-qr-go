package qr

import (
	"fmt"
	"strings"
)

type Code struct {
	Kind            Kind     // Required
	Version         Version  // Required
	Charset         [1]byte  // Required
	BIC             [11]byte // Required
	Name            [70]byte // Required
	IBAN            [28]byte // Required
	Amount          Amount
	ValidUntil      Date // Required
	Purpose         [4]byte
	Message         [70]byte
	ShopID          [35]byte
	MerchDevID      [35]byte
	InvoiceID       [35]byte
	CustomerID      [35]byte
	CredTranID      [35]byte
	LoyaltyID       [35]byte
	NavCheckID      [35]byte
	SeparatorLength [17]byte // Required
}

var (
	// KindHCT for send money
	KindHCT Kind = "HCT"

	// KindRTP for request money
	KindRTP Kind = "RTP"
)

// Kind QR Code type
type Kind string

// String .
func (k Kind) String() string {
	return string(k)
}

// Version of the QR code
type Version string

// String .
func (v Version) String() string {
	if v == "" {
		return "001" // Default
	}
	return string(v)
}

// Amount for payment (optional)
type Amount struct {
	currency string
	total    int
}

// AmountHUF return a HUF value
func AmountHUF(total int) Amount {
	return Amount{
		currency: "HUF",
		total:    total,
	}
}

// String .
func (a Amount) String() string {
	return fmt.Sprintf("%s%d", a.currency, a.total)
}

// String .
func (c Code) String() string {
	var sb strings.Builder

	sb.WriteString(c.Kind.String())
	sb.WriteString("\n")

	sb.WriteString(c.Version.String())
	sb.WriteString("\n")

	return sb.String()
}

// NewPaymentSend QR code creation
func NewPaymentSend() Code {
	return Code{
		Kind: KindHCT,
		//Version: "001",
	}
}

// NewPaymentRequest QR code creation
func NewPaymentRequest() Code {
	return Code{
		Kind: KindRTP,
	}
}
