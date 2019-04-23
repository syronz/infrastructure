package routes

import (
	"github.com/syronz/infrastructure/server/controllers"

	"github.com/syronz/ozzo-routing"
)

// All routes for city
func ServeCityResource(rg *routing.RouteGroup){
	r := &controllers.CityResource{}

	rg.Get("/cities/<id>", controllers.AclChecker("cities", "read"), r.Get)
	rg.Get("/cities", controllers.AclChecker("cities", "read"), r.List)
	rg.Post("/cities", controllers.AclChecker("cities", "write"), r.Create)
	rg.Delete("/cities/<id>", controllers.AclChecker("cities", "write"), r.Delete)
	rg.Put("/cities/<id>", controllers.AclChecker("cities", "write"), r.Update)
}
