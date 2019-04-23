package models

type Activity struct {
	ID				int		`json:"id"`
	CreatedAt		string	`json:"created_at"`
	Event			string	`json:"event"`
	UserId			int		`json:"user_id"`
	IP				string	`json:"ip"`
	Description		string	`json:"description"`
	User			string	`json:"user"`
}
