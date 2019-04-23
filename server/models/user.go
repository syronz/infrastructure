package models

import (
	"github.com/syronz/infrastructure/server/database/mysql"
)

type User struct {
	ID				int		`json:"id"`
	Name			string	`json:"name"`
	Username		string	`json:"username"`
	Password		string	`json:"password"`
	Role			string	`json:"role"`
	City			string	`json:"city"`
	Director		string	`json:"director"`
	Language		string	`json:"language"`
}

func (u User) GetID() int {
	return u.ID
}

func (u User) GetName() string {
	return u.Name
}

func (u User) GetLanguage() string {
	return u.Language
}

func (u User) GetIdentityRole() string {
	return u.Role
}

// Get user role for ACL check
func (u User) GetRole(id string) (User, error) {
	db := mysql.DB
	var user User
	stmt, err := db.Prepare("SELECT role FROM users WHERE id = ?")
	if err != nil {
		return user, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&user.Role)
	if err != nil {
		return user, err
	}

	return user, nil

}
