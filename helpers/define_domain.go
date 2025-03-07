package helpers

import (
	"errors"
	constants2 "github.com/LinkdropHQ/linkdrop-go-sdk/constants"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func DefineDomain(token types.Token) (apitypes.TypedDataDomain, error) {
	chainID := token.ChainId
	address := token.Address

	switch chainID {
	case types.ChainIdPolygon: // Polygon
		if address == constants2.TAUsdcBridgedPolygon {
			return constants2.DomainUsdcBridgedPolygon, nil
		}
		if address == constants2.TAUsdcPolygon {
			return constants2.DomainUsdcPolygon, nil
		}

	case types.ChainIdAvalanche: // Avalanche
		if address == constants2.TAUsdcAvalanche {
			return constants2.DomainUsdcAvalanche, nil
		}

	case types.ChainIdOptimism: // Optimism
		if address == constants2.TAUsdcOptimism {
			return constants2.DomainUsdcOptimism, nil
		}

	case types.ChainIdArbitrum: // Arbitrum
		if address == constants2.TAUsdcArbitrum {
			return constants2.DomainUsdcArbitrum, nil
		}

	case types.ChainIdBase: // Base
		if address == constants2.TAUsdcBase {
			return constants2.DomainUsdcBase, nil
		}
		if address == constants2.TAEurcBase {
			return constants2.DomainEurcBase, nil
		}
		if address == constants2.TACbBtcBase {
			return constants2.DomainCbBtcBase, nil
		}
	}
	return apitypes.TypedDataDomain{}, errors.New("domain not found")
}
