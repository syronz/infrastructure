package routes

import (
	"github.com/syronz/infrastructure/server/controllers"

	"github.com/syronz/ozzo-routing"
)

// Rotes for users
func ServeUserResource(rg *routing.RouteGroup){
	r := &controllers.UserResource{}

	rg.Post("/users",controllers.AclChecker("users", "write"), r.Create)
	rg.Delete("/users/<id>",controllers.AclChecker("users", "write"), r.Delete)
	rg.Put("/users/<id>",controllers.AclChecker("users", "write"), r.Update)
	rg.Get("/users",controllers.AclChecker("users", "read"), r.List)
	rg.Get("/users/<id>",controllers.AclChecker("users", "read"), r.Get)

}
