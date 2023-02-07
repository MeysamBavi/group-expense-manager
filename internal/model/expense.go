package model

import "time"

type Share struct {
	MemberName  string
	ShareWeight int
}

type Expense struct {
	Title         string
	Time          time.Time
	PayerName     string
	Amount        Amount
	Shares        []Share
	_totalWeights int
}

func (e *Expense) SumOfWeights() int {
	if e._totalWeights > 0 {
		return e._totalWeights
	}

	e._totalWeights = 0
	for _, share := range e.Shares {
		e._totalWeights += share.ShareWeight
	}

	return e._totalWeights
}
