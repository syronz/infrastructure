package routes

import (
	"github.com/syronz/infrastructure/server/controllers"

	"github.com/syronz/ozzo-routing"
)

// Rotes for customers
func ServeCustomerResource(rg *routing.RouteGroup){
	r := &controllers.CustomerResource{}

	rg.Post("/customers",controllers.AclChecker("customers", "write"), r.Create)
	rg.Delete("/customers/<id>",controllers.AclChecker("customers", "write"), r.Delete)
	rg.Put("/customers/<id>",controllers.AclChecker("customers", "write"), r.Update)
	rg.Get("/customers",controllers.AclChecker("customers", "read"), r.List)
	rg.Get("/customers/<id>",controllers.AclChecker("customers", "read"), r.Get)

}
