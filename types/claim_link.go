package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type CLFee struct {
	Token         Token
	Amount        big.Int
	Authorization *string
}

type CLSource string

const (
	CLSourceUndefined CLSource = ""
	CLSourceD         CLSource = "d"
	CLSourceP2P       CLSource = "p2p"
)

type CLItemStatus string

const (
	CLItemStatusUndefined  CLItemStatus = ""
	CLItemStatusCreated    CLItemStatus = "created"
	CLItemStatusDepositing CLItemStatus = "depositing"
	CLItemStatusDeposited  CLItemStatus = "deposited"
	CLItemStatusRedeemed   CLItemStatus = "redeemed"
	CLItemStatusRedeeming  CLItemStatus = "redeeming"
	CLItemStatusError      CLItemStatus = "error"
	CLItemStatusRefunded   CLItemStatus = "refunded"
	CLItemStatusRefunding  CLItemStatus = "refunding"
	CLItemStatusCancelled  CLItemStatus = "cancelled"
)

type CLOperationStatus string

const (
	CLOperationStatusPending   CLOperationStatus = "pending"
	CLOperationStatusCompleted CLOperationStatus = "completed"
	CLOperationStatusError     CLOperationStatus = "error"
)

type CLOperation struct {
	Type      string
	Timestamp string
	Status    CLOperationStatus
	Receiver  common.Address
	TxHash    *common.Hash
}
