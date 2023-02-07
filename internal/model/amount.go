package model

import "strconv"

type Amount int64

const AmountZero = Amount(0)

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

func ParseAmount(a string) (Amount, error) {
	amount, err := strconv.ParseInt(a, 10, 64)
	return Amount(amount), err
}
