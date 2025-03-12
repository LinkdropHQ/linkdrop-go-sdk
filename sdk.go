package linkdrop

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"github.com/LinkdropHQ/linkdrop-go-sdk/constants"
	"github.com/LinkdropHQ/linkdrop-go-sdk/helpers"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strings"
)

type SenderHistory struct {
	ClaimLinks []ClaimLink `json:"claimLinks"`
	ResultSet  struct {
		Total  int64 `json:"total"`
		Limit  int64 `json:"limit"`
		Offset int64 `json:"offset"`
	} `json:"resultSet"`
}

type SDK struct {
	config SDKConfig
	Client *Client
}

func Init(baseUrl string, apiKey string, opts ...Option) (*SDK, error) {
	if baseUrl == "" {
		return nil, errors.New("baseUrl is required")
	}

	err := helpers.LoadABI()
	if err != nil {
		return nil, err
	}

	var sdkConfig SDKConfig
	sdkConfig.applyDefaults()
	clientConfig := &ClientConfig{
		apiKey: apiKey,
		apiURL: constants.ApiURL,
	}
	for _, opt := range opts {
		opt(&sdkConfig, clientConfig)
	}
	sdkConfig.baseURL = baseUrl

	return &SDK{
		config: sdkConfig,
		Client: &Client{
			config: clientConfig,
		},
	}, nil
}

func (sdk *SDK) Environment() string {
	return sdk.config.environment
}

func (sdk *SDK) GetVersionFromClaimUrl(claimUrl string) (string, error) {
	return helpers.VersionFromClaimUrl(claimUrl)
}

func (sdk *SDK) GetVersionFromEscrowContract(escrowAddress common.Address) (string, error) {
	return helpers.DefineEscrowVersion(escrowAddress)
}

// ClaimLink creates a new ClaimLink generating linkKey using randomBytesCallback
func (sdk *SDK) ClaimLink(
	params ClaimLinkCreationParams,
	randomBytesCallback types.RandomBytesCallback,
) (claimLink *ClaimLink, err error) {
	linkKey, err := helpers.PrivateKey(randomBytesCallback)
	if err != nil {
		return
	}
	claimLink = new(ClaimLink)
	return sdk.ClaimLinkWithLinkKey(params, *linkKey)
}

// ClaimLinkWithTransferId creates a new ClaimLink setting with provided transferId
// NOTE: the generated link will be created without linkKey and will lack some functionality
func (sdk *SDK) ClaimLinkWithTransferId(
	params ClaimLinkCreationParams,
	transferId common.Address,
) (claimLink *ClaimLink, err error) {
	claimLink = new(ClaimLink)
	err = claimLink.new(sdk, &params, nil, transferId)
	return
}

// ClaimLinkWithLinkKey creates a new ClaimLink with pre-generated linkKey
func (sdk *SDK) ClaimLinkWithLinkKey(
	params ClaimLinkCreationParams,
	linkKey ecdsa.PrivateKey,
) (claimLink *ClaimLink, err error) {
	transferId, err := helpers.AddressFromPrivateKey(&linkKey)
	if err != nil {
		return
	}
	claimLink = new(ClaimLink)
	err = claimLink.new(sdk, &params, &linkKey, transferId)
	return
}

func (sdk *SDK) GetSenderHistory(
	token types.Token,
	sender common.Address,
	onlyActive bool,
	offset int64,
	limit int64,
) (history *SenderHistory, err error) {
	err = token.Validate()
	if err != nil {
		return
	}

	apiResponse, err := sdk.Client.GetHistory(token, sender, onlyActive, offset, limit)
	if err != nil {
		return
	}

	history = new(SenderHistory)
	err = json.Unmarshal(apiResponse, history)
	return
}

func (sdk *SDK) GetLimits(token types.Token) (limits *types.TransferLimits, err error) {
	err = token.Validate()
	if err != nil {
		return
	}

	if token.Type == types.TokenTypeERC721 || token.Type == types.TokenTypeERC1155 {
		return nil, errors.New("limits are not available for ERC721 and ERC1155 tokens")
	}

	apiResponse, err := sdk.Client.GetLimits(token)
	if err != nil {
		return
	}

	apiResponseModel := struct {
		Success      bool   `json:"success"`
		Error        string `json:"error"`
		MinAmount    string `json:"min_transfer_amount"`
		MaxAmount    string `json:"max_transfer_amount"`
		MinAmountUSD string `json:"min_transfer_amount_usd"`
		MaxAmountUSD string `json:"max_transfer_amount_usd"`
	}{}
	err = json.Unmarshal(apiResponse, &apiResponseModel)
	if !apiResponseModel.Success {
		return nil, errors.New("error fetching limits: " + apiResponseModel.Error)
	}
	minTransferAmount, _ := new(big.Int).SetString(apiResponseModel.MinAmount, 10)
	maxTransferAmount, _ := new(big.Int).SetString(apiResponseModel.MaxAmount, 10)
	minTransferAmountUsd, _ := new(big.Int).SetString(apiResponseModel.MinAmountUSD, 10)
	maxTransferAmountUsd, _ := new(big.Int).SetString(apiResponseModel.MaxAmountUSD, 10)
	return &types.TransferLimits{
		MinAmount:    minTransferAmount,
		MaxAmount:    maxTransferAmount,
		MinAmountUSD: minTransferAmountUsd,
		MaxAmountUSD: maxTransferAmountUsd,
	}, nil
}

func (sdk *SDK) GetCurrentFee(
	token types.Token,
	sender common.Address,
	transferId common.Address,
	expiration int64,
	amount *big.Int,
) (fee *types.ClaimLinkFee, totalAmount *big.Int, err error) {
	feeB, err := sdk.Client.GetFee(
		token,
		sender,
		transferId,
		expiration,
		amount,
	)
	if err != nil {
		return
	}
	getFeeResp := &struct {
		Success              bool           `json:"success"`
		Error                string         `json:"string"`
		FeeAmount            string         `json:"fee_amount"`
		TotalAmount          string         `json:"total_amount"`
		FeeAuthorization     string         `json:"fee_authorization"`
		FeeToken             common.Address `json:"fee_token"`
		MinTransferAmount    string         `json:"min_transfer_amount"`
		MaxTransferAmount    string         `json:"max_transfer_amount"`
		MinTransferAmountUsd string         `json:"min_transfer_amount_usd"`
		MaxTransferAmountUsd string         `json:"max_transfer_amount_usd"`
	}{}
	err = json.Unmarshal(feeB, getFeeResp)
	if !getFeeResp.Success {
		return nil, nil, errors.New("error fetching fee: " + getFeeResp.Error)
	}

	feeAmount, _ := (&big.Int{}).SetString(getFeeResp.FeeAmount, 10)
	totalAmount, _ = (&big.Int{}).SetString(getFeeResp.TotalAmount, 10)
	tokenType := types.TokenTypeERC20
	if getFeeResp.FeeToken == types.ZeroAddress {
		tokenType = types.TokenTypeNative
	}
	return &types.ClaimLinkFee{
		Token: types.Token{
			Type:    tokenType,
			ChainId: token.ChainId,
			Address: getFeeResp.FeeToken,
		},
		Amount:        feeAmount,
		Authorization: common.Hex2Bytes(strings.TrimPrefix(getFeeResp.FeeAuthorization, "0x")),
	}, totalAmount, err
}

func (sdk *SDK) GetClaimLink(claimUrl string) (claimLink *ClaimLink, err error) {
	linkSource, err := helpers.LinkSourceFromClaimUrl(claimUrl)
	if err != nil {
		return
	}
	if linkSource == types.LinkSourceDashboard {
		// TODO handle
		return nil, errors.New("not implemented yet")
	}

	decodedLink, err := helpers.DecodeLink(claimUrl)
	if err != nil {
		return
	}

	apiResp, err := sdk.Client.GetTransferStatus(decodedLink.ChainId, decodedLink.TransferId)
	if err != nil {
		return
	}

	respModel := struct {
		ClaimLink map[string]any `json:"claim_link"`
	}{}
	err = json.Unmarshal(apiResp, &respModel)
	if err != nil {
		return
	}

	cl := respModel.ClaimLink
	tokenType := types.TokenType(cl["token_type"].(string))

	tokenId, ok := new(big.Int).SetString(cl["token_id"].(string), 10)
	if !ok {
		return nil, errors.New("invalid token_id")
	}
	if !(tokenType == types.TokenTypeERC721 || tokenType == types.TokenTypeERC1155) {
		tokenId = nil
	}

	amount, ok := new(big.Int).SetString(cl["amount"].(string), 10)
	if !ok {
		return nil, errors.New("invalid amount")
	}
	totalAmount, ok := new(big.Int).SetString(cl["total_amount"].(string), 10)
	if !ok {
		return nil, errors.New("invalid total_amount")
	}
	feeAmount, ok := new(big.Int).SetString(cl["fee_amount"].(string), 10)
	if !ok {
		return nil, errors.New("invalid total_amount")
	}
	feeTokenAddress := common.HexToAddress(cl["sender"].(string))
	feeTokenType := types.TokenTypeERC20
	if feeTokenAddress == types.ZeroAddress {
		feeTokenType = types.TokenTypeNative
	}

	claimLink = &ClaimLink{
		SDK:    sdk,
		Sender: common.HexToAddress(cl["sender"].(string)),
		Token: types.Token{
			Type:    tokenType,
			ChainId: types.ChainId(int64(cl["chain_id"].(float64))),
			Address: common.HexToAddress(cl["token"].(string)),
			Id:      tokenId,
		},
		Amount:      amount,
		TotalAmount: totalAmount,
		Fee: types.ClaimLinkFee{
			Token: types.Token{
				Type:    feeTokenType,
				ChainId: types.ChainId(int64(cl["chain_id"].(float64))),
				Address: common.HexToAddress(cl["sender"].(string)),
			},
			Amount: feeAmount,
		},
		Expiration:    int64(cl["expiration"].(float64)),
		TransferId:    common.HexToAddress(cl["transfer_id"].(string)),
		EscrowAddress: common.HexToAddress(cl["escrow"].(string)),
		LinkKey:       &decodedLink.LinkKey,
		Status:        types.ClaimLinkStatusFromString(cl["status"].(string)),
	}
	return
}
