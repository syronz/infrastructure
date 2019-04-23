package controllers

import (
	"fmt"
	"strconv"
	"github.com/syronz/infrastructure/server/database/mysql"
	"github.com/syronz/infrastructure/server/models"
	"github.com/syronz/infrastructure/server/utils/debug"
	"github.com/syronz/infrastructure/server/dict"
	"github.com/syronz/infrastructure/server/errors"
	"github.com/syronz/infrastructure/server/utils/validator"
	"github.com/syronz/ozzo-routing"
	"github.com/syronz/infrastructure/server/utils/escape"
	"github.com/syronz/infrastructure/server/utils/activity"

)


// Director type
type DirectorResource struct {}

// Return a director by id
func (r *DirectorResource) Get(c *routing.Context) error {
	id, err := validator.CheckId(c.Param("id"))
	if err != nil {
		errCustom := errors.CustomError(c, 403, dict.T(c,"ID must be a number"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	db := mysql.DB
	stmt, err := db.Prepare("SELECT * FROM directors WHERE id = ?")
	defer stmt.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Error in geting director"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	var director models.Director
	err = stmt.QueryRow(id).Scan(&director.ID, &director.Director, &director.Detail)
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Problem in scanning director"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	return c.Write(&models.Result{
		Status:true,
		Message: dict.T(c,"Data of the director"),
		Data:director,
	})
}


// Return directors by pagination and search
func (r *DirectorResource) List(c *routing.Context) error {
	perPage := validator.PaginationPageSize(c.Query("perPage"))
	page := validator.PaginationPageNumber(c.Query("page"))

	search := c.Query("search")
	var where string
	if search != "" {
		where = " WHERE detail like '" + search +
		"%' OR director like '" + search + "%' " +
		" OR id like '" + search + "'"
	}

	sortField := escape.MakeSafe(c.Query("sortField"))
	sortDirection := escape.MakeSafe(c.Query("sortDirection"))
	var sort string
	if sortDirection != "" {
		sort = fmt.Sprintf(" ORDER BY %s %s ", sortField, sortDirection)
	}

	skip := (page - 1) * perPage
	limit := perPage
	db := mysql.DB
	rows, err := db.Query("SELECT id, detail, director FROM directors " + where +
	sort + " LIMIT " +	fmt.Sprint(skip) + "," + fmt.Sprint(limit))
	defer rows.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Unable to retrieve directors data"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	var directors []models.Director
	for rows.Next() {
		var director models.Director
		err = rows.Scan(&director.ID, &director.Detail, &director.Director)
		if err != nil {
			errCustom := errors.CustomError(c, 500, dict.T(c,"Problem in scan to the struct in directors"), nil, err)
			debug.Log(err.Error())
			return c.Write(errCustom)
		}
		directors = append(directors, director)
	}

	var count int
	row, err := db.Query("SELECT COUNT(*) AS count from directors " + where)
	defer row.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Count can't be calculated for directors"), nil, err)
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
		Message: dict.T(c,"List of directors"),
		Data:directors,
	})

}

// Function that create director inside the database, match json with Director model
// after validation add it to the database, at the end print proper message or 
// data
func (r *DirectorResource) Create(c *routing.Context) error {
	data := &models.Director{}

	if err := c.Read(&data); err != nil {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Data is required for create the director"), []string{"director","detail"}, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	var fields []string
	if data.Director == "" {
		fields = append(fields, "director")
	}
	if len(fields) > 0 {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Red fields are required"), fields, nil)
		debug.Log("Fields Are required for create director", fields)
		return c.Write(errCustom)
	}

	// validation for director
	if len(data.Director) < 2 {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Name of Director must be at least two character"), []string{"director"}, nil)
		debug.Log("Name of Director must be at least two character")
		return c.Write(errCustom)
	}

	// Insert to database with stmt method
	db := mysql.DB
	stmt, err := db.Prepare("INSERT INTO directors(detail, director) VALUES(?,?)")
	defer stmt.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Error in inserting director to the database"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	res, err := stmt.Exec(data.Detail, data.Director)
	if err != nil {
		errCustom := errors.CustomError(c, 409, dict.T(c,"Duplication happened in director insertion"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	lid, err := res.LastInsertId()
	if err != nil {
		errCustom := errors.CustomError(c, 400, dict.T(c,"Can't return last id for director insertion"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	data.ID = lid

	activity.Log(c, "director/create", fmt.Sprintf("ID=%d ", lid) + fmt.Sprint(data))
	return c.Write(&models.Result{
		Status:true,
		Message: dict.T(c,"Director created successfully"),
		Data:data,
	})

}

// Delete director by id
func (r *DirectorResource) Delete(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errCustom := errors.CustomError(c, 403, dict.T(c,"ID must be number"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	db := mysql.DB
	stmt, err := db.Prepare("DELETE FROM directors WHERE id = ?")
	defer stmt.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Error in deleting director"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	_, err = stmt.Exec(id)
	defer stmt.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 409, dict.T(c,"Problem in excec stmt for deleting director"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	activity.Log(c, "director/delete", fmt.Sprintf("ID=%d ", id))
	return c.Write(&models.Result{
		Status: true,
		Message: dict.T(c,"Director deleted successfully"),
	})

}

// Edit director by id
func (r *DirectorResource) Update(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errCustom := errors.CustomError(c, 403, dict.T(c,"ID must be number"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	data := &models.Director{}
	if err = c.Read(&data); err != nil {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Data is required for update the director"), []string{"director","detail"}, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	var fields []string
	if data.Director == "" {
		fields = append(fields, "director")
	}
	if len(fields) > 0 {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Red fields are required"), fields, nil)
		debug.Log("Fields Are required for create director", fields)
		return c.Write(errCustom)
	}

	// validation for director
	if len(data.Director) < 2 {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Name of Director must be at least two character"), []string{"director"}, nil)
		debug.Log("Name of Director must be at least two character")
		return c.Write(errCustom)
	}

	// Update director with stmt method
	db := mysql.DB
	stmt, err := db.Prepare("UPDATE directors SET detail = ?, director = ? WHERE id = ?")
	defer stmt.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Error in updating the director"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	_, err = stmt.Exec(data.Detail, data.Director, id)
	if err != nil {
		errCustom := errors.CustomError(c, 409, dict.T(c,"Duplication for director"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}


	activity.Log(c, "director/update", fmt.Sprintf("ID=%d ", id) + fmt.Sprint(data))
	return c.Write(&models.Result{
		Status: true,
		Message: dict.T(c,"Director updated successfully"),
		Data: data,
	})
}
