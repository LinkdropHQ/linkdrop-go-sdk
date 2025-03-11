package helpers

import (
	"errors"
	"github.com/LinkdropHQ/linkdrop-go-sdk/constants"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func DefineDomain(token types.Token) (apitypes.TypedDataDomain, error) {
	chainID := token.ChainId
	address := token.Address

	switch chainID {
	case types.ChainIdBase: // Base
		if address == constants.TAUsdcBase {
			return constants.DomainUsdcBase, nil
		}
		if address == constants.TAEurcBase {
			return constants.DomainEurcBase, nil
		}
		if address == constants.TACbBtcBase {
			return constants.DomainCbBtcBase, nil
		}
	}
	return apitypes.TypedDataDomain{}, errors.New("domain not found")
}
