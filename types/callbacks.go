package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"math/big"
)

type RandomBytesCallback func(length int64) []byte

type SignTypedDataCallback func(typedData apitypes.TypedData) ([]byte, error)

type SendTransactionCallback func(chainId *big.Int, to common.Address, value *big.Int, data []byte) (*Transaction, error)
