package model

type Algorithm struct {
	ID       int  `json:"id"`
	ClientID int  `json:"clientID"`
	VWAP     bool `json:"vwap"`
	TWAP     bool `json:"twap"`
	HFT      bool `json:"hft"`
}
