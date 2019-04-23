package routes

import (
	"github.com/syronz/infrastructure/server/controllers"
	"github.com/syronz/ozzo-routing"
)

// Routes for auth
func ServeAuthResource(rg *routing.RouteGroup){
	r := &controllers.AuthResource{}

	rg.Get("/auth/profile", r.GetProfile)
	rg.Put("/auth/profile", r.UpdateProfile)
	rg.Post("/auth/avatar", r.UplodadAvatar)
}
