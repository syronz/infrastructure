package controllers

import (
	"fmt"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/syronz/infrastructure/server/errors"
	"github.com/syronz/infrastructure/server/models"
	"github.com/syronz/infrastructure/server/dict"
	"github.com/syronz/infrastructure/server/database/mysql"
	"github.com/syronz/infrastructure/server/utils/debug"
	"github.com/syronz/infrastructure/server/utils/hasher"
	"github.com/syronz/infrastructure/server/app"
	"github.com/syronz/infrastructure/server/utils/activity"
	"github.com/syronz/infrastructure/server/utils"

	"github.com/syronz/ozzo-routing"
	"github.com/syronz/ozzo-routing/auth"
	"os"
	//"io"
	"net/http"
)

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Auth type
type AuthResource struct {}

// Auth middleware for login
func Auth(signinKey string) routing.Handler {
	return func(c *routing.Context) error {

		var credential Credential
		if err := c.Read(&credential); err != nil {
			errCustom := errors.CustomError(c, 403, dict.T(c,"Username or password is incorrect"), nil, err)
			debug.Log(err.Error())
			return c.Write(errCustom)
		}

		identity := authenticate(credential)
		if identity == nil {
			errCustom := errors.CustomError(c, 403, dict.T(c,"Username or password is incorrect"), nil, nil)
			return c.Write(errCustom)
		}

		acls := app.Acl[identity.GetIdentityRole()]


		/* DSH: base method */
		tokenString, err := auth.NewJWT(jwt.MapClaims{
			"id": identity.GetID(),
			"name": identity.GetName(),
			"lang": identity.GetLanguage(),
			"exp": time.Now().Add(time.Hour * 72).Unix(),
		}, signinKey)
		if err != nil {
			errCustom := errors.CustomError(c, 500, dict.T(c,"Error in generate JWT"), nil, err)
			debug.Log(err.Error())
			return c.Write(errCustom)
		}


		//c.Response.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		//c.Response.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization")

		return c.Write(map[string]interface{}{
			"token": tokenString,
			"language": identity.GetLanguage(),
			"acls": acls,
			"status": true,
		})
	}
}


func JWTHandler(c *routing.Context, j *jwt.Token) error {
	//userID := j.Claims.(jwt.MapClaims)["id"].(float64)
	//app.GetRequestScope(c).SetUserID(fmt.Sprint(userID))
	//userLanguage := j.Claims.(jwt.MapClaims)["ln"].(string)
	//app.GetRequestScope(c).SetUserLanguage(userLanguage)

	claim, ok := j.Claims.(jwt.MapClaims)
	if ok != true {
		debug.Log("ERROR IN CLAIMING JWT", ok)
	}
	userID := claim["id"].(float64)
	app.GetRequestScope(c).SetUserID(fmt.Sprint(userID))

	userLanguage := claim["lang"]
	app.GetRequestScope(c).SetUserLanguage(fmt.Sprint(userLanguage))

	return nil
}

// Authenticate user if credential being exist
func authenticate(c Credential) models.Identity {
	db := mysql.DB
	stmt, err := db.Prepare("SELECT id, name, username, password, role, city, director, language FROM users WHERE username = ?")
	defer stmt.Close()
	if err != nil {
		debug.Log(err.Error())
		return nil
	}

	var user models.User
	err = stmt.QueryRow(c.Username).Scan(&user.ID,
	&user.Name, &user.Username, &user.Password, &user.Role, &user.City,
	&user.Director, &user.Language)

	if hasher.CheckPasswordHash(c.Password, user.Password) {
		return user
	}

	return nil
}


// Get Profile
func (r *AuthResource) GetProfile(c *routing.Context) error {
	userID := app.GetRequestScope(c).UserID()
	db := mysql.DB
	stmt, err := db.Prepare("SELECT id, name, username, role, city, director, language FROM users WHERE id = ?")
	defer stmt.Close()
	if err != nil {
		debug.Log(err.Error())
		return nil
	}

	var user models.User
	err = stmt.QueryRow(userID).Scan(&user.ID, &user.Name, &user.Username, &user.Role, &user.City,
		&user.Director, &user.Language)
	if err != nil {
		debug.Log(err.Error())
		return nil
	}

	return c.Write(&models.Result{
		Status: true,
		Message: dict.T(c, "User's profile"),
		Data: user,
	})

}

// Update Profile
func (r *AuthResource) UpdateProfile(c *routing.Context) error {
	id := app.GetRequestScope(c).UserID()
	var err error

	data := &models.User{}
	if err = c.Read(&data); err != nil {
		errCustom := errors.CustomError(c, 403, dict.T(c,"Data is required for update the user"), []string{"name", "username", "role", "city"}, err)
		debug.Log(err)
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
		query = "UPDATE users SET password = ?, name = ?, Language = ? WHERE id = ?"
	} else {
		query = "UPDATE users SET name = ?, Language = ? WHERE id = ?"
	}
	stmt, err := db.Prepare(query)
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c,"Error in updating the user"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}
	defer stmt.Close()

	if encryptedPassword != "" {
		_, err = stmt.Exec(encryptedPassword, data.Name, data.Language, id)
	} else {
		_, err = stmt.Exec(data.Name, data.Language, id)
	}
	data.Password = ""
	if err != nil {
		errCustom := errors.CustomError(c, 409, dict.T(c,"Problem in updating profile"), []string{"username"}, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	activity.Log(c, "profile/update", fmt.Sprintf("ID=%v ", id) + fmt.Sprint(data))

	return c.Write(&models.Result{
		Status: true,
		Message: dict.T(c,"Profile updated successfully"),
		Data: data,
	})
}


// Upload avatar
func (r *AuthResource) UplodadAvatar(c *routing.Context) error {
	userID := app.GetRequestScope(c).UserID()
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c, "Error in getting file"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}
	defer file.Close()

	if header.Size > 5 * 1000 * 1000 {
		errCustom := errors.CustomError(c, 500, dict.T(c, "File must be less than 5MB"), []string{"avatar"}, nil)
		debug.Log("File is too big")
		return c.Write(errCustom)
	}

	buff := make([]byte, header.Size)
	_, err = file.Read(buff)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	filetype := http.DetectContentType(buff)
	isImage := utils.Contains([]string{"image/png", "image/jpeg", "image/jpg", "image/gif"},
		filetype)
	if !isImage {
		errCustom := errors.CustomError(c, 403, dict.T(c, "File is not image"), []string{"avatar"}, nil)
		debug.Log("File is not image")
		return c.Write(errCustom)
	}


	var fileName string
	switch filetype {
	case "image/jpeg", "image/jpg":
		fileName = fmt.Sprintf("%v.jpg", userID)
	case "image/gif":
		fileName = fmt.Sprintf("%v.gif", userID)
	default:
		fileName = fmt.Sprintf("%v.png", userID)
	}

	f, err := os.OpenFile("./uploads/images/avatars/" + fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c, "Can not create the file"), nil, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}
	defer f.Close()

	if _, err := f.Write(buff); err != nil {
		errCustom := errors.CustomError(c, 500, dict.T(c, "Can not save the file"), []string{"avatar"}, err)
		debug.Log(err.Error())
		return c.Write(errCustom)
	}

	return c.Write(&models.Result{
		Status: true,
		Message: dict.T(c, "Avatar uploaded successfully"),
		Data: header,
	})



}
