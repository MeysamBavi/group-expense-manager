package model

import (
	"math/big"
	"strconv"
	"strings"
)

type Amount struct {
	r *big.Rat
}

func AmountOf(a int64) Amount {
	return Amount{big.NewRat(a, 1)}
}

func AmountZero() Amount {
	return Amount{zeroRat()}
}

func (a Amount) IsZero() bool {
	return a.r == nil || a.r.Sign() == 0
}

func (a Amount) IsNegative() bool {
	return a.r != nil && a.r.Sign() == -1
}

func (a Amount) IsPositive() bool {
	return a.r != nil && a.r.Sign() == 1
}

func zeroRat() *big.Rat {
	return new(big.Rat)
}

func (a Amount) Negative() Amount {
	if a.IsZero() {
		return Amount{zeroRat()}
	}
	b := new(big.Rat).Set(a.r)
	return Amount{b.Neg(b)}
}

func (a Amount) Add(b Amount) Amount {
	switch {
	case a.IsZero() && b.IsZero():
		return Amount{zeroRat()}
	case a.IsZero():
		return Amount{zeroRat().Set(b.r)}
	case b.IsZero():
		return Amount{zeroRat().Set(a.r)}
	default:
		return Amount{zeroRat().Add(a.r, b.r)}
	}
}

func (a Amount) Sub(b Amount) Amount {
	return a.Add(b.Negative())
}

func (a Amount) Multiply(b int) Amount {
	if a.IsZero() {
		return Amount{zeroRat()}
	}
	r := zeroRat().SetInt64(int64(b))
	return Amount{r.Mul(r, a.r)}
}

func (a Amount) Divide(c int) Amount {
	if a.IsZero() {
		return Amount{zeroRat()}
	}
	r := zeroRat().SetInt64(int64(c))
	return Amount{r.Quo(a.r, r)}
}

func (a Amount) LessThan(b Amount) bool {
	switch {
	case a.IsZero() && b.IsZero():
		return false
	case a.IsZero():
		return b.r.Sign() == 1
	case b.IsZero():
		return a.r.Sign() == -1
	default:
		return a.r.Cmp(b.r) == -1
	}
}

func (a Amount) ToNumeral() int64 {
	if a.IsZero() {
		return 0
	}
	return (big.NewInt(0).Set(a.r.Num())).Div(a.r.Num(), a.r.Denom()).Int64()
}

func (a Amount) String() string {
	if a.IsZero() {
		return "0"
	}

	return a.r.FloatString(0)
}

func ParseAmount(a string) (Amount, error) {
	a = strings.TrimSpace(a)
	a = strings.ReplaceAll(a, ",", "")
	if a == "" {
		return Amount{zeroRat()}, nil
	}
	amount, err := strconv.ParseInt(a, 10, 64)
	return AmountOf(amount), err
}
