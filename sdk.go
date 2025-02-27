package linkdrop

import (
	"encoding/json"
	"errors"
	"github.com/LinkdropHQ/linkdrop-go-sdk/internal/helpers"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type SenderHistory struct {
	ClaimLinks []ClaimLink     `json:"claimLinks"`
	ResultSet  types.ResultSet `json:"resultSet"`
}

type SDK struct {
	Client         *Client
	Deployment     types.Deployment
	GetRandomBytes types.RandomBytesCallback
}

func Init(baseUrl string, deployment types.Deployment, getRandomBytes types.RandomBytesCallback, opts ...Option) (*SDK, error) {
	if baseUrl == "" {
		return nil, errors.New("baseUrl is required")
	}
	if deployment != types.DeploymentLD && deployment != types.DeploymentCBW {
		return nil, errors.New("deployment is invalid, should be one of: LD, CBW")
	}
	if getRandomBytes == nil {
		return nil, errors.New("getRandomBytes is required")
	}

	err := helpers.LoadABI()
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	cfg.applyDefaults()
	for _, opt := range opts {
		opt(cfg)
	}
	cfg.baseURL = baseUrl

	return &SDK{
		Client:         &Client{config: cfg},
		Deployment:     deployment,
		GetRandomBytes: getRandomBytes,
	}, nil
}

func (sdk *SDK) GetVersionFromClaimUrl(claimUrl string) (string, error) {
	return helpers.VersionFromClaimUrl(claimUrl)
}

func (sdk *SDK) GetVersionFromEscrowContract(escrowAddress common.Address) (string, error) {
	return helpers.DefineEscrowVersion(escrowAddress)
}

func (sdk *SDK) GetLinkSourceFromClaimUrl(claimUrl string) (types.CLSource, error) {
	return helpers.LinkSourceFromClaimUrl(claimUrl)
}

func (sdk *SDK) CreateClaimLink(
	token types.Token,
	amount *big.Int,
	sender common.Address,
	expiration *big.Int,
) (claimLink *ClaimLink, err error) {
	err = token.Validate()
	if err != nil {
		return
	}

	if amount.Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("claim link requires amount")
	}

	return sdk.initializeClaimLink(
		token,
		amount,
		expiration,
		sender,
		types.CLSourceP2P,
		nil,
	)
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
	if !apiResponseModel.Success { // Will be empty string since success is a bool field
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

func (sdk *SDK) initializeClaimLink(
	token types.Token,
	amount *big.Int,
	expiration *big.Int,
	sender common.Address,
	source types.CLSource,
	transferId *common.Address,
) (claimLink *ClaimLink, err error) {
	if transferId == nil {
		pk, err := helpers.PrivateKey(sdk.GetRandomBytes)
		if err != nil {
			return nil, err
		}
		address, err := helpers.AddressFromPrivateKey(pk)
		if err != nil {
			return nil, err
		}
		transferId = &address
	}
	claimLink = &ClaimLink{
		SDK:        sdk,
		Token:      token,
		Amount:     amount,
		Expiration: expiration,
		Sender:     sender,
		Source:     source,
		TransferId: *transferId,
	}
	err = claimLink.Validate()
	return
}

func (sdk *SDK) GetCurrentFee(
	token types.Token,
	sender common.Address,
	transferId common.Address,
	expiration *big.Int,
	amount *big.Int,
) (fee *big.Int, err error) {
	feeB, err := sdk.Client.GetFee(
		token,
		sender,
		transferId,
		expiration,
		amount,
	)
	getFeeResp := &struct {
		FeeAmount            []byte `json:"fee_amount"`
		TotalAmount          []byte `json:"total_amount"`
		FeeAuthorization     []byte `json:"fee_authorization"`
		FeeToken             []byte `json:"fee_token"`
		PendingTxs           []byte `json:"pending_txs"`
		PendingBlocks        []byte `json:"pending_blocks"`
		PendingTxSubmittedBn []byte `json:"pending_tx_submitted_bn"`
		PendingTxSubmittedAt []byte `json:"pending_tx_submitted_at"`
		MinTransferAmount    []byte `json:"min_transfer_amount"`
		MaxTransferAmount    []byte `json:"max_transfer_amount"`
	}{}
	err = json.Unmarshal(feeB, getFeeResp)
	return &big.Int{}, err
}

func (sdk *SDK) GetClaimLink(claimUrl string) (claimLink *ClaimLink, err error) {
	linkSource, err := sdk.GetLinkSourceFromClaimUrl(claimUrl)
	if err != nil {
		return
	}
	if linkSource == types.CLSourceD {
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
	tokenId, ok := new(big.Int).SetString(cl["token_id"].(string), 10)
	if !ok {
		return nil, errors.New("invalid token_id")
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
		SDK:     sdk,
		Sender:  common.HexToAddress(cl["sender"].(string)),
		ChainId: types.ChainId(int64(cl["chain_id"].(float64))),
		Token: types.Token{
			Type:    types.TokenType(cl["token_type"].(string)),
			ChainId: types.ChainId(int64(cl["chain_id"].(float64))),
			Address: common.HexToAddress(cl["token"].(string)),
			Id:      tokenId,
		},
		Amount:      amount,
		TotalAmount: totalAmount,
		Fee: &types.CLFee{
			Token: types.Token{
				Type:    feeTokenType,
				ChainId: types.ChainId(int64(cl["chain_id"].(float64))),
				Address: common.HexToAddress(cl["sender"].(string)),
			},
			Amount: feeAmount,
		},
		Expiration:    big.NewInt(int64(cl["expiration"].(float64))),
		TransferId:    common.HexToAddress(cl["transfer_id"].(string)),
		EscrowAddress: common.HexToAddress(cl["escrow"].(string)),
		Operations:    nil, // TODO
		LinkKey:       decodedLink.LinkKey,
		ClaimUrl:      &claimUrl,
		ForRecipient:  false,
		Status:        types.ClItemStatusFromString(cl["status"].(string)),
		Source:        types.CLSource(cl["escrow"].(string)),
	}
	err = claimLink.Validate()
	return
}

func (sdk *SDK) RetrieveClaimLink() {}
