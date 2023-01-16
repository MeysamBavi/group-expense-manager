package model

import "time"

type Transaction struct {
	PayerID    pid
	ReceiverID pid
	Amount     Amount
	Time       time.Time
}
