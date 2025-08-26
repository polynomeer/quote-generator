package model

type Quote struct {
	Symbol   string  `json:"symbol"`
	Price    float64 `json:"price"`
	Bid      float64 `json:"bid"`
	Ask      float64 `json:"ask"`
	Volume   int64   `json:"volume"`
	Currency string  `json:"currency"`
	Ts       int64   `json:"ts"` // ms epoch
}

type Bar struct {
	Symbol string  `json:"symbol"`
	TF     string  `json:"tf"` // e.g. "1s"
	O      float64 `json:"o"`
	H      float64 `json:"h"`
	L      float64 `json:"l"`
	C      float64 `json:"c"`
	V      int64   `json:"v"`
	Ts     int64   `json:"ts"`
}
