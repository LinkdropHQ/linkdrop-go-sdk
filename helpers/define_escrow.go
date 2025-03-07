package helpers

import (
	"fmt"
	"github.com/LinkdropHQ/linkdrop-go-sdk/constants"
	"github.com/ethereum/go-ethereum/common"
)

func DefineEscrowVersion(address common.Address) (version string, err error) {
	for ver, addresses := range constants.Escrows {
		for _, addr := range addresses {
			if addr == address {
				return ver, nil
			}
		}
	}
	return "", fmt.Errorf("address not found in escrows")
}
