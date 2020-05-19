package qr

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
)

type Code struct {
	Kind       kind    // Required
	Version    version // Required
	Charset    int     // Required
	BIC        string  // Required
	Name       string  // Required
	IBAN       string  // Required
	Amount     amount
	Valid      date // Required
	Purpose    [4]byte
	Message    [70]byte
	ShopID     [35]byte
	MerchDevID [35]byte
	InvoiceID  [35]byte
	CustomerID [35]byte
	CredTranID [35]byte
	LoyaltyID  [35]byte
	NavCheckID [35]byte
	//SeparatorLength [17]byte // Required, placeholder
}

var (
	// KindHCT for send money
	KindHCT kind = "HCT"

	// KindRTP for request money
	KindRTP kind = "RTP"
)

// kind QR Code type
type kind string

// String .
func (k kind) String() string {
	return string(k)
}

// version of the QR code
type version string

// String .
func (v version) String() string {
	if v == "" {
		return "001" // Default
	}
	return string(v)
}

// Amount for payment (optional)
type amount struct {
	currency string
	total    int
}

// String .
func (a amount) String() string {
	currency := a.currency
	if currency == "" {
		currency = "HUF"
	}
	return fmt.Sprintf("%s%d", currency, a.total)
}

func (c Code) GeneratePNG() ([]byte, error) {
	if time.Now().After(time.Time(c.Valid)) {
		return nil, errors.New("negative validity period")
	}
	q, err := qrcode.New(c.String(), qrcode.Medium)
	if err != nil {
		return nil, err
	}

	if q.VersionNumber > 13 { // TODO: add testing for this part
		return nil, errors.New("generated image is too big")
	}
	return q.PNG(128)
}

// String .
func (c Code) String() string {
	var sb strings.Builder

	sb.WriteString(c.Kind.String())
	sb.WriteString("\n")

	sb.WriteString(c.Version.String())
	sb.WriteString("\n")

	if c.Charset == 0 {
		sb.WriteString("1") // Set default to 1
	} else {
		sb.WriteString(fmt.Sprintf("%d", c.Charset))
	}
	sb.WriteString("\n")

	sb.WriteString(c.BIC)
	sb.WriteString("\n")

	sb.WriteString(c.Name)
	sb.WriteString("\n")

	sb.WriteString(c.IBAN)
	sb.WriteString("\n")

	// Optional amount
	if c.Amount.total > 0 {
		sb.WriteString(c.Amount.String())
	}
	sb.WriteString("\n")

	// TODO: This should be valid and be here
	sb.WriteString(c.Valid.String())
	sb.WriteString("\n")

	// TODO: Fill these optional fields

	//c.Purpose // AT-44
	sb.WriteString("\n")

	//c.Message
	sb.WriteString("\n")

	//c.ShopID
	sb.WriteString("\n")

	//c.MerchDevID
	sb.WriteString("\n")

	//c.InvoiceID
	sb.WriteString("\n")

	//c.CustomerID
	sb.WriteString("\n")

	//c.CredTranID
	sb.WriteString("\n")

	//c.LoyaltyID
	sb.WriteString("\n")

	//c.NavCheckID
	sb.WriteString("\n")

	return sb.String()
}

// HUFAmount for the transaction
func (c *Code) HUFAmount(total int) error {
	if total < 0 {
		return errors.New("amount could not be negative")
	}

	if total > 999999999999 {
		return errors.New("amount could not be higher than 999999999999")
	}

	c.Amount = amount{
		currency: "HUF",
		total:    total,
	}
	return nil
}

// ValidUntil .
func (c *Code) ValidUntil(t time.Time) {
	c.Valid = date(t)
}

// NewPaymentSend QR code creation
// The reader of the QR code will send the payment to the generator user.
// In Hungarian: Ez az átutalási megbízás, azaz a kedvezményezett generálja a QR
// kódot, hogy a fizető fél a megfelelő adatokkal tudja elküdeni az összeget.
// In Englis from the doc: If it supports the submission of the credit transfer order – i.e. the payee generates the QR code to
// enable the payer to submit the credit transfer order with the correct data – the “HCT” code must be used.
func NewPaymentSend(bic string, name string, iban string) (*Code, error) {
	c := &Code{
		Kind: KindHCT,
	}

	if err := addRecipient(c, bic, name, iban); err != nil {
		return nil, err
	}
	return c, nil
}

// NewPaymentRequest QR code creation
// The reader of the QR code will send a payment request to the generator user
// In Hungarian: Ez a fizetési kérelem küldése, azaz a fizető fél adja meg a QR-kód generálásával
// a főbb adatait a kedvezményezettnek, hogy az utóbbi fizetési kérelmet tudjon küldeni.
// In Englis from the doc: If it supports the transmission of the request to pay – i.e. the payer generates the QR code to transfer
// his main data to the payee in order to enable the payee to send a request to pay – the RTP code must be used.
func NewPaymentRequest(bic string, name string, iban string) (*Code, error) {
	c := &Code{
		Kind: KindRTP,
	}

	if err := addRecipient(c, bic, name, iban); err != nil {
		return nil, err
	}
	return c, nil
}

func addRecipient(code *Code, bic, name, iban string) error {
	if len(bic) == 8 {
		bic = bic + "XXX" // For SEPA payment the 8 char long SWIFT should be extended with XXX to 11 chars
	}

	if len(bic) != 11 {
		return errors.New("invalid BIC length")
	}
	code.BIC = bic

	if len(name) > 70 {
		return errors.New("name should not be longer than 70")
	}
	code.Name = name

	if len(iban) != 28 {
		return errors.New("invalid IBAN length")
	}
	code.IBAN = iban
	return nil
}
