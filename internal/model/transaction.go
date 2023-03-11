package model

type Transaction struct {
	ReceiverName string
	PayerName    string
	Amount       Amount
	Time         Time
}
