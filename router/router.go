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
	productHandler := handler.NewProductHandler(dbConn)
	orderHandler := handler.NewOrderHandler(dbConn)
	orderDetailHandler := handler.NewOrderDetailHandler(dbConn)

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
		products.GET("", productHandler.GetAllProduct)
		products.GET("/:id", productHandler.GetProductDetail)
		products.POST("", productHandler.CreateProduct)
		products.PUT("/:id", productHandler.UpdateProduct)
		products.DELETE("/:id", productHandler.DeleteProduct)
	}
	auth := r.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.GET("/me", authHandler.Me)
	}
	orders := r.Group("/orders")
	{
		orders.POST("", orderHandler.CreateOrder)
	}
	orderDetails := r.Group("/orderDetails")
	{
		orderDetails.POST("", orderDetailHandler.CreateOrderDetail)
	}

	r.Run()

	return r
}
