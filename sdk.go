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

	limits = new(types.TransferLimits)
	err = json.Unmarshal(apiResponse, limits)
	return
}

func (sdk *SDK) initializeClaimLink(
	token types.Token,
	amount *big.Int,
	expiration *big.Int,
	sender common.Address,
	source types.CLSource,
	transferId *common.Address,
) (claimLink *ClaimLink, err error) {
	claimLink = &ClaimLink{
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

func (sdk *SDK) GetClaimLink() {}

func (sdk *SDK) RetrieveClaimLink() {}
