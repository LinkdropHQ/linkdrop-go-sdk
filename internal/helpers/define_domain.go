package helpers

import (
	"errors"
	"github.com/LinkdropHQ/linkdrop-go-sdk/internal/constants"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func DefineDomain(token types.Token) (apitypes.TypedDataDomain, error) {
	chainID := token.ChainId
	address := token.Address

	switch chainID {
	case types.ChainIdPolygon: // Polygon
		if address == constants.TAUsdcBridgedPolygon {
			return constants.DomainUsdcBridgedPolygon, nil
		}
		if address == constants.TAUsdcPolygon {
			return constants.DomainUsdcPolygon, nil
		}

	case types.ChainIdAvalanche: // Avalanche
		if address == constants.TAUsdcAvalanche {
			return constants.DomainUsdcAvalanche, nil
		}

	case types.ChainIdOptimism: // Optimism
		if address == constants.TAUsdcOptimism {
			return constants.DomainUsdcOptimism, nil
		}

	case types.ChainIdArbitrum: // Arbitrum
		if address == constants.TAUsdcArbitrum {
			return constants.DomainUsdcArbitrum, nil
		}

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
