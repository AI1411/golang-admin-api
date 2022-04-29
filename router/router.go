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
	producthandler := handler.NewProductHandler(dbConn)

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
	products := r.Group("/products")
	{
		products.GET("", producthandler.GetAllProduct)
		products.GET("/:id", producthandler.GetProductDetail)
		products.POST("", producthandler.CreateProduct)
		products.PUT("/:id", producthandler.UpdateProduct)
		products.DELETE("/:id", producthandler.DeleteProduct)
	}
	auth := r.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.GET("/me", authHandler.Me)
	}

	r.Run()

	return r
}
