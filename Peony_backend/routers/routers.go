package routers

import (
	"Peony/Peony_backend/info_controllers"
	"Peony/Peony_backend/jwt_middleware.go"
	"Peony/Peony_backend/user_controllers"

	"github.com/gin-gonic/gin"
)

func SetUserRouter(router *gin.Engine) {
	user_router := router.Group("/user")
	user_router.GET("/", jwt_middleware.JwtAuth(), user_controllers.UserDetail)
	user_router.GET("/auth", user_controllers.AuthHandler)
	user_router.GET("/redir", user_controllers.UserGmail)
	user_router.GET("/login")
	user_router.POST("/", user_controllers.CreateUser)
	user_router.PUT("/:user_id/")
}

func SetInfoRouter(router *gin.Engine) {
	info_router := router.Group("/info/")
	info_router.GET("/")
	info_router.POST("/", jwt_middleware.JwtAuth(), info_controllers.CreateInfo)
	info_router.PUT("/:info_id/")
	info_router.DELETE("/:info_id/")
}

func InitRoute() *gin.Engine {
	router := gin.Default()

	SetUserRouter(router)
	SetInfoRouter(router)
	return router
}
