package filecoin

import (
	"crypto/sha256"
	"fmt"
	"github.com/cbergoon/merkletree"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	big2 "github.com/filecoin-project/go-state-types/big"
	"github.com/icodeface/hdkeyring"
	"github.com/ipfs/go-cid"
	"github.com/shopspring/decimal"
	"math/big"
)

// 将大数的Fil转换为小数
func ToFil(v abi.TokenAmount) decimal.Decimal {
	d := decimal.NewFromBigInt(v.Int, 0)
	return d.DivRound(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18)), 18)
}

// 将小数的Fil转换为大数
func FromFil(v decimal.Decimal) abi.TokenAmount {
	return big2.NewFromGo(v.Mul(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18))).BigInt())
}

func BigIntToTokenAmount(v *big.Int) abi.TokenAmount {
	return big2.NewFromGo(v)
}

func BigIntToFIL(v *big.Int) decimal.Decimal {
	d := decimal.NewFromBigInt(v, 0)
	return d.DivRound(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18)), 18)
}

func ValidateAddress(addr string) bool {
	if _, err := address.NewFromString(addr); err != nil {
		return false
	}
	return true
}

type blockId struct {
	cid.Cid
}

func (bid blockId) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write(bid.Bytes()); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

//Equals tests for equality of two Contents
func (bid blockId) Equals(other merkletree.Content) (bool, error) {
	return bid.Cid.Equals(other.(blockId).Cid), nil
}

func CalculateCidsMerkleRoot(cids []cid.Cid) string {
	var list []merkletree.Content
	for _, id := range cids {
		list = append(list, blockId{id})
	}
	t, err := merkletree.NewTree(list)
	if err != nil {
		panic(err)
	}
	return hexutil.Encode(t.MerkleRoot())
}

func DerivePath(account int64, index int64) hdkeyring.DerivationPath {
	return hdkeyring.MustParseDerivationPath(fmt.Sprintf("m/44'/461'/%d'/0/%d", account, index))
}
