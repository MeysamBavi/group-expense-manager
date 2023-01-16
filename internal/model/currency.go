package model

import "fmt"

// Amount an amount of money with a specific currency
type Amount interface {
	fmt.Stringer
	toUint() uint32
	toFloat64() float64
}

type Rial uint32

func (r Rial) String() string {
	return fmt.Sprintf("%d تومان", r/10)
}

func (r Rial) toUint() uint32 {
	return uint32(r)
}

func (r Rial) toFloat64() float64 {
	return float64(r)
}
