package model

import "time"

type Share struct {
	MemberName  string
	ShareWeight int
}

type Expense struct {
	Title     string
	Time      time.Time
	PayerName string
	Amount    Amount
	Shares    []Share
}
