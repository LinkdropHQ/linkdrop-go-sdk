package helpers

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/LinkdropHQ/linkdrop-go-sdk/internal/crypto"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/mr-tron/base58"
)

func DecryptMessage(
	message *types.EncryptedMessage,
) (string, error) {
	encoded := message.Message[2:] // Remove length from message
	return crypto.Decrypt(encoded, message.EncryptionKey)
}

func EncryptMessage(
	message string,
	transferID common.Address,
	chainID types.ChainId,
	encryptionKeyLength int64,
	getRandomBytes types.RandomBytesCallback,
	signTypedData types.SignTypedDataCallback,
) (*types.EncryptedMessage, error) {
	if encryptionKeyLength > 0xFFFF {
		return nil, errors.New("encryptionKeyLength exceeds 2 bytes")
	}

	encryptionKey, _, err := CreateMessageEncryptionKey(transferID.Hex(), signTypedData, chainID, encryptionKeyLength)
	if err != nil {
		return nil, err
	}

	encryptedSenderMessage, err := crypto.Encrypt([]byte(message), encryptionKey, [24]byte{}, getRandomBytes)
	if err != nil {
		return nil, err
	}

	encryptionKeyLengthB := make([]byte, 2)
	binary.BigEndian.PutUint16(encryptionKeyLengthB, uint16(encryptionKeyLength))

	// Build the result
	result := &types.EncryptedMessage{
		Message:       append(encryptionKeyLengthB[:], encryptedSenderMessage...),
		EncryptionKey: encryptionKey,
	}
	return result, nil
}

// CreateMessageEncryptionKey creates a message encryption key
func CreateMessageEncryptionKey(
	transferId string,
	signTypedData types.SignTypedDataCallback,
	chainId types.ChainId,
	encryptionKeyLength int64,
) (encryptionKey [32]byte, encryptionKeyLinkParam [32]byte, err error) {
	// Generating signature
	td := apitypes.TypedData{
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
			"seed": fmt.Sprintf("Encrypting message (transferId: %s)", transferId),
		},
	}
	signature, err := signTypedData(td)
	if err != nil {
		return [32]byte{}, [32]byte{}, err
	}

	// Calculating the initial encryption key
	encryptionKeyInitial := sha256.Sum256(signature)

	// Encoding to Base58 and trimming
	encryptionKeyModified := base58.Encode(encryptionKeyInitial[:])
	if int64(len(encryptionKeyModified)) > encryptionKeyLength {
		encryptionKeyModified = encryptionKeyModified[:encryptionKeyLength]
	}

	// Converting Base58 code to byte array before hashing
	encryptionKeyModifiedBytes, err := base58.Decode(encryptionKeyModified)
	if err != nil {
		return [32]byte{}, [32]byte{}, err
	}

	// Final encryption
	finalHash := sha256.Sum256(encryptionKeyModifiedBytes)

	return finalHash, encryptionKeyInitial, nil
}

func EncryptionKeyFromLink(
	linkKey string,
	encryptionKeyLength int64,
) ([]byte, error) {
	linkKeyB, err := hex.DecodeString(linkKey)
	if err != nil {
		return nil, err
	}
	if len(linkKeyB) > 32 {
		return nil, fmt.Errorf("decoded byte length is not 32; got %d", len(linkKeyB))
	}
	encryptionKeyModified := base58.Encode(linkKeyB)
	if int64(len(encryptionKeyModified)) > encryptionKeyLength {
		encryptionKeyModified = encryptionKeyModified[:encryptionKeyLength]
	}

	encryptionKeyModifiedBytes, err := base58.Decode(encryptionKeyModified)
	if err != nil {
		return []byte{}, err
	}

	key := sha256.Sum256(encryptionKeyModifiedBytes)
	return key[:], nil
}
