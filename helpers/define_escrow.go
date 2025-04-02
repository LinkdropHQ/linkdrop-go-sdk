package helpers

import (
	"fmt"
	"github.com/LinkdropHQ/linkdrop-go-sdk/constants"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
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

func EscrowAddressByChain(
	chain types.ChainId,
) (escrow common.Address, escrowNFT common.Address, err error) {
	switch chain {
	case types.ChainIdBase:
		escrow = constants.Escrows["3.2"][0]
		escrowNFT = constants.Escrows["3.2"][1]
	case types.ChainIdPolygon, types.ChainIdAvalanche, types.ChainIdOptimism, types.ChainIdArbitrum:
		escrow = constants.Escrows["3.2"][2]
		escrowNFT = constants.Escrows["3.2"][3]
	default:
		err = fmt.Errorf("chain not supported")
	}
	return
}

func EscrowAddressForToken(
	token types.Token,
) (escrowAddress common.Address, err error) {
	escrow, escrowNFT, err := EscrowAddressByChain(token.ChainId)
	if err != nil {
		return
	}
	switch token.Type {
	case types.TokenTypeERC1155, types.TokenTypeERC721:
		escrowAddress = escrowNFT
	default:
		escrowAddress = escrow
	}
	return
}
