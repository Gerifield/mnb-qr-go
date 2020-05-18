package qr

type Code struct {
	Kind            Kind     // Required
	Version         [3]byte  // Required
	Charset         [1]byte  // Required
	BIC             [11]byte // Required
	Name            [70]byte // Required
	IBAN            [28]byte // Required
	Amount          [15]byte
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

type Kind string

var (
	// KindHCT for send money
	KindHCT Kind = "HCT"

	// KindRTP for request money
	KindRTP Kind = "RTP"
)

func NewPaymentSend() Code {
	return Code{
		Kind: KindHCT,
	}
}

func NewPaymentRequest() Code {
	return Code{
		Kind: KindRTP,
	}
}
