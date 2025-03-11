package helpers

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/LinkdropHQ/linkdrop-go-sdk/crypto"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func MessageDecrypt(
	message *types.EncryptedMessage,
) (string, error) {
	if message == nil {
		return "", errors.New("message is nil")
	}
	encryptionKey, err := message.LinkKey.MessageEncryptionKey()
	if err != nil {
		return "", err
	}
	// message.Data[2:] - message data with key len removed
	return crypto.Decrypt(message.Data[2:], encryptionKey)
}

func MessageEncrypt(
	message string,
	initialKey types.MessageInitialKey,
	linkKeyLength uint16,
	nonce [crypto.NonceLength]byte,
) (encryptedMessage *types.EncryptedMessage, err error) {
	if linkKeyLength > 0xFFFF {
		return nil, errors.New("linkKeyLength exceeds 2 bytes")
	}
	linkKey := initialKey.LinkKey(linkKeyLength)
	encryptionKey, err := linkKey.MessageEncryptionKey()
	if err != nil {
		return
	}
	encryptedSenderMessage, err := crypto.Encrypt([]byte(message), encryptionKey, nonce)
	if err != nil {
		return nil, err
	}
	linkKeyLengthB := make([]byte, 2)
	binary.BigEndian.PutUint16(linkKeyLengthB, uint16(len(linkKey)))
	encryptedMessage = &types.EncryptedMessage{
		Data:    append(linkKeyLengthB[:], encryptedSenderMessage...),
		LinkKey: linkKey,
	}
	return
}

// MessageInitialKeyCreate creates a message encryption key
// This function returns initial key which is passed as link parameter
// Use MessageInitialKey.EncryptionKey() to retrieve Encryption Key from Initial Key
func MessageInitialKeyCreate(
	transferId common.Address,
	chainId types.ChainId,
	signTypedData types.SignTypedDataCallback,
) (initialKey types.MessageInitialKey, err error) {
	signature, err := signTypedData(MessageInitialKeyTypedData(transferId, chainId))
	if err != nil {
		return
	}
	initialKey = MessageInitialKeyFromSignature(signature)
	return
}

func MessageInitialKeyTypedData(
	transferId common.Address,
	chainId types.ChainId,
) apitypes.TypedData {
	return apitypes.TypedData{
		Domain: apitypes.TypedDataDomain{
			Name:    "MyEncryptionScheme",
			Version: "1",
			ChainId: math.NewHexOrDecimal256(int64(chainId)),
		},
		PrimaryType: "EncryptionMessage",
		Types: map[string][]apitypes.Type{
			"EIP712Domain": {
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
			},
			"EncryptionMessage": {
				{Name: "seed", Type: "string"},
			},
		},
		Message: map[string]interface{}{
			"seed": fmt.Sprintf("Encrypting message (transferId: %s)", transferId.Hex()),
		},
	}
}

func MessageInitialKeyFromSignature(
	MessageEncryptionKeyTypedDataSignature []byte,
) types.MessageInitialKey {
	return sha256.Sum256(MessageEncryptionKeyTypedDataSignature)
}
