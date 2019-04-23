package models

type Customer struct {
	ID			int		`json:"id"`
	Title		string	`json:"title"`
	Name		string	`json:"name"`
	Phone1		string	`json:"phone1"`
	Phone2		string	`json:"phone2"`
	CreatedAt	string	`json:"created_at"`
	Detail		string	`json:"detail"`
}


