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
	"github.com/syronz/infrastructure/server/utils/hasher"
	"github.com/syronz/infrastructure/server/utils/activity"
	"github.com/syronz/infrastructure/server/utils/escape"
	"github.com/syronz/infrastructure/server/app"

	"github.com/syronz/ozzo-routing"

	"time"
)


// User type
type UserResource struct {}

// Return a user by id
func (r *UserResource) Get(c *routing.Context) error {
	id, err := validator.CheckId(c.Param("id"))
	if err != nil {
		errCustom := errors.CustomError(c, 403, dict.T(c,"ID must be a number"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	db := mysql.DB
	stmt, err := db.Prepare("SELECT id, name, username, role, city, director, language FROM users WHERE id = ?")
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Error in geting user"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}
	defer stmt.Close()

	var user models.User
	err = stmt.QueryRow(id).Scan(&user.ID, &user.Name, &user.Username, &user.Role, &user.City, &user.Director, &user.Language)
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Problem in scanning user"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}


	time.Sleep(1 * time.Millisecond)

	return c.Write(&models.Result{
		Status:true,
		Message: dict.T(c,"User's data"),
		Data:user,
	})
}


// Return users by pagination and search
func (r *UserResource) List(c *routing.Context) error {

	perPage := validator.PaginationPageSize(c.Query("perPage"))
	page := validator.PaginationPageNumber(c.Query("page"))

	search := escape.MakeSafe(c.Query("search"))
	var where string
	if search != "" {
		where = " WHERE username like '" + search +
		"%' OR name like '" + search + "%' " +
		" OR role like '" + search + "%' " +
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
	rows, err := db.Query("SELECT id, name, username, role, city, director, language FROM users " +
	where + sort + " LIMIT " + fmt.Sprint(skip) + "," + fmt.Sprint(limit))
	defer rows.Close()
	//rows, err := db.Query("SELECT, id, name, username, role, city, director, language FROM users ? ? LIMIT ?,?", where, sort, skip, limit)
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Unable to retrieve users data"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	var users []models.User
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.Name, &user.Username, &user.Role, &user.City, &user.Director, &user.Language)
		if err != nil {
			errCustom := errors.CustomError(c, 500, dict.T(c,"Problem in scan to the struct in users"), nil, err)
			debug.Log(err.Error())
			return c.Write(errCustom)
		}
		users = append(users, user)
	}

	var count int
	row, err := db.Query("SELECT COUNT(*) AS count from users " + where)
	defer row.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Count can't be calculated for users"), nil, err)
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

	//DSH TODO: Delete this delay
	//time.Sleep( 1 * time.Second )

	return c.Write(&models.Result{
		Status:true,
		Count:count,
		Message: dict.T(c,"List of users"),
		Data:users,
	})

}

// Function that create user inside the database, match json with User model
// after validation add it to the database, at the end print proper message or 
// data
func (r *UserResource) Create(c *routing.Context) error {
	data := &models.User{}

	if err := c.Read(&data); err != nil {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Password and user are required"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	var fields []string
	if data.Name == "" {
		fields = append(fields, "name")
	}
	if data.Username == "" {
		fields = append(fields, "username")
	}
	if data.Password == "" {
		fields = append(fields, "password")
	}
	if data.Role == "" {
		fields = append(fields, "role")
	}
	if data.City == "" {
		fields = append(fields, "city")
	}
	if data.Director == "" {
		fields = append(fields, "director")
	}
	if data.Language == "" {
		fields = append(fields, "language")
	}

	if len(fields) > 0 {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Red fields are required"), fields, nil)
		debug.Log("Required fields for create user", fields)
		return c.Write(errCustom)
	}


	// validation for username
	if len(data.Username) < 5 {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Username must be at least five character"), []string{"username"}, nil)
		debug.Log("length for username is less than six character")
		return c.Write(errCustom)
	}

	// validation for password
	if len(data.Password) < 8 {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Password must be at leas eight character"), []string{"password"}, nil)
		debug.Log("length for password is less than six character")
		return c.Write(errCustom)
	}

	encryptedPassword, err := hasher.HashPassword(data.Password)
	data.Password = ""
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Error in hashing password"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	// Insert to database with stmt method
	db := mysql.DB
	stmt, err := db.Prepare(`INSERT INTO users(name, username, password, role, city, director, language)
	VALUES(?,?,?,?,?,?,?)`)
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Error in inserting user to the database"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	res, err := stmt.Exec(data.Name, data.Username, encryptedPassword, data.Role, data.City, data.Director, data.Language )
	defer stmt.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 409, dict.T(c,"Duplication happened in user insertion"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	lid, err := res.LastInsertId()
	if err != nil {
		errCustom := errors.CustomError(c, 400, dict.T(c,"Can't return last id for user insertion"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	data.ID = int(lid)

	activity.Log(c, "user/create", fmt.Sprintf("ID=%d ", lid) + fmt.Sprint(data))

	return c.Write(&models.Result{
		Status:true,
		Message: dict.T(c,"User created successfully"),
		Data:data,
	})

}

// Delete user by id
func (r *UserResource) Delete(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errCustom := errors.CustomError(c, 403, dict.T(c,"ID must be number"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	userID := app.GetRequestScope(c).UserID()
	if userID == c.Param("id") {
		errCustom := errors.CustomError(c, 403, dict.T(c,"User can not delete own account"), nil, nil)
		return c.Write(errCustom)
	}

	db := mysql.DB
	stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Error in deleting user"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	_, err = stmt.Exec(id)
	defer stmt.Close()
	if err != nil {
		errCustom := errors.CustomError(c, 409, dict.T(c,"Problem in exec stmt for deleting user"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	activity.Log(c, "user/delete", fmt.Sprintf("ID=%d ", id))
	return c.Write(&models.Result{
		Status: true,
		Message: dict.T(c,"User deleted successfully"),
	})
}

// Edit user by id
func (r *UserResource) Update(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errCustom := errors.CustomError(c, 403, dict.T(c,"ID must be number"), nil, err)
		debug.Log(err)
		return c.Write(errCustom)
	}

	data := &models.User{}
	if err = c.Read(&data); err != nil {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Data is required for update the user"), []string{"name", "username", "role", "city"}, err)
		debug.Log(err)
		return c.Write(errCustom)
	}

	var fields []string
	if data.Name == "" {
		fields = append(fields, "name")
	}
	if data.Username == "" {
		fields = append(fields, "username")
	}
	if data.Role == "" {
		fields = append(fields, "role")
	}
	if data.City == "" {
		fields = append(fields, "city")
	}
	if data.Director == "" {
		fields = append(fields, "director")
	}
	if data.Language == "" {
		fields = append(fields, "language")
	}
	if len(fields) > 0 {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Red fields are required"), fields, nil)
		debug.Log("Required fields for update user", fields)
		return c.Write(errCustom)
	}

	var encryptedPassword string
	if data.Password != "" {
		encryptedPassword, err = hasher.HashPassword(data.Password)
		if err != nil {
			errCustom := errors.CustomError(c, 500, dict.T(c,"Error in hashing password"), nil, err)
			debug.Log(err)
			return c.Write(errCustom)
		}
	}

	// validation for username
	if len(data.Username) < 5 {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Username must be at least five character"), []string{"username"}, nil)
		debug.Log("username is less than six character")
		return c.Write(errCustom)
	}

	// validation for name
	if len(data.Name) < 4 {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Name must be at least four character"), []string{"name"}, nil)
		debug.Log("name is less than six character")
		return c.Write(errCustom)
	}

	// Update user with stmt method
	db := mysql.DB
	var query string
	if encryptedPassword != "" {
		query = "UPDATE users SET password = ?, username = ?, name = ?, role = ?, city = ?, director = ?, Language = ? WHERE id = ?"
	} else {
		query = "UPDATE users SET username = ?, name = ?, role = ?, city = ?, director = ?, Language = ? WHERE id = ?"
	}
	stmt, err := db.Prepare(query)
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Error in updating the user"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}
	defer stmt.Close()

	if encryptedPassword != "" {
		_, err = stmt.Exec(encryptedPassword, data.Username, data.Name, data.Role, data.City, data.Director, data.Language, id)
	} else {
		_, err = stmt.Exec(data.Username, data.Name, data.Role, data.City, data.Director, data.Language, id)
	}
	data.Password = ""
	if err != nil {
		errCustom := errors.CustomError(c, 409, dict.T(c,"This username already exist"), []string{"username"}, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	activity.Log(c, "user/update", fmt.Sprintf("ID=%d ", id) +
	fmt.Sprint(data))

	return c.Write(&models.Result{
		Status: true,
		Message: dict.T(c,"User updated successfully"),
		Data: data,
	})
}

