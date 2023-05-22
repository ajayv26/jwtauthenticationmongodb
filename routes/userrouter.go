package routes

import (
	"jwtauth/controller"
	"jwtauth/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.Use(middleware.Authenticate())
	r.GET("/users", controller.GetUsers())
	r.GET("/users/:user_id", controller.GetUser())
}
