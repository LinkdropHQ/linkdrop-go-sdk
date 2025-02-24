package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type RandomBytesCallback func(length int64) []byte

type SignTypedDataCallback func(domain TypedDataDomain, types map[string][]TypedDataField, value map[string]any) ([]byte, error)

type SendTransactionCallback func(to common.Address, value *big.Int, data []byte) (common.Hash, error)
