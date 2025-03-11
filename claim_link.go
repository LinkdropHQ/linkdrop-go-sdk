package linkdrop

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"github.com/LinkdropHQ/linkdrop-go-sdk/constants"
	"github.com/LinkdropHQ/linkdrop-go-sdk/crypto"
	"github.com/LinkdropHQ/linkdrop-go-sdk/helpers"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"math/big"
	"strings"
)

type ClaimLink struct {
	SDK *SDK

	LinkKey    *ecdsa.PrivateKey
	TransferId common.Address

	Token  types.Token
	Amount *big.Int
	Sender common.Address

	Fee         types.ClaimLinkFee
	TotalAmount *big.Int

	Message *types.EncryptedMessage

	EscrowAddress common.Address
	Expiration    int64
	Operations    []types.ClaimLinkOperation
	Status        types.ClaimLinkStatus
}

type ClaimLinkCreationParams struct {
	Token      types.Token
	Sender     common.Address
	Amount     *big.Int
	Expiration int64
}

func (cl *ClaimLink) newWithLinkKey(
	sdk *SDK,
	params *ClaimLinkCreationParams,
	linkKey ecdsa.PrivateKey,
) (err error) {
	// Transfer Id
	transferId, err := helpers.AddressFromPrivateKey(&linkKey)
	if err != nil {
		return
	}
	return cl.new(sdk, params, &linkKey, transferId)
}

func (cl *ClaimLink) new(
	sdk *SDK,
	params *ClaimLinkCreationParams,
	linkKey *ecdsa.PrivateKey,
	transferId common.Address,
) (err error) {
	// Validating params
	if params == nil {
		return errors.New("params is required")
	}
	// SDK should be passed to access LinkDrop API via Client
	if sdk == nil {
		return errors.New("params.SDK is required")
	}
	// Token
	err = params.Token.Validate()
	if err != nil {
		return
	}
	// Fee
	fee, totalAmount, err := sdk.GetCurrentFee(
		params.Token,
		params.Sender,
		transferId,
		params.Expiration,
		params.Amount,
	)
	if err != nil {
		return
	}

	*cl = ClaimLink{
		SDK: sdk,

		LinkKey:    linkKey,
		TransferId: transferId,

		Token:  params.Token,
		Amount: params.Amount,
		Sender: params.Sender,

		Fee:         *fee,
		TotalAmount: totalAmount,

		EscrowAddress: sdk.config.escrowContractAddress,
		Expiration:    params.Expiration,
		Status:        types.ClaimLinkStatusCreated,
	}
	return
}

func (cl *ClaimLink) AddMessage(
	message string,
	encryptionKeyLength uint16,
	signTypedData types.SignTypedDataCallback,
	nonce [crypto.NonceLength]byte,
) (err error) {
	if signTypedData == nil {
		return errors.New("signTypedData callback is required")
	}
	initialKey, err := helpers.MessageInitialKeyCreate(cl.TransferId, cl.Token.ChainId, signTypedData)
	if err != nil {
		return
	}
	return cl.AddMessageWithInitialKey(message, encryptionKeyLength, initialKey, nonce)
}

func (cl *ClaimLink) AddMessageWithInitialKey(
	message string,
	encryptionKeyLength uint16,
	initialKey types.MessageInitialKey,
	nonce [crypto.NonceLength]byte,
) (err error) {
	if cl.Status >= types.ClaimLinkStatusDeposited {
		return errors.New("cannot add message after deposit")
	}
	if len(message) == 0 {
		return errors.New("message text is required")
	}
	if int64(len(message)) > cl.SDK.config.messageConfig.MaxTextLength {
		return errors.New("message text length is too long")
	}
	if encryptionKeyLength > cl.SDK.config.messageConfig.MaxEncryptionKeyLength ||
		encryptionKeyLength < cl.SDK.config.messageConfig.MinEncryptionKeyLength {
		return errors.New("wrong encryption key length")
	}
	if encryptionKeyLength == 0 {
		encryptionKeyLength = 12
	}
	cl.Message, err = helpers.MessageEncrypt(message, initialKey, encryptionKeyLength, nonce)
	return
}

// AddMessageRaw allows to directly set Encrypted Message for the link
// Verifies the validity of message provided
// NOTE: message.Data should contain message.LinkKey length as the first byte. See helpers.MessageEncrypt
func (cl *ClaimLink) AddMessageRaw(
	originalMessage string,
	message types.EncryptedMessage,
) (err error) {
	if cl.Status >= types.ClaimLinkStatusDeposited {
		return errors.New("cannot add message after deposit")
	}
	decryptedMessage, err := helpers.MessageDecrypt(&message)
	if err != nil || decryptedMessage != originalMessage {
		return errors.New("message is not valid")
	}
	cl.Message = &message
	return
}

func (cl *ClaimLink) Redeem(receiver common.Address) (txHash common.Hash, err error) {
	if receiver == types.ZeroAddress {
		err = errors.New("redeem: receiver is not valid")
		return
	}
	if cl.LinkKey == nil {
		err = errors.New("redeem: can't redeem without linkKey")
		return
	}

	receiverSig, err := helpers.GenerateReceiverSig(cl.LinkKey, receiver)
	if err != nil {
		return
	}

	bApiResp, err := cl.SDK.Client.RedeemLink(
		cl.TransferId,
		cl.Token,
		cl.Sender,
		receiver,
		cl.EscrowAddress,
		receiverSig,
	)
	if err != nil {
		return
	}
	ApiRespModel := struct {
		Success bool   `json:"success"`
		TxHash  string `json:"tx_hash"`
		Error   string `json:"error"`
	}{}
	err = json.Unmarshal(bApiResp, &ApiRespModel)
	if err != nil {
		return
	}
	if !ApiRespModel.Success {
		err = errors.New(ApiRespModel.Error)
		return
	}
	return common.HexToHash(ApiRespModel.TxHash), nil
}

func (cl *ClaimLink) GetStatus() (status types.ClaimLinkStatus, operations []types.ClaimLinkOperation, err error) {
	linkB, err := cl.SDK.Client.GetTransferStatus(cl.Token.ChainId, cl.TransferId)
	claimLink := struct {
		Status     types.ClaimLinkStatus      `json:"status"`
		Operations []types.ClaimLinkOperation `json:"operations"`
	}{}
	err = json.Unmarshal(linkB, &claimLink)
	if err != nil {
		return
	}
	if claimLink.Status != cl.Status {
		cl.Status = claimLink.Status
		cl.Operations = claimLink.Operations
	}
	return claimLink.Status, claimLink.Operations, nil
}

func (cl *ClaimLink) DecryptSenderMessage() (message string, err error) {
	if cl.Message == nil {
		return "", errors.New("message is not set")
	}
	if cl.Message.LinkKey == "" {
		return "", errors.New("message link key is not set")
	}
	return helpers.MessageDecrypt(cl.Message)
}

func (cl *ClaimLink) GetDepositParams() (params *types.ClaimLinkDepositParams, err error) {
	var messageData []byte
	if cl.Message != nil {
		messageData = cl.Message.Data
	}

	var data []byte
	switch cl.Token.Type {
	case types.TokenTypeNative:
		data, err = constants.EscrowTokenAbi.Pack(
			"depositETH",
			cl.TransferId,
			cl.TotalAmount.String(),
			big.NewInt(cl.Expiration),
			cl.Fee.Amount.String(),
			cl.Fee.Authorization,
			messageData,
		)
	case types.TokenTypeERC20:
		data, err = constants.EscrowTokenAbi.Pack(
			"deposit",
			cl.Token.Address,
			cl.TransferId,
			cl.TotalAmount,
			big.NewInt(cl.Expiration),
			cl.Fee.Token.Address,
			cl.Fee.Amount,
			cl.Fee.Authorization,
			messageData,
		)
	default:
		return nil, errors.New("invalid token type")
	}
	if err != nil {
		return nil, err
	}

	value, err := helpers.DefineValue(cl.Token, cl.Fee, cl.TotalAmount)
	if err != nil {
		return nil, err
	}

	return &types.ClaimLinkDepositParams{
		ChainId: cl.Token.ChainId,
		Value:   value,
		Data:    data,
		To:      cl.EscrowAddress,
	}, nil
}

func (cl *ClaimLink) Deposit(sendTransaction types.SendTransactionCallback) (txHash common.Hash, err error) {
	params, err := cl.GetDepositParams()
	if err != nil {
		return
	}

	transaction, err := sendTransaction(big.NewInt(int64(params.ChainId)), params.To, params.Value, params.Data)
	if err != nil {
		return
	}
	return transaction.Hash, cl.DepositRegister(*transaction)
}

func (cl *ClaimLink) DepositRegister(transaction types.Transaction) (err error) {
	var messageData []byte
	if cl.Message != nil {
		messageData = cl.Message.Data
	}
	_, err = cl.SDK.Client.Deposit(
		cl.Token,
		cl.Sender,
		cl.EscrowAddress,
		cl.TransferId,
		cl.Expiration,
		transaction,
		cl.Fee,
		cl.Amount,
		cl.TotalAmount,
		messageData,
	)
	if err != nil {
		return
	}
	cl.Status = types.ClaimLinkStatusDeposited
	return
}

func (cl *ClaimLink) GetCurrentFee() (fee *types.ClaimLinkFeeData, err error) {
	return cl.getFee(cl.Amount)
}

func (cl *ClaimLink) UpdateAmount(amount *big.Int) (err error) {
	if cl.Token.Type == types.TokenTypeERC721 {
		return errors.New("can't update amount for ERC721 token")
	}
	if cl.Status >= types.ClaimLinkStatusDeposited {
		return errors.New("can't update amount for claim link with status " + cl.Status.String())
	}

	feeData, err := cl.getFee(amount)
	if err != nil {
		return
	}

	if amount.Cmp(feeData.MinTransferAmount) < 0 {
		return errors.New("amount should be greater than " + feeData.MinTransferAmount.String() + "")
	}
	if amount.Cmp(feeData.MaxTransferAmount) > 0 {
		return errors.New("amount should be less than " + feeData.MaxTransferAmount.String() + "")
	}

	cl.Amount = amount
	cl.TotalAmount = feeData.TotalAmount
	cl.Fee = feeData.Fee
	return
}

func (cl *ClaimLink) GenerateRecoveredClaimUrl(
	getRandomBytes types.RandomBytesCallback,
	signTypedData types.SignTypedDataCallback,
) (link string, err error) {
	newLinkKey, err := helpers.PrivateKey(getRandomBytes)
	if err != nil {
		return
	}
	linkKeyId, err := helpers.AddressFromPrivateKey(newLinkKey)
	if err != nil {
		return
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
	senderSignature, err := signTypedData(helpers.LinkSignatureTypedData(linkKeyId, cl.TransferId, *escrowPaymentDomain))
	if err != nil {
		return
	}
	return cl.GenerateClaimUrl(senderSignature)
}

// GenerateClaimUrl generates a claim URL.
// If senderSignature is provided, the link will be recovered. Otherwise, the original link for the associated linkKey will be generated.
// It returns the constructed URL as a string or an error if the process fails.
func (cl *ClaimLink) GenerateClaimUrl(
	senderSignature []byte,
) (link string, err error) {
	if cl.LinkKey == nil {
		return "", errors.New("can't generate claim url without linkKey")
	}
	linkParams := types.Link{
		LinkKey:         *cl.LinkKey,
		TransferId:      cl.TransferId,
		ChainId:         cl.Token.ChainId,
		SenderSignature: senderSignature,
		Sender:          &cl.Sender,
		Message:         cl.Message,
	}
	link = helpers.EncodeLink(
		cl.SDK.config.baseURL,
		linkParams,
	)
	return
}

func (cl *ClaimLink) getFee(amount *big.Int) (fee *types.ClaimLinkFeeData, err error) {
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
	fee = &types.ClaimLinkFeeData{
		Amount:            amount,
		TotalAmount:       totalAmount,
		MaxTransferAmount: maxTransferAmount,
		MinTransferAmount: minTransferAmount,
		Fee: types.ClaimLinkFee{
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
