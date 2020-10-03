package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gerifield/mnb-qr-go/src/qr"
)

var (
	errInvalidKind = errors.New("invalid kind (should be RTP or HCT)")
	errInvalidSize = errors.New("invalid PNG size")
)

type Srv struct {
}

func New() *Srv {
	return &Srv{}
}

func (s *Srv) GenerateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, http.StatusMethodNotAllowed, errors.New("invalid method"))
		return
	}

	var input struct {
		Kind    string `json:"kind"` // HCT/RTP
		BIC     string `json:"bic"`
		Name    string `json:"name"`
		IBAN    string `json:"iban"`
		Expire  int    `json:"expire"`  // Expire (duration) in seconds
		PNGSize int    `json:"pngSize"` // Size in pixel

		Amount     int    `json:"amount"`     // Optional, HUF only
		Purpose    string `json:"purpose"`    // Optional
		Message    string `json:"message"`    // Optional
		ShopID     string `json:"shopID"`     // Optional
		MerchDevID string `json:"merchDevID"` // Optional
		InvoiceID  string `json:"invoiceID"`  // Optional
		CustomerID string `json:"customerID"` // Optional
		CredTranID string `json:"credTranID"` // Optional
		LoyaltyID  string `json:"loyaltyID"`  // Optional
		NavCheckID string `json:"navCheckID"` // Optional
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	if input.PNGSize == 0 {
		sendError(w, http.StatusBadRequest, errInvalidSize)
		return
	}

	var c *qr.Code
	switch input.Kind {
	case string(qr.KindRTP):
		c, err = qr.NewPaymentRequest(input.BIC, input.Name, input.IBAN)
	case string(qr.KindHCT):
		c, err = qr.NewPaymentSend(input.BIC, input.Name, input.IBAN)
	default:
		sendError(w, http.StatusBadRequest, errInvalidKind)
		return
	}
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	// Set all the fields
	err = c.ValidUntil(time.Now().Add(time.Second * time.Duration(input.Expire)))
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	if input.Amount > 0 {
		err = c.HUFAmount(input.Amount)
		if err != nil {
			sendError(w, http.StatusBadRequest, err)
			return
		}
	}

	if input.Purpose != "" {
		err = c.Purpose(input.Purpose)
		if err != nil {
			sendError(w, http.StatusBadRequest, err)
			return
		}
	}

	err = c.Message(input.Message)
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	err = c.ShopID(input.ShopID)
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	err = c.MerchDevID(input.MerchDevID)
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	err = c.InvoiceID(input.InvoiceID)
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	err = c.CustomerID(input.CustomerID)
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	err = c.CredTranID(input.CredTranID)
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	err = c.CustomerID(input.CustomerID)
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	err = c.NavCheckID(input.NavCheckID)
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	// Generate the image
	b, err := c.GeneratePNG(input.PNGSize)
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	// Display the image with disabled cache
	w.Header().Add("Content-Type", "image/png")
	w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Add("Pragma", "no-cache")
	w.Header().Add("Expires", "0")
	_, _ = w.Write(b)
}
