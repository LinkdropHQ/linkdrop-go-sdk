package types

type EncryptedMessage struct {
	Message                string
	EncryptionKey          [32]byte
	EncryptionKeyLinkParam [32]byte
	EncryptionKeyLength    int64
}
