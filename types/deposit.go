package types

type DepositRequest struct {
	Sender                 string `json:"sender"`
	Escrow                 string `json:"escrow"`
	TransferID             string `json:"transferID"`
	Token                  string `json:"token"`
	TokenType              string `json:"tokenType"`
	Expiration             string `json:"expiration"`
	TxHash                 string `json:"txHash"`
	FeeAuthorization       string `json:"feeAuthorization"`
	Amount                 string `json:"amount"`
	FeeAmount              string `json:"feeAmount"`
	TotalAmount            string `json:"totalAmount"`
	FeeToken               string `json:"feeToken"`
	EncryptedSenderMessage string `json:"encryptedSenderMessage"`
}
