package constants

type Selector string

const (
	SelectorUndefined                   Selector = ""
	SelectorApproveWithAuthorization    Selector = "0xe1560fd3"
	SelectorReceiveWithAuthorizationEOA Selector = "0xef55bec6"
	SelectorReceiveWithAuthorization    Selector = "0x88b7ab63"
)
