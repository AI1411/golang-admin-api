package router

import (
	"log"

	"github.com/AI1411/golang-admin-api/db"
	"github.com/AI1411/golang-admin-api/handler"
	"github.com/AI1411/golang-admin-api/middleware"
	"github.com/AI1411/golang-admin-api/models"
	logger "github.com/AI1411/golang-admin-api/server"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Router() *gin.Engine {
	zapLogger, err := logger.NewLogger(false)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = zapLogger.Sync() }()

	dbConn := db.Init()
	uuidGen := &models.RandomUUIDGenerator{}
	todoHandler := handler.NewTodoHandler(dbConn, zapLogger)
	userHandler := handler.NewUserHandler(dbConn, zapLogger)
	authHandler := handler.NewAuthHandler(dbConn, zapLogger)
	productHandler := handler.NewProductHandler(dbConn, uuidGen, zapLogger)
	orderHandler := handler.NewOrderHandler(dbConn, zapLogger)
	orderDetailHandler := handler.NewOrderDetailHandler(dbConn, zapLogger)
	couponHandler := handler.NewCouponHandler(dbConn, zapLogger)
	qrcodeHandler := handler.NewQrcodeHandler(dbConn)
	userGroupHandler := handler.NewUserGroupHandler(dbConn)
	milestoneHandler := handler.NewMilestoneHandler(dbConn, zapLogger)
	epicHandler := handler.NewEpicHandler(dbConn, zapLogger)
	projectHandler := handler.NewProjectHandler(dbConn, uuidGen, zapLogger)

	r := gin.Default()

	r.Use(middleware.Cors())
	r.Use(func(_ *gin.Context) { binding.EnableDecoderUseNumber = true })
	r.Use(middleware.NewTracing())
	r.Use(middleware.NewLogging(zapLogger))
	authorized := r.Group("/")
	authorized.Use(middleware.AuthenticateBearer())
	auth := r.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
		auth.POST("/register", authHandler.Register)
		authorized.POST("/auth/me", authHandler.Me)
	}
	todos := authorized.Group("/todos")
	{
		todos.GET("", todoHandler.GetAll)
		todos.GET("/:id", todoHandler.GetDetail)
		todos.POST("", todoHandler.CreateTodo)
		todos.PUT("/:id", todoHandler.UpdateTodo)
		todos.DELETE("/:id", todoHandler.DeleteTodo)
	}
	users := authorized.Group("/users")
	{
		users.GET("", userHandler.GetAllUser)
		users.GET("/:id", userHandler.GetUserDetail)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
		users.POST("/exportCsv", userHandler.ExportCSV)
	}
	products := authorized.Group("/products")
	{
		products.GET("", productHandler.GetAllProduct)
		products.GET("/:id", productHandler.GetProductDetail)
		products.POST("", productHandler.CreateProduct)
		products.PUT("/:id", productHandler.UpdateProduct)
		products.DELETE("/:id", productHandler.DeleteProduct)
	}
	orders := authorized.Group("/orders")
	{
		orders.POST("", orderHandler.CreateOrder)
		orders.GET("", orderHandler.GetOrders)
		orders.GET("/:id", orderHandler.GetOrder)
		orders.PUT("/:id", orderHandler.UpdateOrder)
		orders.DELETE("/:id", orderHandler.DeleteOrder)
	}
	orderDetails := authorized.Group("/orderDetails")
	{
		orderDetails.GET("/:id", orderDetailHandler.GetOrderDetail)
		orderDetails.POST("", orderDetailHandler.CreateOrderDetail)
		orderDetails.PUT("/:id", orderDetailHandler.UpdateOrderDetail)
		orderDetails.DELETE("/:id", orderDetailHandler.DeleteOrderDetail)
	}
	coupons := authorized.Group("/coupons")
	{
		coupons.GET("/", couponHandler.GetAllCoupon)
		coupons.GET("/:id", couponHandler.GetCouponDetail)
		coupons.POST("", couponHandler.CreateCoupon)
		coupons.PUT("/:id", couponHandler.UpdateCoupon)
		coupons.POST("/:coupon_id/users/:user_id", couponHandler.AcquireCoupon)
	}
	userGroups := authorized.Group("/userGroups")
	{
		userGroups.GET("", userGroupHandler.GetAllUserGroups)
		userGroups.GET("/:id", userGroupHandler.GetUserGroupsDetail)
		userGroups.POST("", userGroupHandler.CreateUserGroup)
	}
	milestones := authorized.Group("/milestones")
	{
		milestones.GET("", milestoneHandler.GetMilestones)
		milestones.GET("/:id", milestoneHandler.GetMilestoneDetail)
		milestones.POST("", milestoneHandler.CreateMilestone)
		milestones.PUT("/:id", milestoneHandler.UpdateMileStone)
	}
	epics := authorized.Group("/epics")
	{
		epics.GET("", epicHandler.GetEpics)
		epics.GET("/:id", epicHandler.GetEpicDetail)
		epics.POST("", epicHandler.CreateEpic)
		epics.PUT("/:id", epicHandler.UpdateEpic)
		epics.DELETE("/:id", epicHandler.DeleteEpic)
	}
	projects := authorized.Group("/projects")
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
