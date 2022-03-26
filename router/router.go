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
	r.GET("/todos", todoHandler.GetAll)
	r.GET("/todos/:id", todoHandler.GetDetail)
	r.POST("/todos", todoHandler.CreateTodo)
	r.PUT("/todos/:id", todoHandler.UpdateTodo)
	r.DELETE("/todos/:id", todoHandler.DeleteTodo)

	r.GET("/users", userHandler.GetAllUser)
	r.GET("/users/:id", userHandler.GetUserDetail)
	r.POST("/users", userHandler.CreateUser)
	r.PUT("/users/:id", userHandler.UpdateUser)
	r.DELETE("/users/:id", userHandler.DeleteUser)
	r.PUT("/users/:id/uploadImage", userHandler.UploadUserImage)
	r.POST("/users/exportCsv", userHandler.ExportCSV)

	r.POST("register", authHandler.Register)
	r.POST("login", authHandler.Login)
	r.GET("me", authHandler.Me)

	r.Run()

	return r
}
