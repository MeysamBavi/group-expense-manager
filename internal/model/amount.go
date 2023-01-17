package model

import "strconv"

type Amount uint64

func (a Amount) Add(b Amount) Amount {
	return a + b
}

func (a Amount) Sub(b Amount) Amount {
	return a - b
}

func (a Amount) Multiply(c uint64) Amount {
	return a * Amount(c)
}

func (a Amount) Divide(c uint64) Amount {
	return a / Amount(c)
}

func ParseAmount(a string) (Amount, error) {
	amount, err := strconv.ParseUint(a, 10, 64)
	return Amount(amount), err
}
