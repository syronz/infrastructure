package models

type Result struct {
	Status	bool	`json:"status"`
	Count	int	`json:"count"`
	Message	string	`json:"message"`
	Data	interface{} `json:"data"`
	Error	interface{} `json:"error"`
}
