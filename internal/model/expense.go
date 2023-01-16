package model

import "time"

type Share struct {
	PersonID    pid
	ShareWeight uint8
}

type Expense struct {
	Title   string
	Time    time.Time
	PayerID pid
	Amount  Amount
	Shares  []Share
}
