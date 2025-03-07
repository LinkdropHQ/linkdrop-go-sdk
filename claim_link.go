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
	"strconv"
	"strings"
)

type ClaimLink struct {
	SDK *SDK

	LinkKey ecdsa.PrivateKey

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

type ClaimLinkParams struct {
	SDK           *SDK
	Token         types.Token
	Sender        common.Address
	Amount        *big.Int
	EscrowAddress common.Address
	Expiration    int64
}

func NewClaimLink(
	params *ClaimLinkParams,
	randomBytesCallback types.RandomBytesCallback,
) (claimLink *ClaimLink, err error) {
	linkKey, err := helpers.PrivateKey(randomBytesCallback)
	if err != nil {
		return
	}
	return NewClaimLinkWithLinkKey(params, *linkKey)
}

func NewClaimLinkWithLinkKey(
	params *ClaimLinkParams,
	linkKey ecdsa.PrivateKey,
) (claimLink *ClaimLink, err error) {
	// Validating params
	// SDK should be passed to access LinkDrop API via Client
	if params.SDK == nil {
		return nil, errors.New("params.SDK is required")
	}
	// Token
	err = params.Token.Validate()
	if err != nil {
		return
	}
	// Transfer Id
	transferId, err := helpers.AddressFromPrivateKey(&linkKey)
	if err != nil {
		return
	}
	// Fee
	fee, totalAmount, err := params.SDK.GetCurrentFee(
		params.Token,
		params.Sender,
		transferId,
		params.Expiration,
		params.Amount,
	)
	if err != nil {
		return nil, err
	}

	claimLink = &ClaimLink{
		SDK: params.SDK,

		LinkKey: linkKey,

		Token:  params.Token,
		Amount: params.Amount,
		Sender: params.Sender,

		Fee:         *fee,
		TotalAmount: totalAmount,

		EscrowAddress: params.EscrowAddress,
		Expiration:    params.Expiration,
		Status:        types.ClaimLinkStatusCreated,
	}
	return
}

// TODO separate signing
func (cl *ClaimLink) AddMessage(
	message string,
	encryptionKeyLength int64,
	signTypedData types.SignTypedDataCallback,
) (err error) {
	if encryptionKeyLength == 0 {
		encryptionKeyLength = 12
	}
	if cl.Status >= types.ClaimLinkStatusDeposited {
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

	transferId, err := helpers.AddressFromPrivateKey(&cl.LinkKey)
	if err != nil {
		return
	}
	cl.Message, err = helpers.EncryptMessage(
		message,
		transferId,
		cl.Token.ChainId,
		encryptionKeyLength,
		cl.SDK.GetRandomBytes,
		signTypedData,
	)
	return
}

func (cl *ClaimLink) Redeem(receiver common.Address) (txHash common.Hash, err error) {
	if receiver == types.ZeroAddress {
		return txHash, errors.New("redeem: receiver is not valid")
	}

	receiverSig, err := helpers.GenerateReceiverSig(&cl.LinkKey, receiver)
	if err != nil {
		return
	}

	transferId, err := helpers.AddressFromPrivateKey(&cl.LinkKey)
	if err != nil {
		return
	}
	bApiResp, err := cl.SDK.Client.RedeemLink(
		transferId,
		cl.Token,
		cl.Sender,
		receiver,
		cl.EscrowAddress,
		receiverSig,
	)
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
	transferId, err := helpers.AddressFromPrivateKey(&cl.LinkKey)
	if err != nil {
		return
	}
	linkB, err := cl.SDK.Client.GetTransferStatus(cl.Token.ChainId, transferId)
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

// TODO Separate signing logic
func (cl *ClaimLink) DecryptSenderMessage(
	signTypedData types.SignTypedDataCallback,
) (message string, err error) {
	if cl.Message == nil {
		return "", nil
	}
	transferId, err := helpers.AddressFromPrivateKey(&cl.LinkKey)
	if err != nil {
		return
	}
	if cl.Message.EncryptionKey == [crypto.KeyLength]byte{} {
		messageInitialKey, err := helpers.MessageInitialKeyCreate(
			transferId,
			cl.Token.ChainId,
			signTypedData,
		)
		if err != nil {
			return "", err
		}
		cl.Message.InitialKey = messageInitialKey
		cl.Message.EncryptionKey, err = messageInitialKey.MessageEncryptionKey(int64(cl.Message.Data[0]))
		if err != nil {
			return "", err
		}
	}
	return helpers.DecryptMessage(cl.Message)
}

func (cl *ClaimLink) GetDepositParams() (params *types.ClaimLinkDepositParams, err error) {
	transferId, err := helpers.AddressFromPrivateKey(&cl.LinkKey)
	if err != nil {
		return
	}

	var messageData []byte
	if cl.Message != nil {
		messageData = cl.Message.Data
	}

	var data []byte
	switch cl.Token.Type {
	case types.TokenTypeNative:
		data, err = constants.EscrowTokenAbi.Pack(
			"depositETH",
			transferId,
			cl.TotalAmount.String(),
			strconv.Itoa(int(cl.Expiration)),
			cl.Fee.Amount.String(),
			cl.Fee.Authorization,
			messageData,
		)
	case types.TokenTypeERC20:
		data, err = constants.EscrowTokenAbi.Pack(
			"deposit",
			cl.Token.Address,
			transferId,
			cl.TotalAmount,
			cl.Expiration,
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
	if cl.Message == nil {
		cl.Message = new(types.EncryptedMessage)
	}

	transaction, err := sendTransaction(big.NewInt(int64(params.ChainId)), params.To, params.Value, params.Data)
	if err != nil {
		return
	}
	return transaction.Hash, cl.DepositRegister(*transaction)
}

func (cl *ClaimLink) DepositRegister(transaction types.Transaction) (err error) {
	transferId, err := helpers.AddressFromPrivateKey(&cl.LinkKey)
	if err != nil {
		return
	}
	_, err = cl.SDK.Client.Deposit(
		cl.Token,
		cl.Sender,
		cl.EscrowAddress,
		transferId,
		cl.Expiration,
		transaction,
		cl.Fee,
		cl.Amount,
		cl.TotalAmount,
		cl.Message.Data,
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
	transferId, err := helpers.AddressFromPrivateKey(&cl.LinkKey)
	if err != nil {
		return
	}
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
	senderSignature, err := signTypedData(helpers.LinkSignatureTypedData(linkKeyId, transferId, *escrowPaymentDomain))
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
	transferId, err := helpers.AddressFromPrivateKey(&cl.LinkKey)
	if err != nil {
		return
	}
	linkParams := types.Link{
		LinkKey:         cl.LinkKey,
		TransferId:      transferId,
		ChainId:         cl.Token.ChainId,
		SenderSignature: senderSignature,
		Sender:          &cl.Sender,
		Message:         cl.Message,
	}
	link = helpers.EncodeLink(
		cl.SDK.Client.config.baseURL,
		linkParams,
	)
	return
}

func (cl *ClaimLink) getFee(amount *big.Int) (fee *types.ClaimLinkFeeData, err error) {
	transferId, err := helpers.AddressFromPrivateKey(&cl.LinkKey)
	if err != nil {
		return
	}
	feeB, err := cl.SDK.Client.GetFee(
		cl.Token,
		cl.Sender,
		transferId,
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
