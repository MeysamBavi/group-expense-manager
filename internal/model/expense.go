package model

import "time"

type Share struct {
	MemberID    MID
	ShareWeight uint8
}

type Expense struct {
	Title   string
	Time    time.Time
	PayerID MID
	Amount  Amount
	Shares  []Share
}
