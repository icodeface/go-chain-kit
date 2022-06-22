package solana

import (
	"github.com/shopspring/decimal"
	"math/big"
)

// 将大数的Fil转换为小数
func ToSOL(v *big.Int) decimal.Decimal {
	d := decimal.NewFromBigInt(v, 0)
	return d.DivRound(decimal.NewFromInt(10).Pow(decimal.NewFromInt(9)), 9)
}

// 将小数的Fil转换为大数
func FromSOL(v decimal.Decimal) *big.Int {
	r := v.Mul(decimal.NewFromInt(10).Pow(decimal.NewFromInt(9)))
	return r.BigInt()
}
