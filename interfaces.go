package linkdrop

import "github.com/ethereum/go-ethereum/common"

type IClaimLinkRedeemable interface {
	Redeem(receiver common.Address) (txHash common.Hash, err error)
}
