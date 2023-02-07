package model

import "time"

type Transaction struct {
	ReceiverName string
	PayerName    string
	Amount       Amount
	Time         time.Time
}
