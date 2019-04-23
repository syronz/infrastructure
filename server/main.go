// Organized-Restfull is a simple api for authentication with golang
package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"flag"

	logrus "github.com/Sirupsen/logrus"
	"github.com/syronz/infrastructure/server/app"
	"github.com/syronz/infrastructure/server/errors"
	"github.com/syronz/infrastructure/server/controllers"


	"github.com/syronz/infrastructure/server/database/mysql"
	"github.com/syronz/infrastructure/server/routes"
	"github.com/syronz/infrastructure/server/dict"
	"github.com/syronz/infrastructure/server/utils/debug"
	"github.com/syronz/ozzo-routing"
	"github.com/syronz/ozzo-routing/access"
	"github.com/syronz/ozzo-routing/auth"
	"github.com/syronz/ozzo-routing/content"
	"github.com/syronz/ozzo-routing/cors"
)

func init(){
  logrus.SetFormatter(&logrus.JSONFormatter{})
  logrus.SetLevel(logrus.WarnLevel)
}

var migrate = flag.Bool("migrate", false, "Create tables inside database if not exist")
var reset = flag.Bool("reset", false, "Delete tables inside database if they are exist, this used after migrate")

func main() {



	// load application configuration
	if err := app.LoadConfig("./config"); err != nil {
		panic(fmt.Errorf("Invalid application configuration: %s", err))
	}

	// load error messages
	if err := errors.LoadMessages(app.Config.ErrorFile); err != nil {
		panic(fmt.Errorf("Failed to read the error message file: %s", err))
	}

	dict.LoadWords()
	app.LoadAcl()

	// Initiate connect with database
	mysql.DBTi.Connect()
	flag.Parse()
	if *migrate {
		mysql.DBTi.Migrate(*reset)
	}


	logger := logrus.New()

	//type DbPattern interface {
		//ShowIt()
	//}


	router := routing.New()
	router.Use(
		access.Logger(log.Printf),
	)

	api := router.Group("")
	api.Use(
		app.Init(logger),
		content.TypeNegotiator(content.JSON),

		cors.Handler(cors.Options{
			AllowOrigins: "*",
			AllowHeaders: "X-Requested-With, Content-Type, Authorization",
			AllowMethods: "GET, POST, PUT, HEAD, OPTIONS",
		}),
	)
	api.Post("/auth", controllers.Auth(app.Config.JWTSigningKey))
	//api.Use(controllers.AuthGuard)
	api.Use(auth.JWT(app.Config.JWTVerificationKey, auth.JWTOptions{
		SigningMethod: app.Config.JWTSigningMethod,
		TokenHandler: controllers.JWTHandler,
	}))
	routes.ServeCityResource(api)
	routes.ServeUserResource(api)
	routes.ServeAuthResource(api)
	routes.ServeActivityResource(api)
	routes.ServeDirectorResource(api)
	http.Handle("/", router)
	errServe := http.ListenAndServe(":8000", nil)
	if errServe != nil {
		debug.Log(errServe)
	}

}

func getCities(w http.ResponseWriter, r *http.Request) {
	logger := logrus.New()
	logger.Info("inside get Cities ................................")

	response := make(map[string]string)
	response["msg"] = "this is inside getCities"
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
}
