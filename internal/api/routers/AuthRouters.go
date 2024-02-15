package router

import (
	"app/internal/api/controllers"
	"github.com/gin-gonic/gin"
)

func (router *Router) AuthRoutes(group *gin.RouterGroup) {
	group.POST("/create", controller.UserCreate)
	group.POST("/login", controller.Login)
	group.POST("/forgot", controller.ForgotPassword)
	//group.GET("/newpass", controller.NewPassword)
	group.POST("/newpass", controller.PostNewPassword)
}
