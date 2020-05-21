package qr

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHUFAmount(t *testing.T) {
	c := Code{}

	assert.Equal(t, "HUF0", c.Amount.String())
	assert.NoError(t, c.HUFAmount(100))
	assert.Equal(t, "HUF100", c.Amount.String())

	assert.Equal(t, "amount could not be negative", c.HUFAmount(-1).Error())
	assert.Equal(t, "amount could not be higher than 999999999999", c.HUFAmount(1234567890123).Error())
}

func TestPurpose(t *testing.T) {
	c := &Code{}

	assert.Equal(t, "purpose has invalid length", c.Purpose("a").Error())
	assert.Equal(t, "invalid purpose code", c.Purpose("abcd").Error())
	assert.NoError(t, c.Purpose("ACCT"))
	assert.Equal(t, "ACCT", c.purpose)

}

func TestAddRecipientChecks(t *testing.T) {
	testTable := []struct {
		inputBIC    string
		inputName   string
		inputIBAN   string
		expectedErr string
	}{
		// BIC checks
		{"a", "", "", "invalid BIC length"},
		{"abcdefgh", "", "HU00123456789012345678901234", ""}, // 8 char -> auto extend
		{"abcdefghi", "", "", "invalid BIC length"},
		{"abcdefghijk", "", "HU00123456789012345678901234", ""},
		{"abcdefghijke", "", "", "invalid BIC length"},

		// Name Checks
		{"abcdefghijk", "Test User", "HU00123456789012345678901234", ""},
		{"abcdefghijk", strings.Repeat("a", 70), "HU00123456789012345678901234", ""},
		{"abcdefghijk", strings.Repeat("a", 71), "", "name should not be longer than 70"},

		// IBAN check
		{"abcdefghijk", "Test User", "HU0012345678901234567890123", "invalid IBAN length"},
		{"abcdefghijk", "Test User", "HU001234567890123456789012345", "invalid IBAN length"},
		{"abcdefghijk", "Test User", "HU00123456789012345678901234", ""},
	}

	for _, tt := range testTable {
		c := &Code{}
		if tt.expectedErr != "" {
			assert.Equal(t, tt.expectedErr, addRecipient(c, tt.inputBIC, tt.inputName, tt.inputIBAN).Error())
		} else {
			assert.NoError(t, addRecipient(c, tt.inputBIC, tt.inputName, tt.inputIBAN))
		}
	}

	// One more overall checks
	c := &Code{}
	assert.NoError(t, addRecipient(c, "abcdefgh", "Test User", "HU00123456789012345678901234"))
	assert.Equal(t, "abcdefghXXX", c.BIC) // Check auto extend here too
	assert.Equal(t, "Test User", c.Name)
	assert.Equal(t, "HU00123456789012345678901234", c.IBAN)
}

func TestCodeFormat(t *testing.T) {
	c, err := NewPaymentSend("abcdefgh", "Test User", "HU00123456789012345678901234")
	assert.NoError(t, err)

	output := c.String()
	assert.Equal(t, 17, strings.Count(output, "\n"), "has all (even empty) lines")
	assert.False(t, strings.HasPrefix(output, "\n"), "does not start with new line")
	assert.True(t, strings.HasSuffix(output, "\n"), "ends with new line")
}

func TestCodeFormatDetailed(t *testing.T) {
	c, err := NewPaymentRequest("abcdefgh", "Test User", "HU00123456789012345678901234")
	assert.NoError(t, err)

	assert.NoError(t, c.HUFAmount(500))
	assert.NoError(t, c.Purpose("AGRT"))
	assert.NoError(t, c.Message("hello!"))
	assert.NoError(t, c.ShopID("shopIDHere"))
	assert.NoError(t, c.MerchDevID("merchDevID"))
	assert.NoError(t, c.InvoiceID("invoiceID"))
	assert.NoError(t, c.CustomerID("cccustomer"))
	assert.NoError(t, c.CredTranID("credTransID"))
	assert.NoError(t, c.LoyaltyID("loyID"))
	assert.NoError(t, c.NavCheckID("navhere"))

	output := strings.Split(c.String(), "\n")
	assert.Len(t, output, 18)

	// Field checks
	assert.Equal(t, KindRTP.String(), output[0])
	assert.Equal(t, "001", output[1])
	assert.Equal(t, "1", output[2])
	assert.Equal(t, "abcdefghXXX", output[3])
	assert.Equal(t, "Test User", output[4])
	assert.Equal(t, "HU00123456789012345678901234", output[5])
	assert.Equal(t, "HUF500", output[6]) // Amount

	// Valid checks, trim timezone, parse and check with now, it was empty so it should be somewhere now+1
	valid := strings.Split(output[7], "+")
	assert.Len(t, valid, 2)
	vt, err := time.Parse("20060102150405", valid[0])
	assert.NoError(t, err)
	assert.True(t, vt.After(time.Now()))

	assert.Equal(t, "AGRT", output[8])         // Purpose
	assert.Equal(t, "hello!", output[9])       // Message
	assert.Equal(t, "shopIDHere", output[10])  // shopID
	assert.Equal(t, "merchDevID", output[11])  // merchDevID
	assert.Equal(t, "invoiceID", output[12])   // invoiceID
	assert.Equal(t, "cccustomer", output[13])  // customerID
	assert.Equal(t, "credTransID", output[14]) // credTranID
	assert.Equal(t, "loyID", output[15])       // loyaltyID
	assert.Equal(t, "navhere", output[16])     // navCheckID
	assert.Equal(t, "", output[17])            // Empty line at the end
}

func TestCodeFormatDateCheck(t *testing.T) {
	c, err := NewPaymentRequest("abcdefgh", "Test User", "HU00123456789012345678901234")
	assert.NoError(t, err)

	ts := time.Now().Add(4 * time.Hour).UTC()
	assert.NoError(t, c.ValidUntil(ts))
	output := strings.Split(c.String(), "\n")
	assert.Len(t, output, 18)

	assert.Equal(t, date(ts).String(), output[7])
}

func TestCodeSetErrors(t *testing.T) {
	c := &Code{}

	assert.Equal(t, "negative validity period", c.ValidUntil(time.Now().Add(-time.Minute)).Error())
	assert.Equal(t, "message is too long", c.Message(strings.Repeat("a", 71)).Error())
	assert.Equal(t, "shopID is too long", c.ShopID(strings.Repeat("a", 36)).Error())
	assert.Equal(t, "merchDevID is too long", c.MerchDevID(strings.Repeat("a", 36)).Error())
	assert.Equal(t, "invoiceID is too long", c.InvoiceID(strings.Repeat("a", 36)).Error())
	assert.Equal(t, "customerID is too long", c.CustomerID(strings.Repeat("a", 36)).Error())
	assert.Equal(t, "credTranID is too long", c.CredTranID(strings.Repeat("a", 36)).Error())
	assert.Equal(t, "loyaltyID is too long", c.LoyaltyID(strings.Repeat("a", 36)).Error())
	assert.Equal(t, "navCheckID is too long", c.NavCheckID(strings.Repeat("a", 36)).Error())
}

func TestGeneratePNG(t *testing.T) {
	c, err := NewPaymentSend("abcdefgh", "Test User", "HU00123456789012345678901234")
	assert.NoError(t, err)

	_, err = c.GeneratePNG(256)
	assert.Equal(t, "negative validity period", err.Error())

	_ = c.ValidUntil(time.Now().Add(time.Hour))
	_, err = c.GeneratePNG(256)
	assert.NoError(t, err)

	// Fill all the fields and gen again and try to hit the version error
	c = genFullCode(t)
	_, err = c.GeneratePNG(64)
	assert.Equal(t, "qr content is too large", err.Error())

	// Test the max size (generated size with full content: 483, so remove some fields here
	c.navCheckID = ""   // -35
	c.loyaltyID = ""    // -35
	c.credTranID = ""   // -35
	c.customerID = "12" // remaining -33

	assert.Equal(t, qrContentMaxSize, len(c.String())) // This should be fine and check the qr after this
	_, err = c.GeneratePNG(64)
	assert.NoError(t, err)
}

func TestFullCode(t *testing.T) {
	c := genFullCode(t)
	genStr := c.String()
	splitted := strings.Split(genStr, "\n")
	assert.Equal(t, 3, len(splitted[0]))
	assert.Equal(t, 3, len(splitted[1]))
	assert.Equal(t, 1, len(splitted[2]))
	assert.Equal(t, 11, len(splitted[3]))
	assert.Equal(t, 70, len(splitted[4]))
	assert.Equal(t, 28, len(splitted[5]))
	assert.Equal(t, 15, len(splitted[6]))
	assert.Equal(t, 16, len(splitted[7]))
	assert.Equal(t, 4, len(splitted[8]))
	assert.Equal(t, 70, len(splitted[9]))
	assert.Equal(t, 35, len(splitted[10]))
	assert.Equal(t, 35, len(splitted[11]))
	assert.Equal(t, 35, len(splitted[12]))
	assert.Equal(t, 35, len(splitted[13]))
	assert.Equal(t, 35, len(splitted[14]))
	assert.Equal(t, 35, len(splitted[15]))
	assert.Equal(t, 35, len(splitted[16]))

	assert.Equal(t, 0, len(splitted[17])) // Empty line at the end

	assert.Equal(t, 483, len(genStr)) // The fully packed code's length, but the standard won't allow that
}

func genFullCode(t *testing.T) *Code {
	c, err := NewPaymentSend("abcdefgh", strings.Repeat("a", 70), "HU00123456789012345678901234")
	assert.NoError(t, err)

	c.Version = version("111")
	c.Charset = 2
	assert.NoError(t, c.HUFAmount(999999999999))
	assert.NoError(t, c.ValidUntil(time.Date(2120, 03, 30, 10, 11, 12, 0, time.FixedZone("testZone", 11))))
	assert.NoError(t, c.Purpose("ACCT"))
	assert.NoError(t, c.Message(strings.Repeat("b", 70)))
	assert.NoError(t, c.ShopID(strings.Repeat("c", 35)))
	assert.NoError(t, c.MerchDevID(strings.Repeat("d", 35)))
	assert.NoError(t, c.InvoiceID(strings.Repeat("e", 35)))
	assert.NoError(t, c.CustomerID(strings.Repeat("f", 35)))
	assert.NoError(t, c.CredTranID(strings.Repeat("g", 35)))
	assert.NoError(t, c.LoyaltyID(strings.Repeat("h", 35)))
	assert.NoError(t, c.NavCheckID(strings.Repeat("i", 35)))
	return c
}
