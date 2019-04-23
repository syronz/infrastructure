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
	//"github.com/syronz/infrastructure/server/app"

)


// City type
type CityResource struct {}

// Return a city by id
func (r *CityResource) Get(c *routing.Context) error {
	id, err := validator.CheckId(c.Param("id"))
	if err != nil {
		errCustom := errors.CustomError(c, 403, dict.T(c,"ID must be a number"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	db := mysql.DB
	stmt, err := db.Prepare("SELECT * FROM cities WHERE id = ?")
	defer stmt.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Error in geting city"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	var city models.City
	err = stmt.QueryRow(id).Scan(&city.ID, &city.Governorate, &city.City)
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Problem in scanning city"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	return c.Write(&models.Result{
		Status:true,
		Message: dict.T(c,"Data of the city"),
		Data:city,
	})
}


// Return cities by pagination and search
func (r *CityResource) List(c *routing.Context) error {
	perPage := validator.PaginationPageSize(c.Query("perPage"))
	page := validator.PaginationPageNumber(c.Query("page"))

	search := c.Query("search")
	var where string
	if search != "" {
		where = " WHERE governorate like '" + search +
		"%' OR city like '" + search + "%' " +
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
	rows, err := db.Query("SELECT id, governorate, city FROM cities " + where +
	sort + " LIMIT " +	fmt.Sprint(skip) + "," + fmt.Sprint(limit))
	defer rows.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Unable to retrieve cities data"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	var cities []models.City
	for rows.Next() {
		var city models.City
		err = rows.Scan(&city.ID, &city.Governorate, &city.City)
		if err != nil {
			errCustom := errors.CustomError(c, 500, dict.T(c,"Problem in scan to the struct in cities"), nil, err)
			debug.Log(err.Error())
			return c.Write(errCustom)
		}
		cities = append(cities, city)
	}

	var count int
	row, err := db.Query("SELECT COUNT(*) AS count from cities " + where)
	defer row.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Count can't be calculated for cities"), nil, err)
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
		Message: dict.T(c,"List of cities"),
		Data:cities,
	})

}

// Function that create city inside the database, match json with City model
// after validation add it to the database, at the end print proper message or 
// data
func (r *CityResource) Create(c *routing.Context) error {
	data := &models.City{}

	if err := c.Read(&data); err != nil {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Data is required for create the city"), []string{"city","governorate"}, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	var fields []string
	if data.Governorate == "" {
		fields = append(fields, "governorate")
	}
	if data.City == "" {
		fields = append(fields, "city")
	}
	if len(fields) > 0 {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Red fields are required"), fields, nil)
		debug.Log("Fields Are required for create city", fields)
		return c.Write(errCustom)
	}

	// validation for city
	if len(data.City) < 2 {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Name of City must be at least two character"), []string{"city"}, nil)
		debug.Log("Name of City must be at least two character")
		return c.Write(errCustom)
	}

	// validation for governorate
	if len(data.Governorate) < 2 {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Name of Governorate must be at least two character"), []string{"governorate"}, nil)
		debug.Log("Name of Governorate must be at least two character")
		return c.Write(errCustom)
	}


	// Insert to database with stmt method
	db := mysql.DB
	stmt, err := db.Prepare("INSERT INTO cities(governorate, city) VALUES(?,?)")
	defer stmt.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Error in inserting city to the database"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	res, err := stmt.Exec(data.Governorate, data.City)
	if err != nil {
		errCustom := errors.CustomError(c, 409, dict.T(c,"Duplication happened in city insertion"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	lid, err := res.LastInsertId()
	if err != nil {
		errCustom := errors.CustomError(c, 400, dict.T(c,"Can't return last id for city insertion"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	data.ID = lid

	activity.Log(c, "city/create", fmt.Sprintf("ID=%d ", lid) + fmt.Sprint(data))

	return c.Write(&models.Result{
		Status:true,
		Message: dict.T(c,"City created successfully"),
		Data:data,
	})

}

// Delete city by id
func (r *CityResource) Delete(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errCustom := errors.CustomError(c, 403, dict.T(c,"ID must be number"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	db := mysql.DB
	stmt, err := db.Prepare("DELETE FROM cities WHERE id = ?")
	defer stmt.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Error in deleting city"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	_, err = stmt.Exec(id)
	defer stmt.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 409, dict.T(c,"Problem in excec stmt for deleting city"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	activity.Log(c, "city/delete", fmt.Sprintf("ID=%d ", id))
	return c.Write(&models.Result{
		Status: true,
		Message: dict.T(c,"City deleted successfully"),
	})

}

// Edit city by id
func (r *CityResource) Update(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errCustom := errors.CustomError(c, 403, dict.T(c,"ID must be number"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	data := &models.City{}
	if err = c.Read(&data); err != nil {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Data is required for update the city"), []string{"city","governorate"}, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	var fields []string
	if data.Governorate == "" {
		fields = append(fields, "governorate")
	}
	if data.City == "" {
		fields = append(fields, "city")
	}
	if len(fields) > 0 {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Red fields are required"), fields, nil)
		debug.Log("Fields Are required for create city", fields)
		return c.Write(errCustom)
	}

	// validation for city
	if len(data.City) < 2 {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Name of City must be at least two character"), []string{"city"}, nil)
		debug.Log("Name of City must be at least two character")
		return c.Write(errCustom)
	}

	// validation for governorate
	if len(data.Governorate) < 2 {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Name of Governorate must be at least two character"), []string{"governorate"}, nil)
		debug.Log("Name of Governorate must be at least two character")
		return c.Write(errCustom)
	}

	// Update city with stmt method
	db := mysql.DB
	stmt, err := db.Prepare("UPDATE cities SET governorate = ?, city = ? WHERE id = ?")
	defer stmt.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Error in updating the city"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	_, err = stmt.Exec(data.Governorate, data.City, id)
	if err != nil {
		errCustom := errors.CustomError(c, 409, dict.T(c,"Duplication for city"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}


	activity.Log(c, "city/update", fmt.Sprintf("ID=%d ", id) + fmt.Sprint(data))
	return c.Write(&models.Result{
		Status: true,
		Message: dict.T(c,"City updated successfully"),
		Data: data,
	})
}
