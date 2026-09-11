package router

import (
	"home-box/backend/controller"
	"home-box/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		// 无需登录鉴权
		api.POST("/register", controller.Register)
		api.POST("/login", controller.Login)

		// 需要 JWT 鉴权
		auth := api.Group("")
		auth.Use(middleware.JWTAuth())
		{
			auth.GET("/profile", controller.GetProfile)
			auth.PUT("/profile", controller.UpdateProfile)
			auth.POST("/medicine", controller.AddMedicine)
			auth.PUT("/medicine", controller.UpdateMedicine)
			auth.DELETE("/medicine", controller.DeleteMedicine)
			auth.GET("/medicines", controller.GetMedicinesByUserID)
			auth.GET("/diseases", controller.GetDiseases)
		}
	}

	return r
}
