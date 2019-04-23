package models

type Error struct {
	Code		int	`json:"code"`
	Message		string	`json:"message"`
	Params		[]string	`json:"params"`
	Fields		[]string	`json:"fields"`
	Extra		string	`json:"extra"`
}



