package router

import (
	_ "app/docs"
	controller "app/internal/api/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"path/filepath"
	"runtime"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter() Router {
	return Router{engine: gin.Default()}
}

func (router *Router) Run(port string) error {
	router.Setup()
	return router.engine.Run(":" + port)
}

func (router *Router) Setup() {
	gin.SetMode(gin.DebugMode)
	_, currentFilePath, _, _ := runtime.Caller(1)
	templatesPath := filepath.Join(filepath.Dir(currentFilePath), "../../../templates")
	router.engine.LoadHTMLGlob(filepath.Join(templatesPath, "*.html"))
	router.engine.Static("/assets", filepath.Join(templatesPath, "assets"))
	router.engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))
	router.engine.GET("v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.engine.GET("/", controller.AuthPage)
	router.engine.GET("/forgot", controller.ForgotPage)
	router.engine.GET("newpass", controller.NewPassPage)
	router.engine.GET("home", controller.MainPage)
	v1 := router.engine.Group("/v1")
	router.AuthRoutes(v1)
}
