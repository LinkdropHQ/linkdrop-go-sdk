package types

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type ClaimLinkFee struct {
	Token         Token
	Amount        *big.Int
	Authorization []byte
}

func (clf *ClaimLinkFee) Validate() error {
	if !(clf.Token.Type == TokenTypeNative || clf.Token.Type == TokenTypeERC20) {
		return errors.New("fee token type is invalid, should be one of: native, ERC20")
	}
	return clf.Token.Validate()
}

type ClaimLinkFeeData struct {
	Amount            *big.Int     `json:"amount"`
	TotalAmount       *big.Int     `json:"total_amount"`
	MaxTransferAmount *big.Int     `json:"max_transfer_amount"`
	MinTransferAmount *big.Int     `json:"min_transfer_amount"`
	Fee               ClaimLinkFee `json:"fee"`
}

type ClaimLinkStatus int64

const (
	ClaimLinkStatusUndefined ClaimLinkStatus = iota
	ClaimLinkStatusCreated
	ClaimLinkStatusDepositing
	ClaimLinkStatusDeposited
	ClaimLinkStatusRedeeming
	ClaimLinkStatusRedeemed
	ClaimLinkStatusRefunding
	ClaimLinkStatusRefunded
	ClaimLinkStatusCancelled
	ClaimLinkStatusError
)

func (clis ClaimLinkStatus) String() string {
	switch clis {
	case ClaimLinkStatusUndefined:
		return ""
	case ClaimLinkStatusCreated:
		return "created"
	case ClaimLinkStatusDepositing:
		return "depositing"
	case ClaimLinkStatusDeposited:
		return "deposited"
	case ClaimLinkStatusRedeeming:
		return "redeeming"
	case ClaimLinkStatusRedeemed:
		return "redeemed"
	case ClaimLinkStatusRefunding:
		return "refunding"
	case ClaimLinkStatusRefunded:
		return "refunded"
	case ClaimLinkStatusCancelled:
		return "cancelled"
	case ClaimLinkStatusError:
		return "error"
	}
	return ""
}

func ClaimLinkStatusFromString(value string) ClaimLinkStatus {
	switch value {
	case "":
		return ClaimLinkStatusUndefined
	case "created":
		return ClaimLinkStatusCreated
	case "depositing":
		return ClaimLinkStatusDepositing
	case "deposited":
		return ClaimLinkStatusDeposited
	case "redeeming":
		return ClaimLinkStatusRedeeming
	case "redeemed":
		return ClaimLinkStatusRedeemed
	case "refunding":
		return ClaimLinkStatusRefunding
	case "refunded":
		return ClaimLinkStatusRefunded
	case "cancelled":
		return ClaimLinkStatusCancelled
	case "error":
		return ClaimLinkStatusError
	}
	return ClaimLinkStatusUndefined
}

type ClaimLinkOperationStatus string

const (
	LinkOperationStatusPending   ClaimLinkOperationStatus = "pending"
	LinkOperationStatusCompleted ClaimLinkOperationStatus = "completed"
	LinkOperationStatusError     ClaimLinkOperationStatus = "error"
)

type ClaimLinkOperation struct {
	Type      string                   `json:"type"`
	Timestamp string                   `json:"timestamp"`
	Status    ClaimLinkOperationStatus `json:"status"`
	Receiver  common.Address           `json:"receiver"`
	TxHash    *common.Hash             `json:"txHash"`
}

type ClaimLinkDepositParams struct {
	ChainId ChainId        `json:"chainId"`
	Value   *big.Int       `json:"value"`
	Data    []byte         `json:"data"`
	To      common.Address `json:"to"`
}

func (cldp ClaimLinkDepositParams) MarshalJSON() ([]byte, error) {
	type Alias ClaimLinkDepositParams
	return json.Marshal(&struct {
		Data  string `json:"data"`
		Value string `json:"value"`
		*Alias
	}{
		Data:  "0x" + common.Bytes2Hex(cldp.Data),
		Value: cldp.Value.String(),
		Alias: (*Alias)(&cldp),
	})
}

// Link
// Represents the parsed structure of the link
type Link struct {
	SenderSignature []byte
	LinkKey         ecdsa.PrivateKey
	TransferId      common.Address
	ChainId         ChainId
	Version         string
	Message         *EncryptedMessage
	Sender          *common.Address
}
