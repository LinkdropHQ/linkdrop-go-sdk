package helpers

import (
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

func TransferId(getRandomBytes types.RandomBytesCallback) (transferId common.Address, err error) {
	pk, err := PrivateKey(getRandomBytes)
	if err != nil {
		return
	}
	transferId, err = AddressFromPrivateKey(pk)
	if err != nil {
		return
	}
	return
}
