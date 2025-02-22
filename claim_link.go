package linkdrop

import (
	"errors"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
)

type ClaimLink struct {
	SDK                    *SDK
	Sender                 common.Address // Required
	Token                  types.Token
	Amount                 big.Int
	PendingTxs             *int64
	PendingTxSubmittedBn   *int64
	PendingTxSubmittedAt   *int64
	PendingBlocks          *int64
	Fee                    *types.CLFee
	TotalAmount            *big.Int
	Expiration             time.Duration
	TransferId             string
	EscrowAddress          *common.Address
	Operations             []types.CLOperation
	LinkKey                *string
	ClaimUrl               *string
	ForRecipient           bool
	Status                 types.CLItemStatus // CLItemStatusUndefined by default
	Source                 types.CLSource     // CLSourceUndefined by default
	EncryptedSenderMessage *string
	SenderMessage          string
}

func (cl *ClaimLink) Validate() (err error) {
	err = cl.Token.Validate()
	if err != nil {
		return
	}

	if cl.Fee != nil {
		if !(cl.Fee.Token.Type == types.TokenTypeNative || cl.Fee.Token.Type == types.TokenTypeERC20) {
			return errors.New("fee token type is invalid, should be one of: native, ERC20")
		}
		err = cl.Fee.Token.Validate()
		if err != nil {
			return
		}
		if cl.Fee.Token.ChainId != cl.Token.ChainId {
			return errors.New("fee token chain id is invalid")
		}
	}

	if cl.Source == types.CLSourceUndefined {
		cl.Source = types.CLSourceP2P
	}

	if cl.Amount.Cmp(big.NewInt(0)) == 0 {
		return errors.New("amount is required")
	}

	if cl.Source == types.CLSourceD {
		if cl.EscrowAddress == nil {
			return errors.New("escrow address is required for source D")
		}
		// TODO validate EscrowAddress - defineIfEscrowAddressIsCorrect
	}

	return
}
