package types

type TransferLimits struct {
	MinAmount    string `json:"minAmount"`
	MaxAmount    string `json:"maxAmount"`
	MinAmountUSD string `json:"minAmountUSD"`
	MaxAmountUSD string `json:"maxAmountUSD"`
}
