package types

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var ZeroAddress = common.Address{}

type TokenType string

const (
	TokenTypeNative  TokenType = "native"
	TokenTypeERC20   TokenType = "ERC20"
	TokenTypeERC721  TokenType = "ERC721"
	TokenTypeERC1155 TokenType = "ERC1155"
)

type Token struct {
	Type    TokenType      `json:"type"`
	ChainId ChainId        `json:"chainId"`
	Address common.Address `json:"address"`
	Id      *big.Int       `json:"id"`
}

func (t *Token) Validate() error {
	if !IsChainSupported(t.ChainId) {
		return errors.New("chain is not supported")
	}

	if t.Type == TokenTypeNative && t.Address != ZeroAddress {
		return errors.New("native token should not have address")
	}

	if t.Type != TokenTypeNative && t.Address == ZeroAddress {
		return errors.New("address is not provided")
	}

	if (t.Type == TokenTypeERC721 || t.Type == TokenTypeERC1155) && t.Id == nil {
		return errors.New("id is not provided")
	}

	return nil
}
