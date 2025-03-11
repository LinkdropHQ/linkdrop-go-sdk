package types

import "github.com/ethereum/go-ethereum/common"

type ChainConfig struct {
	ChainId       ChainId
	EscrowAddress common.Address
}

type ChainId int64

const (
	ChainIdBase ChainId = 8453
)

func (cid *ChainId) IsSupported() bool {
	return *cid == ChainIdBase
}
