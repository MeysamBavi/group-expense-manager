package model

type pid uint32

type Person struct {
	ID         pid
	Name       string
	CardNumber string
}
