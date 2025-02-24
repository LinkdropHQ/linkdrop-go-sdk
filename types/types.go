package types

type ResultSet struct {
	Total  int64 `json:"total"`
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}
