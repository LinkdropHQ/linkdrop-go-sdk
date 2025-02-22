package constants

const (
	ApiURL                      = "https://escrow-api.linkdrop.io/v3"
	DashboardApiUrl             = "https://escrow-api.linkdrop.io/dashboard"
	DevDashboardApiUrl          = "https://escrow-api.linkdrop.io/staging"
	NativeTokenAddress          = "0x0000000000000000000000000000000000000000"
	EscrowContractAddress       = "0xbe7b40eb3a9d85d3a76142cb637ab824f0d35ead"
	EscrowNFTContractAddress    = "0x5fc1316119a1b7cec52a2984c62764343dca70c9"
	CbwEscrowContractAddress    = "0x5badb0143f69015c5c86cbd9373474a9c8ab713b"
	CbwEscrowNFTContractAddress = "0x3c74782de03c0402d207fe41307fe50fe9b6b5c7"
)

var SupportedStableCoins = map[TokenAddress]Selector{
	UsdcBase:           ReceiveWithAuthorization,
	EurcBase:           ReceiveWithAuthorization,
	UsdcBridgedPolygon: ApproveWithAuthorization,
	UsdcPolygon:        ReceiveWithAuthorization,
	UsdcArbitrum:       ReceiveWithAuthorization,
	UsdcOptimism:       ReceiveWithAuthorization,
	UsdcAvalanche:      ReceiveWithAuthorizationEOA,
	CbBTC:              ReceiveWithAuthorizationEOA,
}
