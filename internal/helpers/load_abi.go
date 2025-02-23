package helpers

import (
	"github.com/LinkdropHQ/linkdrop-go-sdk/internal/constants"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"os"
	"strings"
)

func LoadABI() (err error) {
	abiFile, err := os.ReadFile("./abi/LinkdropEscrowNFT.json")
	if err != nil {
		return
	}
	constants.EscrowNFTAbi, err = abi.JSON(strings.NewReader(string(abiFile)))
	if err != nil {
		return err
	}

	abiFile, err = os.ReadFile("./abi/LinkdropEscrowToken.json")
	if err != nil {
		return
	}
	constants.EscrowTokenAbi, err = abi.JSON(strings.NewReader(string(abiFile)))
	if err != nil {
		return err
	}
	return
}
