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
	"github.com/mr-tron/base58"
)

type MessageInitialKey [32]byte
type MessageEncryptionKey [32]byte

func DecryptMessage(
	message *types.EncryptedMessage,
) (string, error) {
	encoded := message.Data[2:] // Remove length from message
	return crypto.Decrypt(encoded, message.EncryptionKey)
}

func EncryptMessage(
	message string,
	transferID common.Address,
	chainID types.ChainId,
	encryptionKeyLength int64,
	getRandomBytes types.RandomBytesCallback,
	signTypedData types.SignTypedDataCallback,
) (encryptedMessage *types.EncryptedMessage, err error) {
	if encryptionKeyLength > 0xFFFF {
		return nil, errors.New("encryptionKeyLength exceeds 2 bytes")
	}

	initialKey, err := MessageInitialKeyCreate(transferID, chainID, signTypedData)
	if err != nil {
		return
	}

	encryptionKey, err := initialKey.MessageEncryptionKey(encryptionKeyLength)

	encryptedSenderMessage, err := crypto.Encrypt([]byte(message), encryptionKey, [24]byte{}, getRandomBytes)
	if err != nil {
		return nil, err
	}

	encryptionKeyLengthB := make([]byte, 2)
	binary.BigEndian.PutUint16(encryptionKeyLengthB, uint16(encryptionKeyLength))

	// Build the result
	result := &types.EncryptedMessage{
		Data:          append(encryptionKeyLengthB[:], encryptedSenderMessage...),
		EncryptionKey: encryptionKey,
	}
	return result, nil
}

// MessageInitialKeyCreate creates a message encryption key
// This function returns initial key which is passed as link parameter
// Use MessageInitialKey.EncryptionKey() to retrieve Encryption Key from Initial Key
func MessageInitialKeyCreate(
	transferId common.Address,
	chainId types.ChainId,
	signTypedData types.SignTypedDataCallback,
) (initialKey MessageInitialKey, err error) {
	signature, err := signTypedData(MessageInitialKeyTypedData(transferId, chainId))
	if err != nil {
		return
	}
	return MessageInitialKeyFromSignature(signature)
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
) (initialKey MessageInitialKey, err error) {
	return sha256.Sum256(MessageEncryptionKeyTypedDataSignature), nil
}

func (meki *MessageInitialKey) MessageEncryptionKey(
	encryptionKeyLength int64,
) (encryptionKey MessageEncryptionKey, err error) {
	encryptionKeyModified := base58.Encode(meki[:])
	if int64(len(encryptionKeyModified)) > encryptionKeyLength {
		encryptionKeyModified = encryptionKeyModified[:encryptionKeyLength]
	}
	encryptionKeyModifiedBytes, err := base58.Decode(encryptionKeyModified)
	if err != nil {
		return
	}
	encryptionKey = sha256.Sum256(encryptionKeyModifiedBytes)
	return
}
