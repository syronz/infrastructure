package controllers

import (
	"fmt"
	"strings"

	"github.com/syronz/infrastructure/server/database/mysql"
	"github.com/syronz/infrastructure/server/models"
	"github.com/syronz/infrastructure/server/utils/debug"
	"github.com/syronz/infrastructure/server/dict"
	"github.com/syronz/infrastructure/server/utils/validator"
	"github.com/syronz/infrastructure/server/utils/escape"
	"github.com/syronz/infrastructure/server/errors"

	"github.com/syronz/ozzo-routing"
)

// Activity type
type ActivityResource struct {}

// List of activities
func (r *ActivityResource) List(c *routing.Context) error {
	perPage := validator.PaginationPageSize(c.Query("perPage"))
	page := validator.PaginationPageNumber(c.Query("page"))

	search := escape.MakeSafe(c.Query("search"))
	var preSearch []string
	var preWhere string

	if strings.Contains(search, ">") {
		preSearch = strings.Split(search, ">")
		preWhere = " (event LIKE  '%" + preSearch[0] +
			"%' AND description LIKE '%" + preSearch[1] + "%') OR "
	}

	debug.Log(preWhere)


	var where string
	if search != "" {
		where = " WHERE "+ preWhere +" u.username like '" + search +
		"%' OR u.name like '" + search + "%' " +
		" OR a.event like '%" + search + "%' " +
		" OR a.created_at like '%" + search + "%' " +
		" OR a.ip like '" + search + "' " +
		" OR a.description like '%" + search + "%' " +
		" OR a.id like '" + search + "'"
	}

	sortField := escape.MakeSafe(c.Query("sortField"))
	sortDirection := escape.MakeSafe(c.Query("sortDirection"))
	var sort string
	if sortDirection != "" {
		sort = fmt.Sprintf(" ORDER BY %s %s ", sortField, sortDirection)
	}

	deletedWord := dict.T(c, "DELETED USER")
	skip := (page - 1) * perPage
	limit := perPage
	db := mysql.DB
	rows, err := db.Query("SELECT a.id, a.event, ifnull(u.name, '"+deletedWord+"') as name, a.created_at, a.ip," +
		"a.description FROM activity a left join users u on u.id = a.user_id " +
		where + sort + " LIMIT " + fmt.Sprint(skip) + "," + fmt.Sprint(limit))
	defer rows.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Unable to retrieve activities data"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	var activities []models.Activity
	for rows.Next() {
		var act models.Activity
		err = rows.Scan(&act.ID, &act.Event, &act.User, &act.CreatedAt, &act.IP, &act.Description)
		if err != nil {
			errCustom := errors.CustomError(c, 500, dict.T(c,"Problem in scan to the struct in activities"), nil, err)
			debug.Log(err.Error())
			return c.Write(errCustom)
		}
		activities = append(activities, act)
	}

	var count int
	row, err := db.Query("SELECT COUNT(*) AS count from activity a left join users u on u.id = a.user_id" + where)
	defer row.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Count can't be calculated for activities"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}
	row.Next()
	err = row.Scan(&count)
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Scan count faced problem"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	return c.Write(&models.Result{
		Status:true,
		Count:count,
		Message: dict.T(c,"List of activities"),
		Data:activities,
	})
}
