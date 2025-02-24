package constants

import "github.com/ethereum/go-ethereum/common"

const (
	ApiURL             = "https://escrow-api.linkdrop.io/v3"
	DashboardApiUrl    = "https://escrow-api.linkdrop.io/dashboard"
	DevDashboardApiUrl = "https://escrow-api.linkdrop.io/staging"
)

var (
	NativeTokenAddress          = common.HexToAddress("0x0000000000000000000000000000000000000000")
	EscrowContractAddress       = common.HexToAddress("0xbe7b40eb3a9d85d3a76142cb637ab824f0d35ead")
	EscrowNFTContractAddress    = common.HexToAddress("0x5fc1316119a1b7cec52a2984c62764343dca70c9")
	CbwEscrowContractAddress    = common.HexToAddress("0x5badb0143f69015c5c86cbd9373474a9c8ab713b")
	CbwEscrowNFTContractAddress = common.HexToAddress("0x3c74782de03c0402d207fe41307fe50fe9b6b5c7")
)

var SupportedStableCoins = map[common.Address]Selector{
	TAUsdcBase:           SelectorReceiveWithAuthorization,
	TAEurcBase:           SelectorReceiveWithAuthorization,
	TAUsdcBridgedPolygon: SelectorApproveWithAuthorization,
	TAUsdcPolygon:        SelectorReceiveWithAuthorization,
	TAUsdcArbitrum:       SelectorReceiveWithAuthorization,
	TAUsdcOptimism:       SelectorReceiveWithAuthorization,
	TAUsdcAvalanche:      SelectorReceiveWithAuthorizationEOA,
	TACbBTC:              SelectorReceiveWithAuthorizationEOA,
}
