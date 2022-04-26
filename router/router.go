package router

import (
	"github.com/AI1411/golang-admin-api/db"
	"github.com/AI1411/golang-admin-api/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	dbConn := db.Init()
	todoHandler := handler.NewTodoHandler(dbConn)
	userHandler := handler.NewUserHandler(dbConn)
	authHandler := handler.NewAuthHandler(dbConn)

	r := gin.Default()
	r.Use(cors.Default())
	todos := r.Group("/todos")
	{
		todos.GET("", todoHandler.GetAll)
		todos.GET("/:id", todoHandler.GetDetail)
		todos.POST("", todoHandler.CreateTodo)
		todos.PUT("/:id", todoHandler.UpdateTodo)
		todos.DELETE("/:id", todoHandler.DeleteTodo)
	}
	users := r.Group("/users")
	{
		users.GET("", userHandler.GetAllUser)
		users.GET("/:id", userHandler.GetUserDetail)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
		users.PUT("/:id/uploadImage", userHandler.UploadUserImage)
		users.POST("/exportCsv", userHandler.ExportCSV)
	}

	r.POST("register", authHandler.Register)
	r.POST("login", authHandler.Login)
	r.GET("me", authHandler.Me)

	r.Run()

	return r
}
