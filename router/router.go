package router

import (
	"github.com/AI1411/golang-admin-api/db"
	"github.com/AI1411/golang-admin-api/handler"
	"github.com/AI1411/golang-admin-api/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	dbConn := db.Init()
	uuidGen := &models.RandomUUIDGenerator{}
	todoHandler := handler.NewTodoHandler(dbConn)
	userHandler := handler.NewUserHandler(dbConn)
	authHandler := handler.NewAuthHandler(dbConn)
	productHandler := handler.NewProductHandler(dbConn)
	orderHandler := handler.NewOrderHandler(dbConn)
	orderDetailHandler := handler.NewOrderDetailHandler(dbConn)
	couponHandler := handler.NewCouponHandler(dbConn)
	qrcodeHandler := handler.NewQrcodeHandler(dbConn)
	userGroupHandler := handler.NewUserGroupHandler(dbConn)
	milestoneHandler := handler.NewMilestoneHandler(dbConn)
	epicHandler := handler.NewEpicHandler(dbConn)
	projectHandler := handler.NewProjectHandler(dbConn, uuidGen)

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
		orders.GET("", orderHandler.GetOrders)
		orders.GET("/:id", orderHandler.GetOrder)
		orders.PUT("/:id", orderHandler.UpdateOrder)
		orders.DELETE("/:id", orderHandler.DeleteOrder)
	}
	orderDetails := r.Group("/orderDetails")
	{
		orderDetails.GET("/:id", orderDetailHandler.GetOrderDetail)
		orderDetails.POST("", orderDetailHandler.CreateOrderDetail)
		orderDetails.PUT("/:id", orderDetailHandler.UpdateOrderDetail)
		orderDetails.DELETE("/:id", orderDetailHandler.DeleteOrderDetail)
	}
	coupons := r.Group("/coupons")
	{
		coupons.GET("/", couponHandler.GetAllCoupon)
		coupons.GET("/:id", couponHandler.GetCouponDetail)
		coupons.POST("", couponHandler.CreateCoupon)
		coupons.PUT("/:id", couponHandler.UpdateCoupon)
		coupons.POST("/:coupon_id/users/:user_id", couponHandler.AcquireCoupon)
	}
	userGroups := r.Group("/userGroups")
	{
		userGroups.GET("", userGroupHandler.GetAllUserGroups)
		userGroups.GET("/:id", userGroupHandler.GetUserGroupsDetail)
		userGroups.POST("", userGroupHandler.CreateUserGroup)
	}
	milestones := r.Group("/milestones")
	{
		milestones.GET("", milestoneHandler.GetMilestones)
		milestones.GET("/:id", milestoneHandler.GetMilestoneDetail)
		milestones.POST("", milestoneHandler.CreateMilestone)
		milestones.PUT("/:id", milestoneHandler.UpdateMileStone)
	}
	epics := r.Group("/epics")
	{
		epics.GET("", epicHandler.GetEpics)
		epics.GET("/:id", epicHandler.GetEpicDetail)
		epics.POST("", epicHandler.CreateEpic)
		epics.PUT("/:id", epicHandler.UpdateEpic)
		epics.DELETE("/:id", epicHandler.DeleteEpic)
	}
	projects := r.Group("/projects")
	{
		projects.GET("", projectHandler.GetProjects)
		projects.GET("/:id", projectHandler.GetProjectDetail)
		projects.POST("", projectHandler.CreateProject)
		projects.PUT("/:id", projectHandler.UpdateProject)
		projects.DELETE("/:id", projectHandler.DeleteProject)
	}

	r.GET("/qrcode", qrcodeHandler.GenerateQrcode)

	if err := r.Run(); err != nil {
		panic(err)
	}

	return r
}
