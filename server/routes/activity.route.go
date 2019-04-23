package routes

import (
	"github.com/syronz/infrastructure/server/controllers"

	"github.com/syronz/ozzo-routing"
)

// Routes for activity
func ServeActivityResource(rg *routing.RouteGroup){
	r := &controllers.ActivityResource{}

	rg.Get("/activities", controllers.AclChecker("activities", "read"), r.List)

}
