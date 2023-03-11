package model_test

import (
	"github.com/MeysamBavi/group-expense-manager/internal/model"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestAmount_Add(t *testing.T) {
	assert := assert2.New(t)
	assert.Equal(int64(34), model.AmountOf(17).Add(model.AmountOf(17)).ToNumeral())
	assert.Equal(int64(50), model.AmountOf(34).Add(model.AmountOf(16)).ToNumeral())
	assert.Equal(int64(34), model.AmountOf(-1).Add(model.AmountOf(35)).ToNumeral())
	assert.Equal(int64(0), model.AmountOf(0).Add(model.AmountOf(0)).ToNumeral())
	assert.Equal(int64(103833), model.AmountOf(0).Add(model.AmountOf(103833)).Add(model.AmountZero()).ToNumeral())
	assert.Equal(int64(207666), model.AmountZero().Add(model.AmountOf(103833)).Add(model.AmountOf(103833)).ToNumeral())
	assert.Equal(int64(1), model.AmountOf(270000).
		Add(model.AmountOf(-270000)).Add(model.AmountOf(1)).ToNumeral())
}

func TestAmount_Sub(t *testing.T) {
	assert := assert2.New(t)
	assert.Equal(int64(43), model.AmountOf(80).Sub(model.AmountOf(37)).ToNumeral())
	assert.Equal(int64(80), model.AmountOf(80).Sub(model.AmountOf(0)).ToNumeral())
	assert.Equal(int64(80), model.AmountOf(0).Sub(model.AmountOf(-80)).ToNumeral())
	assert.Equal(int64(120), model.AmountOf(60).Sub(model.AmountOf(-60)).ToNumeral())
}

func TestAmount_Multiply(t *testing.T) {
	assert := assert2.New(t)
	assert.Equal(int64(64), model.AmountOf(16).Multiply(4).ToNumeral())
	assert.Equal(int64(-64), model.AmountOf(-16).Multiply(4).ToNumeral())
	assert.Equal(int64(0), model.AmountOf(-16).Multiply(0).ToNumeral())
	assert.Equal(int64(429), model.AmountOf(143).Multiply(3).ToNumeral())
	assert.Equal(int64(-6400), model.AmountOf(800).Multiply(-8).ToNumeral())
	assert.Equal(int64(6400), model.AmountOf(-800).Multiply(-8).ToNumeral())
}

func TestAmount_Divide(t *testing.T) {
	assert := assert2.New(t)
	assert.Equal(int64(4), model.AmountOf(102).Divide(25).ToNumeral())
	assert.Equal(int64(4), model.AmountOf(100).Divide(25).ToNumeral())
	assert.Equal(int64(3), model.AmountOf(100).Divide(26).ToNumeral())
	assert.Equal(int64(0), model.AmountOf(0).Divide(26).ToNumeral())
	assert.Equal(int64(4), model.AmountOf(-8).Divide(-2).ToNumeral())
	assert.Equal(int64(2), model.AmountOf(-8).Divide(-3).ToNumeral())
}

func TestAmount_Negative(t *testing.T) {
	assert := assert2.New(t)
	assert.Equal(int64(-1), model.AmountOf(1).Negative().ToNumeral())
	assert.Equal(int64(0), model.AmountOf(0).Negative().ToNumeral())
	assert.Equal(int64(1), model.AmountOf(-1).Negative().ToNumeral())
	assert.Equal(int64(-(1 << 31)), model.AmountOf(1<<31).Negative().ToNumeral())
}

func TestAmount_IsNegative(t *testing.T) {
	assert := assert2.New(t)
	assert.Equal(true, model.AmountOf(-2).IsNegative())
	assert.Equal(true, model.AmountOf(-1).IsNegative())
	assert.Equal(false, model.AmountOf(0).IsNegative())
	assert.Equal(false, model.AmountOf(1).IsNegative())
	assert.Equal(false, model.AmountOf(1<<32).IsNegative())
	assert.Equal(true, model.AmountOf(-(1 << 32)).IsNegative())
}

func TestAmount_LessThan(t *testing.T) {
	assert := assert2.New(t)
	assert.Equal(true, model.AmountOf(64).LessThan(model.AmountOf(66)))
	assert.Equal(true, model.AmountOf(0).LessThan(model.AmountOf(66)))
	assert.Equal(true, model.AmountOf(-1).LessThan(model.AmountOf(0)))
	assert.Equal(true, model.AmountOf(-11).LessThan(model.AmountOf(-5)))
	assert.Equal(true, model.AmountOf(-11).LessThan(model.AmountOf(-10)))
	assert.Equal(false, model.AmountOf(-11).LessThan(model.AmountOf(-11)))
	assert.Equal(false, model.AmountOf(64).LessThan(model.AmountOf(-10)))
	assert.Equal(false, model.AmountOf(64).LessThan(model.AmountOf(63)))
	assert.Equal(false, model.AmountOf(64).LessThan(model.AmountOf(64)))
}

func TestAmount_Immutability(t *testing.T) {
	assert := assert2.New(t)

	a := model.AmountOf(3)
	b := model.AmountOf(13)

	a.Add(b)
	assert.Equal(int64(3), a.ToNumeral())
	assert.Equal(int64(13), b.ToNumeral())

	a.Sub(b)
	assert.Equal(int64(3), a.ToNumeral())
	assert.Equal(int64(13), b.ToNumeral())

	a.Multiply(3)
	assert.Equal(int64(3), a.ToNumeral())
	assert.Equal(int64(13), b.ToNumeral())

	a.Divide(18)
	assert.Equal(int64(3), a.ToNumeral())
	assert.Equal(int64(13), b.ToNumeral())

	a.Divide(2)
	assert.Equal(int64(3), a.ToNumeral())
	assert.Equal(int64(13), b.ToNumeral())

	c := model.AmountOf(19).Divide(2)
	assert.Equal(int64(9), c.ToNumeral())
	assert.Equal(int64(9), c.ToNumeral())

	a.Negative()
	assert.Equal(int64(3), a.ToNumeral())
	assert.Equal(int64(13), b.ToNumeral())
}
