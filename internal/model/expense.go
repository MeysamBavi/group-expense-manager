package model

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
