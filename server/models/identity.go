package models

type Identity interface {
	GetID() int
	GetName() string
	GetLanguage() string
	GetIdentityRole() string
}
