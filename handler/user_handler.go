package handler

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/AI1411/golang-admin-api/util/appcontext"
	"go.uber.org/zap"

	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/errors"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type UserHandler struct {
	Db     *gorm.DB
	logger *zap.Logger
}

type searchUserParams struct {
	FirstName string `form:"first_name" binding:"omitempty,max=64"`
	LastName  string `form:"last_name" binding:"omitempty,max=64"`
	Age       string `form:"age" binding:"omitempty,numeric"`
	Email     string `form:"email" binding:"omitempty,max=64"`
	Offset    string `form:"offset,default=0" binding:"omitempty,numeric"`
	Limit     string `form:"limit,default=10" binding:"omitempty,numeric"`
}

func NewUserHandler(db *gorm.DB, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		Db:     db,
		logger: logger,
	}
}

func (h *UserHandler) GetAllUser(ctx *gin.Context) {
	traceID := appcontext.GetTraceID(ctx)
	var params searchUserParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind query params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}

	var users []models.User
	query := createUserQueryBuilder(params, h)
	if err := query.Preload("Todos").Find(&users).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to get users", err))
		return
	}
	log.Printf("param %+v", params)
	ctx.JSON(http.StatusOK, gin.H{
		"total": len(users),
		"users": users,
	})
}

func (h *UserHandler) GetUserDetail(ctx *gin.Context) {
	var user models.User
	id := ctx.Param("id")
	if err := h.Db.Preload("Todos").Where("id = ?", id).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("user not found"))
		case gorm.ErrInvalidSQL:
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	user := models.User{}
	id := ctx.Param("id")
	if err := h.Db.Where("id = ?", id).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("user not found"))
		case gorm.ErrInvalidSQL:
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}

	if err := h.Db.Save(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to update user", err))
		return
	}
	ctx.JSON(http.StatusAccepted, user)
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	user := models.User{}
	id := ctx.Param("id")
	if err := h.Db.Where("id = ?", id).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("product not found"))
		case gorm.ErrInvalidSQL:
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	if err := h.Db.Delete(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to delete user", err))
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{
		"message": "削除されました",
	})
}

func (h *UserHandler) ExportCSV(ctx *gin.Context) {
	fileName := time.Now().Format("202101011111") + "_users.csv"
	filePath := "assets/csv/users/" + fileName

	if err := h.CreateFile(filePath); err != nil {
		log.Printf("test %+v", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "CSVを出力しました",
	})
}

func (h *UserHandler) CreateFile(filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var users []models.User

	h.Db.Find(&users)

	if err := writer.Write([]string{
		"ID", "LastName", "FirstName", "Email", "Age",
	}); err != nil {
		return err
	}

	for _, user := range users {
		data := []string{
			user.ID,
			user.LastName,
			user.FirstName,
			user.Email,
			strconv.Itoa(int(user.Age)),
		}

		if err = writer.Write(data); err != nil {
			return err
		}
	}
	return nil
}

func createUserQueryBuilder(params searchUserParams, h *UserHandler) *gorm.DB {
	var users []models.User
	query := h.Db.Find(&users)

	if params.FirstName != "" {
		query = query.Where("first_name LIKE ?", "%"+params.FirstName+"%")
	}
	if params.LastName != "" {
		query = query.Where("last_name LIKE ?", "%"+params.LastName+"%")
	}
	if params.Age != "" {
		query = query.Where("age = ?", params.Age)
	}
	if params.Email != "" {
		query = query.Where("email = ?", params.Email)
	}
	if params.Offset != "" {
		query = query.Offset(params.Offset)
	}
	if params.Limit != "" {
		query = query.Limit(params.Limit)
	}
	return query
}
