package helpers

import (
	_ "embed"
	"github.com/LinkdropHQ/linkdrop-go-sdk/constants"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
)

//go:embed abi/LinkdropEscrowNFT.json
var escrowNFTJson []byte

//go:embed abi/LinkdropEscrowToken.json
var escrowTokenJson []byte

func LoadABI() (err error) {
	abiRaw := strings.NewReader(string(escrowNFTJson))
	constants.EscrowNFTAbi, err = abi.JSON(abiRaw)
	if err != nil {
		return err
	}

	constants.EscrowTokenAbi, err = abi.JSON(strings.NewReader(string(escrowTokenJson)))
	if err != nil {
		return err
	}

	return
}
