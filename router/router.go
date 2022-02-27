package router

import (
	"api/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Router(dbConn *gorm.DB) {
	todoHandler := controllers.TodoHandler{Db: dbConn}
	userHandler := controllers.UserHandler{Db: dbConn}
	authHandler := controllers.AuthHandler{Db: dbConn}

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

	r.POST("register", authHandler.Register)
	if err := r.Run(":8084"); err != nil {
		return
	}
}
