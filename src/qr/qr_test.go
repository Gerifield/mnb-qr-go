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

func TestGenerateQR(t *testing.T) {
	c, err := NewPaymentSend("abcdefgh", "Test User", "HU00123456789012345678901234")
	assert.NoError(t, err)

	_, err = c.GeneratePNG(256)
	assert.Equal(t, "negative validity period", err.Error())

	_ = c.ValidUntil(time.Now().Add(time.Hour))
	_, err = c.GeneratePNG(256)
	assert.NoError(t, err)
}
