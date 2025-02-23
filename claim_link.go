package linkdrop

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"github.com/LinkdropHQ/linkdrop-go-sdk/internal/constants"
	"github.com/LinkdropHQ/linkdrop-go-sdk/internal/crypto"
	"github.com/LinkdropHQ/linkdrop-go-sdk/internal/helpers"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type ClaimLink struct {
	SDK                    *SDK
	Sender                 common.Address // Required
	Token                  types.Token
	Amount                 *big.Int
	TotalAmount            *big.Int
	PendingTxs             *int64
	PendingTxSubmittedBn   *int64
	PendingTxSubmittedAt   *int64
	PendingBlocks          *int64
	Fee                    *types.CLFee
	Expiration             *big.Int
	TransferId             common.Address
	EscrowAddress          *common.Address
	Operations             []types.CLOperation
	LinkKey                *ecdsa.PrivateKey
	ClaimUrl               *string
	ForRecipient           bool
	Status                 types.CLItemStatus // CLItemStatusUndefined by default
	Source                 types.CLSource     // CLSourceUndefined by default
	EncryptedSenderMessage *types.EncryptedMessage
	validated              bool
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
	} else {
		cl.Fee = new(types.CLFee)
	}

	if cl.Source == types.CLSourceUndefined {
		cl.Source = types.CLSourceP2P
	}

	if cl.Amount.Cmp(big.NewInt(0)) == 0 {
		return errors.New("amount is required")
	}

	if cl.Source == types.CLSourceP2P {
		if cl.EscrowAddress == nil {
			return errors.New("escrow address is required for source P2P")
		}
		// TODO validate EscrowAddress - defineIfEscrowAddressIsCorrect
	}

	cl.validated = true
	return
}

func (cl *ClaimLink) AddMessage(
	message string,
	encryptionKeyLength int64,
	signTypedData types.SignTypedDataCallback,
) (err error) {
	if !cl.validated {
		return errors.New("claim link is not validated. Run Validate()")
	}
	if encryptionKeyLength == 0 {
		encryptionKeyLength = 12
	}
	if cl.Status >= types.CLItemStatusDeposited {
		return errors.New("cannot add message after deposit")
	}
	if len(message) == 0 {
		return errors.New("message text is required")
	}
	if signTypedData == nil {
		return errors.New("signTypedData callback is required")
	}
	if int64(len(message)) > cl.SDK.Client.config.messageConfig.MaxTextLength {
		return errors.New("message text length is too long")
	}
	// TODO move config to SDK?
	if encryptionKeyLength > cl.SDK.Client.config.messageConfig.MaxEncryptionKeyLength ||
		encryptionKeyLength < cl.SDK.Client.config.messageConfig.MinEncryptionKeyLength {
		return errors.New("message text length is too long")
	}

	cl.EncryptedSenderMessage, err = helpers.EncryptMessage(
		message,
		cl.TransferId,
		cl.Token.ChainId,
		encryptionKeyLength,
		cl.SDK.GetRandomBytes,
		signTypedData,
	)

	return
}

func (cl *ClaimLink) GetDepositParams() (params *types.CLDepositParams, err error) {
	if !cl.validated {
		return nil, errors.New("claim link is not validated. Run Validate()")
	}
	if cl.EscrowAddress == nil {
		return nil, errors.New("escrowAddress is not defined for claim link")
	}
	if *cl.EscrowAddress == types.ZeroAddress {
		return nil, errors.New("escrowAddress can't be zero address")
	}

	var data []byte
	senderMessage := ""
	if cl.EncryptedSenderMessage != nil {
		senderMessage = "0x" + cl.EncryptedSenderMessage.Message
	}

	switch cl.Token.Type {
	case types.TokenTypeNative:
		data, err = constants.EscrowTokenAbi.Pack(
			"depositETH",
			cl.TransferId,
			cl.TotalAmount.String(),
			cl.Expiration, // TODO check expiration
			cl.Fee.Amount.String(),
			cl.Fee.Authorization,
			senderMessage,
		)
	case types.TokenTypeERC20:
		data, err = constants.EscrowTokenAbi.Pack(
			"deposit",
			cl.Token.Address.Hex(),
			cl.TransferId,
			cl.TotalAmount.String(),
			cl.Expiration, // TODO check expiration
			cl.Fee.Token.Address.Hex(),
			cl.Fee.Amount.String(),
			cl.Fee.Authorization,
			senderMessage,
		)
	case types.TokenTypeERC721:
		data, err = constants.EscrowTokenAbi.Pack(
			"depositERC721",
			cl.Token.Address.Hex(),
			cl.TransferId,
			cl.Token.Id.String(),
			cl.Expiration, // TODO check expiration
			cl.Fee.Amount.String(),
			cl.Fee.Authorization,
			senderMessage,
		)
	case types.TokenTypeERC1155:
		data, err = constants.EscrowTokenAbi.Pack(
			cl.Token.Address.Hex(),
			cl.TransferId,
			cl.Token.Id.String(),
			cl.Amount,
			cl.Expiration, // TODO check expiration
			cl.Fee.Amount.String(),
			cl.Fee.Authorization,
			senderMessage,
		)
	default:
		return nil, errors.New("invalid token type")
	}

	// Safely skip error check (if cl.Fee == nil) here
	value, err := cl.defineValue()
	if err != nil {
		return nil, err
	}

	return &types.CLDepositParams{
		Value: value,
		Data:  data,
		To:    *cl.EscrowAddress,
	}, nil
}

func (cl *ClaimLink) Redeem(receiver common.Address) (txHash common.Hash, err error) {
	if !cl.validated {
		return txHash, errors.New("claim link is not validated. Run Validate()")
	}
	if receiver == types.ZeroAddress {
		return txHash, errors.New("dest is not valid")
	}
	if cl.ClaimUrl == nil {
		return txHash, errors.New("cannot redeem before deposit")
	}

	if cl.Source == types.CLSourceD {
		linkKeyHex := crypto.Keccak256(helpers.GetClaimCodeFromDashboardLink(*cl.ClaimUrl))
		linkKey, err := helpers.PrivateKeyFromHex(linkKeyHex)
		receiverSignature, err := helpers.GenerateReceiverSig(linkKey, receiver)
		if err != nil {
			return txHash, err
		}
		bTxHash, err := cl.SDK.Client.RedeemLink(
			receiver,
			cl.TransferId,
			receiverSignature,
			nil, nil, nil,
		)
		return common.BytesToHash(bTxHash), nil
	}
	if cl.EscrowAddress == nil {
		return txHash, errors.New("cannot redeem before deposit")
	}
	decodedLink, err := helpers.DecodeLink(*cl.ClaimUrl)
	if err != nil {
		return txHash, err
	}
	receiverSig, err := helpers.GenerateReceiverSig(decodedLink.LinkKey, receiver)
	if err != nil {
		return txHash, err
	}
	if decodedLink.SenderSig != "" {
		bTxHash, err := cl.SDK.Client.RedeemRecoveredLink(
			receiver,
			cl.TransferId,
			receiverSig,
			decodedLink.SenderSig,
			cl.Sender,
			*cl.EscrowAddress,
			cl.Token,
		)
		if err != nil {
			return txHash, err
		}
		return common.BytesToHash(bTxHash), nil
	}
	bTxHash, err := cl.SDK.Client.RedeemLink(
		receiver,
		cl.TransferId,
		receiverSig,
		&cl.Sender,
		cl.EscrowAddress,
		&cl.Token,
	)
	return common.BytesToHash(bTxHash), nil
}

func (cl *ClaimLink) GetStatus() (types.CLItemStatus, []types.CLOperation, error) {
	if !cl.validated {
		return types.CLItemStatusUndefined, nil, errors.New("claim link is not validated. Run Validate()")
	}
	if cl.TransferId == types.ZeroAddress {
		return types.CLItemStatusUndefined, nil, errors.New("transfer id is not defined for claim link")
	}
	linkB, err := cl.SDK.Client.GetTransferStatus(cl.TransferId)
	claimLink := struct {
		Status     types.CLItemStatus  `json:"status"`
		Operations []types.CLOperation `json:"operations"`
	}{}
	err = json.Unmarshal(linkB, &claimLink)
	if err != nil {
		return types.CLItemStatusUndefined, nil, err
	}
	if claimLink.Status != cl.Status {
		cl.Status = claimLink.Status
		cl.Operations = claimLink.Operations
	}
	return claimLink.Status, claimLink.Operations, nil
}

func (cl *ClaimLink) DecryptSenderMessage() (message string, err error) {
	if !cl.validated {
		return "", errors.New("claim link is not validated. Run Validate()")
	}
	if cl.ForRecipient {
		return "", errors.New("this link can only be redeemed")
	}
	if cl.EncryptedSenderMessage == nil {
		return "", nil
	}
	return crypto.Decrypt(
		cl.EncryptedSenderMessage.Message,
		cl.EncryptedSenderMessage.EncryptionKey,
	)
}

func (cl *ClaimLink) defineValue() (value *big.Int, err error) {
	if !cl.validated {
		return big.NewInt(0), errors.New("claim link is not validated. Run Validate()")
	}

	if cl.Fee.Token.Address == cl.Token.Address && cl.Token.Address != constants.NativeTokenAddress {
		return big.NewInt(0), nil
	}
	if cl.Token.Address == constants.NativeTokenAddress {
		return cl.TotalAmount, nil
	}

	return cl.Fee.Amount, nil
}

func (cl *ClaimLink) Deposit(sendTransaction types.SendTransactionCallback) (txHash common.Hash, transferId common.Address, claimUrl string, err error) {
	if !cl.validated {
		return txHash, transferId, claimUrl, errors.New("claim link is not validated. Run Validate()")
	}
	if cl.EscrowAddress == nil {
		return txHash, transferId, claimUrl, errors.New("escrow address is not defined for claim link")
	}
	if sendTransaction == nil {
		return txHash, transferId, claimUrl, errors.New("sendTransaction callback is required")
	}
	if cl.ForRecipient {
		return txHash, transferId, claimUrl, errors.New("this link can only be redeemed")
	}
	params, err := cl.GetDepositParams()
	if err != nil {
		return txHash, transferId, claimUrl, err
	}
	if cl.Fee == nil {
		cl.Fee = new(types.CLFee)
	}
	if cl.EncryptedSenderMessage == nil {
		cl.EncryptedSenderMessage = new(types.EncryptedMessage)
	}
	txHash, err = sendTransaction(params.To, params.Value, params.Data)
	if err != nil {
		return txHash, transferId, claimUrl, err
	}
	_, err = cl.SDK.Client.Deposit(
		cl.Token,
		cl.Sender,
		*cl.EscrowAddress,
		cl.TransferId,
		cl.Expiration,
		txHash,
		*cl.Fee,
		cl.Amount,
		cl.TotalAmount,
		cl.EncryptedSenderMessage.Message,
	)
	if err != nil {
		return txHash, transferId, claimUrl, err
	}
	var linkKey *ecdsa.PrivateKey
	if cl.LinkKey != nil {
		linkKey = cl.LinkKey
	}
	linkParams := types.Link{
		LinkKey:       linkKey,
		TransferId:    cl.TransferId,
		ChainId:       cl.Token.ChainId,
		Sender:        &cl.Sender,
		EncryptionKey: &cl.EncryptedSenderMessage.EncryptionKey,
	}

	claimUrl = helpers.EncodeLink(
		cl.SDK.Client.config.apiURL,
		linkParams,
	)
	cl.ClaimUrl = &claimUrl
	cl.Status = types.CLItemStatusDeposited
	return txHash, cl.TransferId, claimUrl, err
}

// TODO
func (cl *ClaimLink) IsDepositWithAuthorizationAvailable() {
}

// TODO
func (cl *ClaimLink) DepositWithAuthorization() {
}

func (cl *ClaimLink) GetCurrentFee() (fee *types.CLFeeData, err error) {
	return cl.getFee(cl.Amount)
}

func (cl *ClaimLink) UpdateAmount(amount *big.Int) (*types.CLFeeData, error) {
	if cl.ForRecipient {
		return nil, errors.New("this link can only be redeemed")
	}

	if cl.Token.Type == types.TokenTypeERC721 {
		return nil, errors.New("can't update amount for ERC721 token")
	}

	feeData, err := cl.getFee(amount)
	if err != nil {
		return nil, err
	}

	if amount.Cmp(feeData.MinTransferAmount) < 0 {
		return nil, errors.New("amount should be greater than " + feeData.MinTransferAmount.String() + "")
	}

	if amount.Cmp(feeData.MaxTransferAmount) > 0 {
		return nil, errors.New("amount should be less than " + feeData.MaxTransferAmount.String() + "")
	}

	if cl.LinkKey == nil {
		status, _, err := cl.GetStatus()
		if err != nil {
			return nil, err
		}

		if status == types.CLItemStatusCreated {
			cl.Amount = amount
			cl.Fee.Amount = feeData.Fee.Amount
			cl.TotalAmount = feeData.TotalAmount
			cl.Fee.Authorization = feeData.Fee.Authorization
			cl.Fee.Token.Address = feeData.Fee.Token.Address

			return feeData, nil
		} else {
			return nil, errors.New("can't update amount for claim link with status " + status.String())
		}
	}

	if cl.Status >= types.CLItemStatusDeposited {
		return nil, errors.New("can't update amount for claim link with status " + cl.Status.String())
	}

	cl.Amount = amount
	cl.TotalAmount = feeData.TotalAmount
	cl.Fee.Amount = feeData.Fee.Amount
	cl.Fee.Authorization = feeData.Fee.Authorization
	cl.Fee.Token.Address = feeData.Fee.Token.Address

	return feeData, nil
}

func (cl *ClaimLink) GenerateClaimUrl(signTypedData types.SignTypedDataCallback) (link string, transferId common.Address, err error) {
	if !cl.validated {
		return "", types.ZeroAddress, errors.New("claim link is not validated. Run Validate()")
	}
	if cl.ForRecipient {
		return "", types.ZeroAddress, errors.New("this link can only be redeemed")
	}
	if signTypedData == nil {
		return "", types.ZeroAddress, errors.New("signTypedData callback is required")
	}
	if cl.TransferId == types.ZeroAddress {
		return "", types.ZeroAddress, errors.New("transfer id is not defined for claim link")
	}
	if cl.EscrowAddress == nil {
		return "", types.ZeroAddress, errors.New("escrow address is not defined for claim link")
	}
	version, err := helpers.DefineEscrowVersion(*cl.EscrowAddress)
	if err != nil {
		return
	}
	escrowPaymentDomain := &types.TypedDataDomain{
		Name:              "LinkdropEscrow",
		Version:           version,
		ChainId:           cl.Token.ChainId,
		VerifyingContract: cl.EscrowAddress,
	}
	linkKey, _, senderSig, err := helpers.GenerateLinkKeyAndSignature(
		signTypedData,
		cl.SDK.GetRandomBytes,
		cl.TransferId,
		*escrowPaymentDomain,
	)
	if err != nil {
		return
	}
	linkParams := types.Link{
		SenderSig:  senderSig,
		LinkKey:    linkKey,
		TransferId: cl.TransferId,
		ChainId:    cl.Token.ChainId,
	}
	if cl.EncryptedSenderMessage != nil && len(cl.EncryptedSenderMessage.Message) >= 2 {
		// TODO double-check
		encryptionKeyLength := int64(common.FromHex(cl.EncryptedSenderMessage.Message[:2])[0])
		_, encryptionKeyLinkParam, err := helpers.CreateMessageEncryptionKey(
			cl.TransferId.Hex(),
			signTypedData,
			cl.Token.ChainId,
			encryptionKeyLength,
		)
		if err != nil {
			return "", types.ZeroAddress, err
		}
		linkParams.EncryptionKeyLinkParam = &encryptionKeyLinkParam
	}
	claimUrl := helpers.EncodeLink(
		cl.SDK.Client.config.apiURL,
		linkParams,
	)
	cl.ClaimUrl = &claimUrl
	return claimUrl, cl.TransferId, nil
}

func (cl *ClaimLink) getFee(amount *big.Int) (fee *types.CLFeeData, err error) {
	if !cl.validated {
		return nil, errors.New("claim link is not validated. Run Validate()")
	}
	if cl.ForRecipient {
		return nil, errors.New("this link can only be redeemed")
	}
	feeB, err := cl.SDK.Client.GetFee(
		cl.Token,
		cl.Sender,
		cl.TransferId,
		cl.Expiration,
		amount,
	)
	if err != nil {
		return nil, err
	}

	feeRaw := struct {
		Amount            string `json:"amount"`
		TotalAmount       string `json:"total_amount"`
		MaxTransferAmount string `json:"max_transfer_amount"`
		MinTransferAmount string `json:"min_transfer_amount"`
		FeeToken          string `json:"fee_token"`
		FeeAmount         string `json:"fee_amount"`
		FeeAuthorization  string `json:"fee_authorization"`
	}{}
	err = json.Unmarshal(feeB, &feeRaw)
	if err != nil {
		return nil, err
	}
	amount, ok := new(big.Int).SetString(feeRaw.Amount, 10)
	if !ok {
		return nil, errors.New("invalid amount")
	}
	totalAmount, ok := new(big.Int).SetString(feeRaw.TotalAmount, 10)
	if !ok {
		return nil, errors.New("invalid totalAmount")
	}
	maxTransferAmount, ok := new(big.Int).SetString(feeRaw.MaxTransferAmount, 10)
	if !ok {
		return nil, errors.New("invalid maxTransferAmount")
	}
	minTransferAmount, ok := new(big.Int).SetString(feeRaw.MinTransferAmount, 10)
	if !ok {
		return nil, errors.New("invalid minTransferAmount")
	}
	feeAmount, ok := new(big.Int).SetString(feeRaw.FeeAmount, 10)
	if !ok {
		return nil, errors.New("invalid fee amount")
	}
	fee = &types.CLFeeData{
		Amount:            amount,
		TotalAmount:       totalAmount,
		MaxTransferAmount: maxTransferAmount,
		MinTransferAmount: minTransferAmount,
		Fee: types.CLFee{
			Token: types.Token{
				Type:    types.TokenType(feeRaw.FeeToken),
				Address: common.HexToAddress(feeRaw.FeeToken),
			},
			Amount:        feeAmount,
			Authorization: feeRaw.FeeAuthorization,
		},
	}
	return fee, nil
}

func (cl *ClaimLink) defineDomain() {}
