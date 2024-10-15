package qr

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
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
	purpose    string
	message    string
	shopID     string
	merchDevID string
	invoiceID  string
	customerID string
	credTranID string
	loyaltyID  string
	navCheckID string
	//SeparatorLength [17]byte // Required, placeholder
}

const (
	qrContentMaxSize = 345
)

var (
	// KindHCT for send money
	KindHCT kind = "HCT"

	// KindRTP for request money
	KindRTP kind = "RTP"

	// See the details in the checks, not all codes here maybe
	purposeCodes = []string{"ACCT", "ADVA", "AGRT", "AIRB", "ALMY", "ANNI", "ANTS", "AREN", "BECH", "BENE", "BEXP", "BOCE", "BONU", "BUSB", "CASH", "CBFF", "CBTV", "CCRD", "CDBL", "CFEE", "CHAR", "CLPR", "CMDT", "COLL", "COMC", "COMM", "COMT", "COST", "CPYR", "CSDB", "CSLP", "CVCF", "DBTC", "DCRD", "DEPT", "DERI", "DIVD", "DMEQ", "DNTS", "ELEC", "ENRG", "ESTX", "FERB", "FREX", "GASB", "GDDS", "GDSV", "GOVI", "GOVT", "GSCB", "GVEA", "GVEB", "GVEC", "GVED", "HEDG", "HLRP", "HLTC", "HLTI", "HREC", "HSPC", "HSTX", "ICCP", "ICRF", "IDCP", "IHRP", "INPC", "INSM", "INSU", "INTC", "INTE", "INTX", "LBRI", "LICF", "LIFI", "LIMA", "LOAN", "LOAR", "LTCF", "MDCS", "MSVC", "NETT", "NITX", "NOWS", "NWCH", "NWCM", "OFEE", "OTHR", "OTLC", "PADD", "PAYR", "PENS", "PHON", "POPE", "PPTI", "PRCP", "PRME", "PTSP", "RCKE", "RCPT", "REFU", "RENT", "RINP", "RLWY", "ROYA", "SALA", "SAVG", "SCVE", "SECU", "SSBE", "STDY", "SUBS", "SUPP", "TAXS", "TELI", "TRAD", "TREA", "TRFD", "VATX", "VIEW", "WEBI", "WHLD", "WTER"}
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

// GeneratePNG .
func (c *Code) GeneratePNG(size int) ([]byte, error) {
	if c.Valid.Expired() {
		return nil, errors.New("negative validity period")
	}

	qrContent := c.String()
	if len(qrContent) > qrContentMaxSize {
		return nil, errors.New("qr content is too large")
	}

	q, err := qrcode.NewWith(qrContent,
		//qrcode.WithEncodingMode(qrcode.EncModeByte),
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionMedium))
	//qrcode.WithVersion(13)) -> With the standard's allowed 345 char size it will be bigger than 13, it will be 14
	// -> This is expected since the \n will trigger binary encoding which could be 331 in 13 and 362 in 14 version (with M error correction)
	if err != nil {
		return nil, err
	}

	if q.Dimension() > 73 { // 65 should be the max (ver 13), but 73 (ver 14) is the final size due to the binary encoding
		return nil, fmt.Errorf("generated QR code (width size) %d is too high (content too big)", q.Dimension())
	}

	buf := bytes.NewBuffer(nil)
	wr := standard.NewWithWriter(nopCloser{Writer: buf}, standard.WithQRWidth(uint8(size)))
	err = q.Save(wr)

	return buf.Bytes(), err
}

// String .
func (c *Code) String() string {
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

	if time.Time(c.Valid).IsZero() {
		// Add a default time with one hour expire
		sb.WriteString(date(time.Now().Add(time.Hour)).String())
	} else {
		sb.WriteString(c.Valid.String())
	}
	sb.WriteString("\n")

	if c.purpose != "" {
		sb.WriteString(c.purpose)
	}
	sb.WriteString("\n")

	if c.message != "" {
		sb.WriteString(c.message)
	}
	sb.WriteString("\n")

	if c.shopID != "" {
		sb.WriteString(c.shopID)
	}
	sb.WriteString("\n")

	if c.merchDevID != "" {
		sb.WriteString(c.merchDevID)
	}
	sb.WriteString("\n")

	if c.invoiceID != "" {
		sb.WriteString(c.invoiceID)
	}
	sb.WriteString("\n")

	if c.customerID != "" {
		sb.WriteString(c.customerID)
	}
	sb.WriteString("\n")

	if c.credTranID != "" {
		sb.WriteString(c.credTranID)
	}
	sb.WriteString("\n")

	if c.loyaltyID != "" {
		sb.WriteString(c.loyaltyID)
	}
	sb.WriteString("\n")

	if c.navCheckID != "" {
		sb.WriteString(c.navCheckID)
	}
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
func (c *Code) ValidUntil(t time.Time) error {
	if time.Now().After(t) {
		return errors.New("negative validity period")
	}
	c.Valid = date(t)
	return nil
}

// Purpose for the transaction
// Possible values are the AT-44 codes: https://www.rba.hr/documents/20182/183267/External+purpose+codes+list/8a28f888-1f83-5e29-d6ed-fce05f428689?version=1.1
func (c *Code) Purpose(purpose string) error {
	if len(purpose) != 4 {
		return errors.New("purpose has invalid length")
	}
	purpose = strings.ToUpper(purpose)

	found := false
	for _, c := range purposeCodes {
		if purpose == c {
			found = true
			break
		}
	}

	if !found {
		return errors.New("invalid purpose code")
	}

	c.purpose = purpose
	return nil
}

// Message .
func (c *Code) Message(msg string) error {
	if len(msg) > 70 {
		return errors.New("message is too long")
	}
	c.message = msg
	return nil
}

// ShopID .
func (c *Code) ShopID(shopID string) error {
	if len(shopID) > 35 {
		return errors.New("shopID is too long")
	}
	c.shopID = shopID
	return nil
}

// MerchDevID .
func (c *Code) MerchDevID(merchDevID string) error {
	if len(merchDevID) > 35 {
		return errors.New("merchDevID is too long")
	}
	c.merchDevID = merchDevID
	return nil
}

// InvoiceID .
func (c *Code) InvoiceID(invoiceID string) error {
	if len(invoiceID) > 35 {
		return errors.New("invoiceID is too long")
	}
	c.invoiceID = invoiceID
	return nil
}

// CustomerID .
func (c *Code) CustomerID(customerID string) error {
	if len(customerID) > 35 {
		return errors.New("customerID is too long")
	}
	c.customerID = customerID
	return nil
}

// CredTranID .
func (c *Code) CredTranID(credTranID string) error {
	if len(credTranID) > 35 {
		return errors.New("credTranID is too long")
	}
	c.credTranID = credTranID
	return nil
}

// LoyaltyID .
func (c *Code) LoyaltyID(loyaltyID string) error {
	if len(loyaltyID) > 35 {
		return errors.New("loyaltyID is too long")
	}
	c.loyaltyID = loyaltyID
	return nil
}

// NavCheckID .
func (c *Code) NavCheckID(navCheckID string) error {
	if len(navCheckID) > 35 {
		return errors.New("navCheckID is too long")
	}
	c.navCheckID = navCheckID
	return nil
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

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }
