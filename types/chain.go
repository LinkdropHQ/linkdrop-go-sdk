package types

import (
	"github.com/ethereum/go-ethereum/common"
)

type ChainConfig struct {
	ChainId       ChainId
	EscrowAddress common.Address
}

type ChainId int64

const (
	ChainIdBase      ChainId = 8453
	ChainIdPolygon   ChainId = 137
	ChainIdAvalanche ChainId = 43114
	ChainIdOptimism  ChainId = 10
	ChainIdArbitrum  ChainId = 42161
)

func (cid *ChainId) IsSupported() bool {
	switch *cid {
	case ChainIdBase, ChainIdPolygon, ChainIdAvalanche, ChainIdOptimism, ChainIdArbitrum:
		return true
	}
	return false
}
