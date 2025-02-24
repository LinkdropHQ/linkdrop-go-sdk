package types

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type CLFeeData struct {
	Amount            *big.Int `json:"amount"`
	TotalAmount       *big.Int `json:"total_amount"`
	MaxTransferAmount *big.Int `json:"max_transfer_amount"`
	MinTransferAmount *big.Int `json:"min_transfer_amount"`
	Fee               CLFee    `json:"fee"`
}

type CLFee struct {
	Token         Token
	Amount        *big.Int
	Authorization string // TODO []byte
}

type CLSource string

const (
	CLSourceUndefined CLSource = ""
	CLSourceD         CLSource = "d"
	CLSourceP2P       CLSource = "p2p"
)

type CLItemStatus int64

const (
	CLItemStatusUndefined CLItemStatus = iota
	CLItemStatusCreated
	CLItemStatusDepositing
	CLItemStatusDeposited
	CLItemStatusRedeeming
	CLItemStatusRedeemed
	CLItemStatusRefunding
	CLItemStatusRefunded
	CLItemStatusCancelled
	CLItemStatusError
)

func ClItemStatusFromString(itemStatus string) CLItemStatus {
	switch itemStatus {
	case "created":
		return CLItemStatusCreated
	case "depositing":
		return CLItemStatusDepositing
	case "deposited":
		return CLItemStatusDeposited
	case "redeeming":
		return CLItemStatusRedeeming
	case "redeemed":
		return CLItemStatusRedeemed
	case "refunding":
		return CLItemStatusRefunding
	case "refunded":
		return CLItemStatusRefunded
	case "cancelled":
		return CLItemStatusCancelled
	case "error":
		return CLItemStatusError
	}
	return CLItemStatusUndefined
}

func (clis CLItemStatus) String() string {
	switch clis {
	case CLItemStatusUndefined:
		return ""
	case CLItemStatusCreated:
		return "created"
	case CLItemStatusDepositing:
		return "depositing"
	case CLItemStatusDeposited:
		return "deposited"
	case CLItemStatusRedeeming:
		return "redeeming"
	case CLItemStatusRedeemed:
		return "redeemed"
	case CLItemStatusRefunding:
		return "refunding"
	case CLItemStatusRefunded:
		return "refunded"
	case CLItemStatusCancelled:
		return "cancelled"
	case CLItemStatusError:
		return "error"
	}
	return ""
}

type CLOperationStatus string

const (
	CLOperationStatusPending   CLOperationStatus = "pending"
	CLOperationStatusCompleted CLOperationStatus = "completed"
	CLOperationStatusError     CLOperationStatus = "error"
)

type CLOperation struct {
	Type      string            `json:"type"`
	Timestamp string            `json:"timestamp"`
	Status    CLOperationStatus `json:"status"`
	Receiver  common.Address    `json:"receiver"`
	TxHash    *common.Hash      `json:"txHash"`
}

type CLDepositParams struct {
	Value *big.Int
	Data  []byte
	To    common.Address
}

type Link struct {
	SenderSig              string
	LinkKey                *ecdsa.PrivateKey
	TransferId             common.Address
	ChainId                ChainId
	Version                string
	EncryptionKey          *[]byte
	EncryptionKeyLinkParam *[]byte
	Sender                 *common.Address
}
