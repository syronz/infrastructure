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
	"github.com/syronz/infrastructure/server/utils/activity"
	"github.com/syronz/infrastructure/server/utils/escape"

	"github.com/syronz/ozzo-routing"

	"time"
)


// Customer type
type CustomerResource struct {}

// Return a customer by id
func (r *CustomerResource) Get(c *routing.Context) error {
	id, err := validator.CheckId(c.Param("id"))
	if err != nil {
		errCustom := errors.CustomError(c, 403, dict.T(c,"ID must be a number"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	db := mysql.DB
	stmt, err := db.Prepare("SELECT id, name, title, phone1, phone2, created_at, detail FROM customers WHERE id = ?")
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Error in geting customer"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}
	defer stmt.Close()

	var customer models.Customer
	err = stmt.QueryRow(id).Scan(&customer.ID, &customer.Name, &customer.Title, &customer.Phone1, &customer.Phone2, &customer.CreatedAt, &customer.Detail)
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Problem in scanning customer"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}


	time.Sleep(1 * time.Millisecond)

	return c.Write(&models.Result{
		Status:true,
		Message: dict.T(c,"Customer's data"),
		Data:customer,
	})
}


// Return customers by pagination and search
func (r *CustomerResource) List(c *routing.Context) error {

	perPage := validator.PaginationPageSize(c.Query("perPage"))
	page := validator.PaginationPageNumber(c.Query("page"))

	search := escape.MakeSafe(c.Query("search"))
	var where string
	if search != "" {
		where = " WHERE title like '" + search +
		"%' OR name like '%" + search + "%' " +
		" OR phone1 like '" + search + "%' " +
		" OR phone2 like '" + search + "%' " +
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
	rows, err := db.Query("SELECT id, name, title, phone1, phone2, created_at, detail FROM customers " +
	where + sort + " LIMIT " + fmt.Sprint(skip) + "," + fmt.Sprint(limit))
	defer rows.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Unable to retrieve customers data"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	var customers []models.Customer
	for rows.Next() {
		var customer models.Customer
		err = rows.Scan(&customer.ID, &customer.Name, &customer.Title, &customer.Phone1, &customer.Phone2, &customer.CreatedAt, &customer.Detail)
		if err != nil {
			errCustom := errors.CustomError(c, 500, dict.T(c,"Problem in scan to the struct in customers"), nil, err)
			debug.Log(err.Error())
			return c.Write(errCustom)
		}
		customers = append(customers, customer)
	}

	var count int
	row, err := db.Query("SELECT COUNT(*) AS count from customers " + where)
	defer row.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Count can't be calculated for customers"), nil, err)
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
		Message: dict.T(c,"List of customers"),
		Data:customers,
	})

}

// Function that create customer inside the database, match json with Customer model
// after validation add it to the database, at the end print proper message or 
// data
func (r *CustomerResource) Create(c *routing.Context) error {
	data := &models.Customer{}

	if err := c.Read(&data); err != nil {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Data for customer is required"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	var fields []string
	if data.Name == "" {
		fields = append(fields, "name")
	}
	if data.Phone1 == "" {
		fields = append(fields, "phone1")
	}

	if len(fields) > 0 {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Red fields are required"), fields, nil)
		debug.Log("Required fields for create customer", fields)
		return c.Write(errCustom)
	}


	// validation for phone1
	if !(len(data.Phone1) == 7 || len(data.Phone1) == 11) {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Phone is wrong"), []string{"phone1"}, nil)
		debug.Log("Phone1 must be 7 or 11 character")
		return c.Write(errCustom)
	}

	// validation for name
	if len(data.Name) < 8 {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Name must be at leas eight character"), []string{"name"}, nil)
		debug.Log("length for customer name is less than eight character")
		return c.Write(errCustom)
	}

	// Insert to database with stmt method
	db := mysql.DB
	stmt, err := db.Prepare(`INSERT INTO customers(name, title, phone1, phone2, detail)
	VALUES(?,?,?,?,?)`)
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Error in inserting customer to the database"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	res, err := stmt.Exec(data.Name, data.Title, data.Phone1, data.Phone2, data.Detail )
	defer stmt.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 409, dict.T(c,"Duplication happened, use another number as phone1"), []string{"phone1"}, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	lid, err := res.LastInsertId()
	if err != nil {
		errCustom := errors.CustomError(c, 400, dict.T(c,"Can't return last id for customer insertion"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	data.ID = int(lid)

	activity.Log(c, "customer/create", fmt.Sprintf("ID=%d ", lid) + fmt.Sprint(data))

	return c.Write(&models.Result{
		Status:true,
		Message: dict.T(c,"Customer created successfully"),
		Data:data,
	})

}

// Delete customer by id
func (r *CustomerResource) Delete(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errCustom := errors.CustomError(c, 403, dict.T(c,"ID must be number"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	db := mysql.DB
	stmt, err := db.Prepare("DELETE FROM customers WHERE id = ?")
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Error in deleting customer"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	_, err = stmt.Exec(id)
	defer stmt.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 409, dict.T(c,"Problem in exec stmt for deleting customer"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	activity.Log(c, "customer/delete", fmt.Sprintf("ID=%d ", id))
	return c.Write(&models.Result{
		Status: true,
		Message: dict.T(c,"Customer deleted successfully"),
	})
}

// Edit customer by id
func (r *CustomerResource) Update(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errCustom := errors.CustomError(c, 403, dict.T(c,"ID must be number"), nil, err)
		debug.Log(err)
		return c.Write(errCustom)
	}

	data := &models.Customer{}
	if err = c.Read(&data); err != nil {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Data is required for update the customer"), []string{"name", "phone1"}, err)
		debug.Log(err)
		return c.Write(errCustom)
	}

	var fields []string
	if data.Name == "" {
		fields = append(fields, "name")
	}
	if data.Phone1 == "" {
		fields = append(fields, "phone1")
	}

	if len(fields) > 0 {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Red fields are required"), fields, nil)
		debug.Log("Required fields for create customer", fields)
		return c.Write(errCustom)
	}


	// validation for phone1
	if !(len(data.Phone1) == 7 || len(data.Phone1) == 11) {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Phone is wrong"), []string{"phone1"}, nil)
		debug.Log("Phone1 must be 7 or 11 character")
		return c.Write(errCustom)
	}

	// validation for name
	if len(data.Name) < 8 {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Name must be at leas eight character"), []string{"name"}, nil)
		debug.Log("length for customer name is less than eight character")
		return c.Write(errCustom)
	}


	// Update customer with stmt method
	db := mysql.DB
	var query string
	query = "UPDATE customers SET title = ?, name = ?, phone1 = ?, phone2 = ?, Detail = ? WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Error in updating the customer"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}
	defer stmt.Close()

	_, err = stmt.Exec(data.Title, data.Name, data.Phone1, data.Phone2, data.Detail, id)
	if err != nil {
		errCustom := errors.CustomError(c, 409, dict.T(c,"Dupplication happened"), []string{"phone1"}, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	activity.Log(c, "customer/update", fmt.Sprintf("ID=%d ", id) +
	fmt.Sprint(data))

	return c.Write(&models.Result{
		Status: true,
		Message: dict.T(c,"Customer updated successfully"),
		Data: data,
	})
}


