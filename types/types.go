package types

// ClaimLink TODO
type ClaimLink interface{}

type ResultSet struct {
	Total  int64 `json:"total"`
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

type SenderHistory struct {
	ClaimLinks []ClaimLink `json:"claimLinks"`
	ResultSet  ResultSet   `json:"resultSet"`
}
