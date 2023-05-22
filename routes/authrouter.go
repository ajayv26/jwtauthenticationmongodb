package routes

import (
	"jwtauth/controller"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("users/signup", controller.SignUp())
	r.POST("users/login", controller.Login())
}
