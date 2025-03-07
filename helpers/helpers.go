package helpers

import (
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"math/big"
)

func DefineValue(
	token types.Token,
	fee types.ClaimLinkFee,
	totalAmount *big.Int,
) (*big.Int, error) {
	if token.Type == types.TokenTypeNative {
		return totalAmount, nil
	}
	if fee.Token.Address == token.Address {
		return big.NewInt(0), nil
	}

	return fee.Amount, nil
}
