package types

type EncryptedMessage struct {
	Message                []byte
	EncryptionKey          [32]byte
	EncryptionKeyLinkParam [32]byte
	EncryptionKeyLength    int64
}
