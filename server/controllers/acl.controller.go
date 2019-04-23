package controllers

import (
	"net/http"
	"github.com/syronz/ozzo-routing"
	"github.com/syronz/infrastructure/server/app"
	"github.com/syronz/infrastructure/server/utils/debug"
	"github.com/syronz/infrastructure/server/models"
	//"github.com/syronz/infrastructure/server/errors"
	"github.com/syronz/infrastructure/server/dict"
)

func AclChecker(part string, action string) routing.Handler {
	return func(c *routing.Context) error {
		id := app.GetRequestScope(c).UserID()

		var user models.User
		user, err := user.GetRole(id)
		if err != nil {
			debug.Log(err.Error())
			return routing.NewHTTPError(http.StatusInternalServerError, dict.T(c,"Problem with getting user's information"))
		}

		roles, ok := app.Acl[user.Role]
		if ok != true {
			debug.Log("Role was not defined")
			return routing.NewHTTPError(http.StatusUnauthorized, dict.T(c,"Role was not defined"))
		}

		switch action {
		case "read":
			if contains(roles.Read, part) {
				return nil
			} else {
				return routing.NewHTTPError(http.StatusMethodNotAllowed, dict.T(c,"You don't have permission"))
			}
		case "write":
			if contains(roles.Write, part) {
				return nil
			} else {
				return routing.NewHTTPError(http.StatusMethodNotAllowed, dict.T(c,"You don't have permission"))
			}
		}


		return routing.NewHTTPError(http.StatusMethodNotAllowed, dict.T(c,"Action not allowed"))

	}
}

func contains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}
