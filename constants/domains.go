package constants

import (
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

var (
	DomainUsdcArbitrum = apitypes.TypedDataDomain{
		Name:              "USD Coin",
		Version:           "2",
		ChainId:           math.NewHexOrDecimal256(42161),
		VerifyingContract: TAUsdcArbitrum.Hex(),
	}

	DomainUsdcAvalanche = apitypes.TypedDataDomain{
		Name:              "USD Coin",
		Version:           "2",
		ChainId:           math.NewHexOrDecimal256(43114),
		VerifyingContract: TAUsdcAvalanche.Hex(),
	}

	DomainCbBtcBase = apitypes.TypedDataDomain{
		Name:              "Coinbase Wrapped BTC",
		Version:           "2",
		ChainId:           math.NewHexOrDecimal256(8453),
		VerifyingContract: TACbBtcBase.Hex(),
	}

	DomainEurcBase = apitypes.TypedDataDomain{
		Name:              "EURC",
		Version:           "2",
		ChainId:           math.NewHexOrDecimal256(8453),
		VerifyingContract: TAEurcBase.Hex(),
	}

	DomainUsdcBase = apitypes.TypedDataDomain{
		Name:              "USD Coin",
		Version:           "2",
		ChainId:           math.NewHexOrDecimal256(8453),
		VerifyingContract: TAUsdcBase.Hex(),
	}

	DomainUsdcOptimism = apitypes.TypedDataDomain{
		Name:              "USD Coin",
		Version:           "2",
		ChainId:           math.NewHexOrDecimal256(10),
		VerifyingContract: TAUsdcOptimism.Hex(),
	}

	DomainUsdcBridgedPolygon = apitypes.TypedDataDomain{
		Name:              "USD Coin (PoS)",
		Version:           "1",
		ChainId:           math.NewHexOrDecimal256(137),
		VerifyingContract: TAUsdcBridgedPolygon.Hex(),
		Salt:              "0x0000000000000000000000000000000000000000000000000000000000000089",
	}

	DomainUsdcPolygon = apitypes.TypedDataDomain{
		Name:              "USD Coin",
		Version:           "2",
		ChainId:           math.NewHexOrDecimal256(137),
		VerifyingContract: TAUsdcPolygon.Hex(),
	}
)
