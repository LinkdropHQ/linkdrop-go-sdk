package types

import "math/big"

type TransferLimits struct {
	MinAmount    *big.Int `json:"min_transfer_amount"`
	MaxAmount    *big.Int `json:"max_transfer_amount"`
	MinAmountUSD *big.Int `json:"min_transfer_amount_usd"`
	MaxAmountUSD *big.Int `json:"max_transfer_amount_usd"`
}
