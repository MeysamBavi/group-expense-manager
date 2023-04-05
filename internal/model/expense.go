package model

import (
	"strconv"
	"strings"
)

type Share struct {
	MemberName  string
	ShareWeight int
}

type Expense struct {
	Title     string
	Time      Time
	PayerName string
	Amount    Amount
	Shares    []Share
}

func (e *Expense) SumOfWeights() int {
	sum := 0
	for _, share := range e.Shares {
		sum += share.ShareWeight
	}
	return sum
}

func ParseShareWeight(weightStr string) (int, error) {
	weightStr = strings.TrimSpace(weightStr)

	if weightStr == "" {
		return 0, nil
	}

	if b, err := strconv.ParseBool(weightStr); err == nil {
		var weight int
		if b {
			weight = 1
		} else {
			weight = 0
		}
		return weight, nil
	}

	return strconv.Atoi(weightStr)
}
