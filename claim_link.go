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
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"math/big"
	"strings"
	"time"
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
	EscrowAddress          common.Address
	Operations             []types.CLOperation
	LinkKey                *ecdsa.PrivateKey
	ClaimUrl               *string
	ForRecipient           bool
	Status                 types.CLItemStatus // CLItemStatusUndefined by default
	Source                 types.CLSource     // CLSourceUndefined by default
	EncryptedSenderMessage *types.EncryptedMessage
	ChainId                types.ChainId
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

	if cl.EscrowAddress == types.ZeroAddress {
		if cl.Token.Type == types.TokenTypeERC1155 || cl.Token.Type == types.TokenTypeERC721 {
			// TODO refactor - move escrow address to SDK config
			cl.EscrowAddress = cl.SDK.Client.config.escrowNFTContractAddress
		} else {
			cl.EscrowAddress = cl.SDK.Client.config.escrowContractAddress
		}
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
	if cl.EncryptedSenderMessage == nil {
		cl.EncryptedSenderMessage = new(types.EncryptedMessage)
	}

	var data []byte

	switch cl.Token.Type {
	case types.TokenTypeNative:
		data, err = constants.EscrowTokenAbi.Pack(
			"depositETH",
			cl.TransferId,
			cl.TotalAmount.String(),
			cl.Expiration.String(),
			cl.Fee.Amount.String(),
			cl.Fee.Authorization,
			cl.EncryptedSenderMessage.Message,
		)
	case types.TokenTypeERC20:
		data, err = constants.EscrowTokenAbi.Pack(
			"deposit",
			cl.Token.Address,
			cl.TransferId,
			cl.TotalAmount,
			cl.Expiration,
			cl.Fee.Token.Address,
			cl.Fee.Amount,
			cl.Fee.Authorization,
			cl.EncryptedSenderMessage.Message,
		)
	case types.TokenTypeERC721:
		data, err = constants.EscrowTokenAbi.Pack(
			"depositERC721",
			cl.Token.Address.Hex(),
			cl.TransferId,
			cl.Token.Id.String(),
			cl.Expiration.String(),
			cl.Fee.Amount.String(),
			cl.Fee.Authorization,
			cl.EncryptedSenderMessage.Message,
		)
	case types.TokenTypeERC1155:
		data, err = constants.EscrowTokenAbi.Pack(
			cl.Token.Address.Hex(),
			cl.TransferId,
			cl.Token.Id.String(),
			cl.Amount,
			cl.Expiration.String(),
			cl.Fee.Amount.String(),
			cl.Fee.Authorization,
			cl.EncryptedSenderMessage.Message,
		)
	default:
		return nil, errors.New("invalid token type")
	}
	if err != nil {
		return nil, err
	}

	value, err := cl.defineValue()
	if err != nil {
		return nil, err
	}

	return &types.CLDepositParams{
		ChainId: cl.Token.ChainId,
		Value:   value,
		Data:    data,
		To:      cl.EscrowAddress,
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
			cl.ChainId,
			receiver,
			cl.TransferId,
			receiverSignature,
			nil, nil, nil,
		)
		return common.BytesToHash(bTxHash), nil
	}
	decodedLink, err := helpers.DecodeLink(*cl.ClaimUrl)
	if err != nil {
		return txHash, err
	}
	receiverSig, err := helpers.GenerateReceiverSig(decodedLink.LinkKey, receiver)
	if err != nil {
		return txHash, err
	}
	if decodedLink.SenderSig != nil {
		bTxHash, err := cl.SDK.Client.RedeemRecoveredLink(
			receiver,
			cl.TransferId,
			receiverSig,
			decodedLink.SenderSig,
			cl.Sender,
			cl.EscrowAddress,
			cl.Token,
		)
		if err != nil {
			return txHash, err
		}
		return common.BytesToHash(bTxHash), nil
	}
	bApiResp, err := cl.SDK.Client.RedeemLink(
		cl.ChainId,
		receiver,
		cl.TransferId,
		receiverSig,
		&cl.Sender,
		&cl.EscrowAddress,
		&cl.Token,
	)
	ApiRespModel := struct {
		Success bool   `json:"success"`
		TxHash  string `json:"tx_hash"`
		Error   string `json:"error"`
	}{}
	err = json.Unmarshal(bApiResp, &ApiRespModel)
	if err != nil {
		return common.HexToHash(ApiRespModel.TxHash), err
	}
	if !ApiRespModel.Success {
		return common.HexToHash(ApiRespModel.TxHash), errors.New(ApiRespModel.Error)
	}
	return common.HexToHash(ApiRespModel.TxHash), nil
}

func (cl *ClaimLink) GetStatus() (types.CLItemStatus, []types.CLOperation, error) {
	if !cl.validated {
		return types.CLItemStatusUndefined, nil, errors.New("claim link is not validated. Run Validate()")
	}
	if cl.TransferId == types.ZeroAddress {
		return types.CLItemStatusUndefined, nil, errors.New("transfer id is not defined for claim link")
	}
	linkB, err := cl.SDK.Client.GetTransferStatus(cl.Token.ChainId, cl.TransferId)
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

func (cl *ClaimLink) DecryptSenderMessage(
	signTypedData types.SignTypedDataCallback,
) (message string, err error) {
	if !cl.validated {
		return "", errors.New("claim link is not validated. Run Validate()")
	}
	if cl.ForRecipient {
		return "", errors.New("this link can only be redeemed")
	}
	if cl.EncryptedSenderMessage == nil {
		return "", nil
	}
	if cl.EncryptedSenderMessage.EncryptionKey == [crypto.KeyLength]byte{} {
		encryptionKeyLength := int64(cl.EncryptedSenderMessage.Message[0])
		cl.EncryptedSenderMessage.EncryptionKey, _, err = helpers.CreateMessageEncryptionKey(
			cl.TransferId.Hex(),
			signTypedData,
			cl.Token.ChainId,
			encryptionKeyLength,
		)
		if err != nil {
			return "", err
		}
	}
	return helpers.DecryptMessage(cl.EncryptedSenderMessage)
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

func (cl *ClaimLink) Deposit(sendTransaction types.SendTransactionCallback) (txHash common.Hash, transferId common.Address, err error) {
	if !cl.validated {
		return txHash, transferId, errors.New("claim link is not validated. Run Validate()")
	}
	if sendTransaction == nil {
		err = errors.New("sendTransaction callback is required")
		return
	}
	if cl.ForRecipient {
		return txHash, transferId, errors.New("this link can only be redeemed")
	}

	params, err := cl.GetDepositParams()
	if err != nil {
		return
	}
	if cl.Fee == nil {
		cl.Fee = new(types.CLFee)
	}
	if cl.EncryptedSenderMessage == nil {
		cl.EncryptedSenderMessage = new(types.EncryptedMessage)
	}

	transaction, err := sendTransaction(big.NewInt(int64(params.ChainId)), params.To, params.Value, params.Data)
	if err != nil {
		return
	}

	_, err = cl.SDK.Client.Deposit(
		cl.Token,
		cl.Sender,
		cl.EscrowAddress,
		cl.TransferId,
		cl.Expiration,
		transaction,
		*cl.Fee,
		cl.Amount,
		cl.TotalAmount,
		cl.EncryptedSenderMessage.Message,
	)
	if err != nil {
		return
	}
	cl.Status = types.CLItemStatusDeposited
	return transaction.Hash, cl.TransferId, err
}

func (cl *ClaimLink) IsDepositWithAuthorizationAvailable(address common.Address) bool {
	return constants.SupportedStableCoins[address] != constants.SelectorUndefined
}

func (cl *ClaimLink) DepositWithAuthorization(
	signTypedData types.SignTypedDataCallback,
	authConfig *types.AuthorizationConfig,
) (txHash common.Hash, transferId common.Address, err error) {
	if !cl.validated {
		err = errors.New("claim link is not validated. Run Validate()")
		return
	}
	if signTypedData == nil {
		err = errors.New("signTypedData callback is required")
		return
	}
	if cl.Expiration == nil {
		err = errors.New("expiration is not defined for claim link")
		return
	}
	if cl.Amount == nil {
		err = errors.New("amount is not defined for claim link")
		return
	}
	if cl.Status >= types.CLItemStatusDeposited {
		err = errors.New("can't deposit with authorization for claim link with status " + cl.Status.String())
		return
	}
	if cl.Token.Type == types.TokenTypeNative {
		err = errors.New("can't deposit with authorization for native token")
		return
	}

	var authSelector string
	if authConfig == nil {
		domain, err := helpers.DefineDomain(cl.Token)
		if err != nil {
			return txHash, transferId, err
		}
		authConfig = &types.AuthorizationConfig{
			Domain: domain,
		}
		authSelector = string(constants.SupportedStableCoins[cl.Token.Address])
	} else {
		if authConfig.AuthorizationMethod == nil {
			err = errors.New("authorization method is not defined for authorization config")
			return
		}
		authSelector, err = authConfig.AuthorizationMethod.Selector()
		if err != nil {
			return
		}
	}
	now := time.Now().Unix()
	validAfter := now - 60*60
	validBefore := now + 60*60*24
	authorization, err := helpers.GetDepositAuthorization(
		signTypedData,
		cl.Sender,
		cl.EscrowAddress,
		cl.Amount,
		validAfter,
		validBefore,
		cl.TransferId,
		cl.Expiration,
		authConfig.Domain,
		cl.Token,
		cl.Fee.Amount,
		constants.Selector(authSelector),
		authConfig.AuthorizationMethod,
	)
	if err != nil {
		return
	}

	var message []byte
	if cl.EncryptedSenderMessage != nil {
		message = cl.EncryptedSenderMessage.Message // signing key len included [2:] to remove
	}
	result, err := cl.SDK.Client.DepositWithAuthorization(
		cl.Token,
		cl.Sender,
		cl.EscrowAddress,
		cl.TransferId,
		cl.Expiration,
		authorization,
		authSelector,
		*cl.Fee,
		cl.Amount,
		cl.TotalAmount,
		message,
	)
	if err != nil {
		return
	}

	cl.Status = types.CLItemStatusDeposited
	return common.BytesToHash(result), cl.TransferId, err
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

func (cl *ClaimLink) GenerateRecoveredClaimUrl(
	signTypedData types.SignTypedDataCallback,
) (link string, transferId common.Address, err error) {
	return cl.generateClaimUrl(true, signTypedData)
}

func (cl *ClaimLink) GenerateClaimUrl(signTypedData types.SignTypedDataCallback) (link string, transferId common.Address, err error) {
	return cl.generateClaimUrl(false, signTypedData)
}

func (cl *ClaimLink) generateClaimUrl(
	linkRecover bool,
	signTypedData types.SignTypedDataCallback,
) (link string, transferId common.Address, err error) {
	if signTypedData == nil {
		return "", types.ZeroAddress, errors.New("signTypedData callback is required")
	}
	if !cl.validated {
		return "", types.ZeroAddress, errors.New("claim link is not validated. Run Validate()")
	}
	if cl.ForRecipient {
		return "", types.ZeroAddress, errors.New("this link can only be redeemed")
	}
	if cl.TransferId == types.ZeroAddress {
		return "", types.ZeroAddress, errors.New("transfer id is not defined for claim link")
	}

	version, err := helpers.DefineEscrowVersion(cl.EscrowAddress)
	if err != nil {
		return
	}
	escrowPaymentDomain := &apitypes.TypedDataDomain{
		Name:              "LinkdropEscrow",
		Version:           version,
		ChainId:           math.NewHexOrDecimal256(int64(cl.Token.ChainId)),
		VerifyingContract: cl.EscrowAddress.Hex(),
	}

	var linkKeyId common.Address
	var senderSig []byte
	// If link key doesn't exist it's a new link or a recovered one
	// recover == true forces re-generation
	if cl.LinkKey == nil || linkRecover {
		// Just a new link
		cl.LinkKey, linkKeyId, err = helpers.GenerateLinkKey(
			cl.SDK.GetRandomBytes,
		)
		if err != nil {
			return
		}
	}
	// If recover - the link key was re-generated so we include sig
	if linkRecover {
		senderSig, err = helpers.GenerateLinkSignature(
			linkKeyId,
			signTypedData,
			cl.TransferId,
			*escrowPaymentDomain,
		)
	}

	linkParams := types.Link{
		LinkKey:    cl.LinkKey,
		TransferId: cl.TransferId,
		ChainId:    cl.Token.ChainId,
		SenderSig:  senderSig,
		Sender:     &cl.Sender,
	}
	if cl.EncryptedSenderMessage != nil && len(cl.EncryptedSenderMessage.Message) >= 2 {
		_, encryptionKeyLinkParam, err := helpers.CreateMessageEncryptionKey(
			cl.TransferId.Hex(),
			signTypedData,
			cl.Token.ChainId,
			int64(len(cl.EncryptedSenderMessage.EncryptionKey)),
		)
		if err != nil {
			return "", types.ZeroAddress, err
		}
		linkParams.EncryptionKeyLinkParam = &encryptionKeyLinkParam
	}
	claimUrl := helpers.EncodeLink(
		cl.SDK.Client.config.baseURL,
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
			Authorization: common.Hex2Bytes(strings.TrimPrefix(feeRaw.FeeAuthorization, "0x")),
		},
	}
	return fee, nil
}
