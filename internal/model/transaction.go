package model

import "time"

type Transaction struct {
	ReceiverID MID
	PayerID    MID
	Amount     Amount
	Time       time.Time
}
