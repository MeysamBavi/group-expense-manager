package model

import "time"

type Share struct {
	MemberID    MID
	ShareWeight int
}

type Expense struct {
	Title   string
	Time    time.Time
	PayerID MID
	Amount  Amount
	Shares  []Share
}
