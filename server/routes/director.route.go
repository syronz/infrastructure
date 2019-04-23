package routes

import (
	"github.com/syronz/infrastructure/server/controllers"

	"github.com/syronz/ozzo-routing"
)

// All routes for director
func ServeDirectorResource(rg *routing.RouteGroup){
	r := &controllers.DirectorResource{}

	rg.Get("/directors/<id>", controllers.AclChecker("directors", "read"), r.Get)
	rg.Get("/directors", controllers.AclChecker("directors", "read"), r.List)
	rg.Post("/directors", controllers.AclChecker("directors", "write"), r.Create)
	rg.Delete("/directors/<id>", controllers.AclChecker("directors", "write"), r.Delete)
	rg.Put("/directors/<id>", controllers.AclChecker("directors", "write"), r.Update)
}

