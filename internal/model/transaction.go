package model

import "time"

type Transaction struct {
	PayerID    MID
	ReceiverID MID
	Amount     Amount
	Time       time.Time
}
