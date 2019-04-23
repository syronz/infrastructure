package models

type Message struct {
	Message		string	`json:"message"`
	Params		[]string	`json:"params"`
	Extra		string	`json:"extra"`
}
