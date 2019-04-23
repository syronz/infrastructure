package models

import (
	"strconv"
)

type City struct {
	ID			int64	`json:"id"`
	Governorate	string	`json:"governorate"`
	City		string	`json:"city"`
}

func CityReturnMapOne(city City) map[string]string {
	tmp := make(map[string]string)
	tmp["id"] = strconv.FormatInt(city.ID, 10)
	tmp["governorate"] = city.Governorate
	tmp["city"] = city.City

	return tmp
}



