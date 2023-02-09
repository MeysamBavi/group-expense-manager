package model

import (
	"strconv"
	"strings"
)

type Amount int64

const AmountZero = Amount(0)

func (a Amount) Negative() Amount {
	return -a
}

func (a Amount) Add(b Amount) Amount {
	return a + b
}

func (a Amount) Sub(b Amount) Amount {
	return a - b
}

func (a Amount) Multiply(c int) Amount {
	return a * Amount(c)
}

func (a Amount) Divide(c int) Amount {
	return a / Amount(c)
}

func (a Amount) LessThan(b Amount) bool {
	return a < b
}

func (a Amount) ToNumeral() int64 {
	return int64(a)
}

func ParseAmount(a string) (Amount, error) {
	if a == "" {
		return 0, nil
	}
	a = strings.ReplaceAll(a, ",", "")
	amount, err := strconv.ParseInt(a, 10, 64)
	return Amount(amount), err
}
