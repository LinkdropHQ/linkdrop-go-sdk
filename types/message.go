package types

import (
	"crypto/sha256"
	"github.com/mr-tron/base58"
)

type EncryptedMessage struct {
	Data          []byte
	EncryptionKey [32]byte
	InitialKey    [32]byte
}

type MessageInitialKey [32]byte
type MessageEncryptionKey [32]byte

func (mki *MessageInitialKey) MessageEncryptionKey(
	encryptionKeyLength int64,
) (encryptionKey MessageEncryptionKey, err error) {
	encryptionKeyModified := base58.Encode(mki[:])
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
