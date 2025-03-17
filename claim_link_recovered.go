package linkdrop

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"github.com/LinkdropHQ/linkdrop-go-sdk/helpers"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type ClaimLinkRecovered struct {
	SDK        *SDK
	TransferId common.Address // TransferId - a unique public identifier of the claim link

	Sender        common.Address
	Token         types.Token
	EscrowAddress common.Address

	Message *types.EncryptedMessage // Message - an optional encrypted message

	LinkKey         *ecdsa.PrivateKey // LinkKey - a re-generated LinkKey that can be "promoted" with SenderSignature
	SenderSignature []byte            // SenderSignature - a sender's signature that allows the link redemption with LinkKeyId instead of an original LinkKey
}

func (clr *ClaimLinkRecovered) GetTypedData(
	linkKeyId common.Address,
) (*apitypes.TypedData, error) {
	escrowVersion, err := helpers.DefineEscrowVersion(clr.EscrowAddress)
	if err != nil {
		return nil, err
	}
	typedData := helpers.RecoveredLinkTypedData(
		linkKeyId,
		clr.TransferId,
		clr.Token.ChainId,
		escrowVersion,
		clr.EscrowAddress,
	)
	return &typedData, nil
}

// GenerateClaimUrl generates a new recovered claim URL
func (clr *ClaimLinkRecovered) GenerateClaimUrl(
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
	typedData, err := clr.GetTypedData(linkKeyId)
	if err != nil {
		return
	}
	senderSignature, err := signTypedData(*typedData)
	if err != nil {
		return
	}
	return clr.ClaimUrl(*newLinkKey, senderSignature)
}

func (clr *ClaimLinkRecovered) ClaimUrl(
	linkKey ecdsa.PrivateKey,
	senderSignature []byte,
) (link string, err error) {
	clr.SenderSignature = senderSignature
	link = helpers.EncodeLink(
		clr.SDK.config.baseURL,
		types.Link{
			SenderSignature: senderSignature,
			LinkKey:         linkKey,
			TransferId:      clr.TransferId,
			ChainId:         clr.Token.ChainId,
			Message:         clr.Message,
		},
	)
	return
}

func (clr *ClaimLinkRecovered) Redeem(
	receiver common.Address,
) (txHash common.Hash, err error) {
	if receiver == types.ZeroAddress {
		err = errors.New("redeem: receiver is not valid")
		return
	}
	if clr.LinkKey == nil || clr.SenderSignature == nil {
		err = errors.New("can't redeem without linkKeyId and sender signature")
		return
	}
	receiverSig, err := helpers.GenerateReceiverSig(clr.LinkKey, receiver)
	if err != nil {
		return
	}

	bApiResp, err := clr.SDK.Client.RedeemLink(
		clr.TransferId,
		clr.Token,
		clr.Sender,
		receiver,
		clr.EscrowAddress,
		receiverSig,
		clr.SenderSignature,
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
