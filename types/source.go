package types

// LinkSource
// Reflects link format
type LinkSource string

const (
	LinkSourceUndefined LinkSource = ""
	LinkSourceP2P       LinkSource = "P2P"
	LinkSourceDashboard LinkSource = "D"
)
