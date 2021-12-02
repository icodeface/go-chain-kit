package ethereum

import (
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/icodeface/hdwallet"
)

func ValidateAddress(addr string) bool {
	return common.IsHexAddress(addr)
}

func DerivePath(account int64, index int64) hdwallet.DerivationPath {
	return hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/%d'/0/%d", account, index))
}

func NewFilterQuery(blockHash *common.Hash, contract common.Address, event common.Hash) ethereum.FilterQuery {
	return ethereum.FilterQuery{
		BlockHash: blockHash,
		Addresses: []common.Address{
			contract,
		},
		Topics: [][]common.Hash{
			{event},
		},
	}
}
