package handler

import (
	"fmt"
	"github.com/goark/koyomi"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AI1411/golang-admin-api/util/appcontext"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/signintech/gopdf"

	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/errors"
)

type OrderHandler struct {
	Db     *gorm.DB
	logger *zap.Logger
}

func NewOrderHandler(db *gorm.DB, logger *zap.Logger) *OrderHandler {
	return &OrderHandler{
		Db:     db,
		logger: logger,
	}
}

type searchOrderParams struct {
	UserID      string `form:"user_id" binding:"omitempty,uuid4"`
	Quantity    string `form:"quantity" binding:"omitempty,numeric"`
	TotalPrice  string `form:"total_price" binding:"omitempty,numeric"`
	OrderStatus string `form:"order_status" binding:"omitempty,oneof=new paid canceled delivered refunded returned partially partially_paid"`
	Offset      string `form:"offset,default=0" binding:"omitempty,numeric"`
	Limit       string `form:"limit,default=10" binding:"omitempty,numeric"`
}

type generatePDFRequest struct {
	OrderID string `json:"order_id" binding:"required,uuid4"`
}

func (h *OrderHandler) GetOrders(ctx *gin.Context) {
	traceID := appcontext.GetTraceID(ctx)
	var params searchOrderParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}

	var orders []models.Order
	query := createOrderQueryBuilder(params, h)
	query.Preload("OrderDetails").Find(&orders)

	ctx.JSON(http.StatusOK, gin.H{
		"total":  len(orders),
		"orders": orders,
	})
}

func (h *OrderHandler) GetOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	var order models.Order
	traceID := appcontext.GetTraceID(ctx)
	if err := h.Db.Preload("OrderDetails").Where("id = ?", id).First(&order).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			h.logger.Error("failed to get order", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.AbortWithStatusJSON(http.StatusNotFound, errors.NewNotFoundError("order not found"))
			return
		}
		h.logger.Error("failed to get order", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to get order", err))
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (h *OrderHandler) CreateOrder(ctx *gin.Context) {
	traceID := appcontext.GetTraceID(ctx)
	order := models.Order{}
	if err := ctx.ShouldBindJSON(&order); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	orderData := models.Order{
		ID:          order.CreateUUID(),
		UserID:      order.UserID,
		Quantity:    order.OrderDetails.TotalQuantity(),
		TotalPrice:  order.OrderDetails.TotalPrice(),
		OrderStatus: order.OrderStatus,
		Remarks:     order.Remarks,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := h.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("orders").Create(&orderData).Error; err != nil {
			h.logger.Error("failed to create order", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to create order", err))
			return err
		}
		for _, detail := range order.OrderDetails {
			detail.ID = order.CreateUUID()
			detail.OrderID = orderData.ID
			orderData.OrderDetails = append(orderData.OrderDetails, detail)
			if err := tx.Create(&detail).Error; err != nil {
				h.logger.Error("failed to create order", zap.Error(err),
					zap.String("trace_id", traceID))
				ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to create order detail", err))
				return err
			}
		}
		return nil
	}); err != nil {
		h.logger.Error("failed to create order", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to create order", err))
		return
	}

	ctx.JSON(http.StatusCreated, orderData)
}

func (h *OrderHandler) UpdateOrder(ctx *gin.Context) {
	order := models.Order{}
	id := ctx.Param("id")
	traceID := appcontext.GetTraceID(ctx)
	if err := h.Db.Where("id = ?", id).First(&order).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			h.logger.Error("failed to update milestone", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("order not found"))
		case gorm.ErrInvalidSQL:
			h.logger.Error("failed to update milestone", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	if err := ctx.ShouldBindJSON(&order); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	if err := h.Db.Save(&order).Error; err != nil {
		h.logger.Error("failed to update milestone", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to update order", err))
		return
	}
	ctx.JSON(http.StatusAccepted, order)
}

func (h *OrderHandler) DeleteOrder(ctx *gin.Context) {
	var order models.Order
	id := ctx.Param("id")
	if err := h.Db.Where("id = ?", id).First(&order).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("order not found"))
		case gorm.ErrInvalidSQL:
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}

	if err := h.Db.Delete(&order).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to delete order", err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (h *OrderHandler) ExportPDF(ctx *gin.Context) {
	traceID := appcontext.GetTraceID(ctx)
	var req generatePDFRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	pdf := gopdf.GoPdf{}
	A4 := *gopdf.PageSizeA4
	A4Yoko := gopdf.Rect{W: A4.H, H: A4.W}
	pdf.Start(gopdf.Config{PageSize: A4Yoko})
	pdf.AddPage()
	template := pdf.ImportPage("./template.pdf", 1, "/MediaBox")
	pageH := 1080 * (A4Yoko.W / 1920)
	pdf.UseImportedTemplate(template, 0, 0, A4Yoko.W, pageH)
	err := pdf.AddTTFFont("ipaexg", "./ipaexg.ttf")
	if err != nil {
		panic(err)
	}
	// 宛名
	pdf.SetFont("ipaexg", "", 28)

	var order models.Order
	if err := h.Db.Where("id = ?", req.OrderID).First(&order).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			h.logger.Error("failed to get order", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.AbortWithStatusJSON(http.StatusNotFound, errors.NewNotFoundError("order not found"))
			return
		}
		h.logger.Error("failed to get order", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to get order", err))
		return
	}

	var user models.User
	if err := h.Db.Where("id = ?", order.UserID).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			h.logger.Error("failed to get user", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.AbortWithStatusJSON(http.StatusNotFound, errors.NewNotFoundError("user not found"))
			return
		}
		h.logger.Error("failed to get user", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to get user", err))
		return
	}

	drawText(&pdf, 300, 140, user.LastName+user.FirstName)
	// 日付
	if err := pdf.SetFont("ipaexg", "", 15); err != nil {
		h.logger.Error("failed to set font", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to write pdf", err))
		return
	}
	year, month, day := convertKoyomi(order.CreatedAt)
	drawText(&pdf, 600, 100, year)  // 年
	drawText(&pdf, 635, 100, month) // 月
	drawText(&pdf, 676, 100, day)   // 日
	// 金額
	if err := pdf.SetFont("ipaexg", "", 28); err != nil {
		h.logger.Error("failed to set font", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to write pdf", err))
		return
	}
	totalPrice := convertPrice(int(order.TotalPrice))
	drawText(&pdf, 280, 200, "¥"+totalPrice+"-")
	// PDFをファイルに書き出す
	fileName := fmt.Sprintf("assets/pdf/%s_%s.pdf", order.ID, time.Now().Format("20060102150405"))
	if err := pdf.WritePdf(fileName); err != nil {
		h.logger.Error("failed to write pdf", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to write pdf", err))
		return
	}
}

func drawText(pdf *gopdf.GoPdf, x float64, y float64, s string) {
	pdf.SetX(x)
	pdf.SetY(y)
	pdf.Cell(nil, s)
}

func convertPrice(price int) string {
	arr := strings.Split(fmt.Sprintf("%d", price), "")
	cnt := len(arr) - 1
	res := ""
	i2 := 0
	for i := cnt; i >= 0; i-- {
		if i2 > 2 && i2%3 == 0 {
			res = fmt.Sprintf(",%s", res)
		}
		res = fmt.Sprintf("%s%s", arr[i], res)
		i2++
	}
	return res
}

func convertKoyomi(date time.Time) (year, month, day string) {
	tm := date
	te := koyomi.NewDate(tm)
	_, y := te.YearEraString()
	year = strings.Split(y, "年")[0]
	return year, strconv.Itoa(int(te.Month())), strconv.Itoa(te.Day())
}

func createOrderQueryBuilder(params searchOrderParams, h *OrderHandler) *gorm.DB {
	var orders []models.Order
	query := h.Db.Order("created_at desc").Find(&orders)
	if params.UserID != "" {
		query = query.Where("user_id = ?", params.UserID)
	}
	if params.Quantity != "" {
		query = query.Where("quantity = ?", params.Quantity)
	}
	if params.TotalPrice != "" {
		query = query.Where("total_price = ?", params.TotalPrice)
	}
	if params.OrderStatus != "" {
		query = query.Where("order_status = ?", params.OrderStatus)
	}
	if params.Offset != "" {
		query = query.Offset(params.Offset)
	}
	if params.Limit != "" {
		query = query.Limit(params.Limit)
	}
	return query
}
