package types

type EncryptedMessage struct {
	Data          []byte
	EncryptionKey [32]byte
	InitialKey    [32]byte
}
