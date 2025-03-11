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

func DecryptMessage(
	message *types.EncryptedMessage,
) (string, error) {
	encoded := message.Data[2:] // Remove length from message
	return crypto.Decrypt(encoded, message.EncryptionKey)
}

func EncryptMessage(
	message string,
	initialKey types.MessageInitialKey,
	encryptionKeyLength int64,
	nonce [crypto.NonceLength]byte,
) (encryptedMessage *types.EncryptedMessage, err error) {
	if encryptionKeyLength > 0xFFFF {
		return nil, errors.New("encryptionKeyLength exceeds 2 bytes")
	}

	encryptionKey, err := initialKey.MessageEncryptionKey(encryptionKeyLength)

	encryptedSenderMessage, err := crypto.Encrypt([]byte(message), encryptionKey, nonce)
	if err != nil {
		return nil, err
	}

	encryptionKeyLengthB := make([]byte, 2)
	binary.BigEndian.PutUint16(encryptionKeyLengthB, uint16(encryptionKeyLength))

	// Build the result
	result := &types.EncryptedMessage{
		Data:          append(encryptionKeyLengthB[:], encryptedSenderMessage...),
		InitialKey:    initialKey,
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
) (initialKey types.MessageInitialKey, err error) {
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
) (initialKey types.MessageInitialKey, err error) {
	return sha256.Sum256(MessageEncryptionKeyTypedDataSignature), nil
}
