package types

type ChainId int64

const (
	ChainIdPolygon   ChainId = 137
	ChainIdBase      ChainId = 8453
	ChainIdArbitrum  ChainId = 42161
	ChainIdOptimism  ChainId = 10
	ChainIdAvalanche ChainId = 43114
)

func IsChainSupported(chainId ChainId) bool {
	return chainId == ChainIdPolygon || chainId == ChainIdBase || chainId == ChainIdArbitrum || chainId == ChainIdOptimism || chainId == ChainIdAvalanche
}
