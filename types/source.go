package types

// LinkSource
// Reflects link format
type LinkSource string

const (
	LinkSourceUndefined LinkSource = ""
	LinkSourceP2P       LinkSource = "p2p"
	LinkSourceDashboard LinkSource = "d"
)
