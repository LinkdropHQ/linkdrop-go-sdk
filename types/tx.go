package types

import "github.com/ethereum/go-ethereum/common"

type TransactionType string

const (
	TransactionTypeUserOp TransactionType = "userOp"
	TransactionTypeTx     TransactionType = "tx"
)

type Transaction struct {
	Hash common.Hash
	Type TransactionType
}
