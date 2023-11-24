package users

import "github.com/gin-gonic/gin"

func UserRoutes(router *gin.Engine) {
	authRoutes := router.Group("/user/auth")
	{
		authRoutes.POST("/signup", Register)
		authRoutes.POST("/signin", Login)
		authRoutes.POST("/logout", JWTAuthMiddleWare(), Logoutuser)
		authRoutes.GET("/getuser", JWTAuthMiddleWare(), GetSingleUser)
		authRoutes.GET("/fetchuser", FetchSingleUser)
		authRoutes.POST("/updateuser", JWTAuthMiddleWare(), UpdateUser)
		authRoutes.GET("/fetchsellers", FetchSellers)
	}
}
