package main

import (
	"meeting-backend/config"
	"meeting-backend/controllers"
	"meeting-backend/middlewares"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	r := gin.Default()

	public := r.Group("/api")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	public.POST("/login", controllers.Login)
	public.GET("/user", controllers.CurrentUser)

	meetings := public.Group("/meetings")
	meetings.Use(middlewares.JwtAuthMiddleWare())

	meetings.GET("/", controllers.GetMeetings)
	meetings.GET("/:id", controllers.GetMeetingByID)
	meetings.POST("/", controllers.PostMeeting)
	meetings.PUT("/:id", controllers.PutMeeting)
	meetings.DELETE("/:id", controllers.DeleteMeeting)

	staffs := public.Group("/staffs")
	staffs.Use(middlewares.JwtAuthMiddleWare())

	staffs.GET("/", controllers.GetStaffs)
	staffs.GET("/:id", controllers.GetStaffByID)
	staffs.POST("/", controllers.PostStaff)
	staffs.PUT("/:id", controllers.PutStaff)
	staffs.DELETE("/:id", controllers.DeleteStaff)

	locations := public.Group("/locations")
	locations.Use(middlewares.JwtAuthMiddleWare())

	locations.GET("/", controllers.GetLocations)
	locations.GET("/:id", controllers.GetLocationByID)
	locations.POST("/", controllers.PostLocation)
	locations.PUT("/:id", controllers.PutLocation)
	locations.DELETE("/:id", controllers.DeleteLocation)

	r.Run(":8080")
}
