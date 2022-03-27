package handler

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/errors"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type UserHandler struct {
	Db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{Db: db}
}

func (h *UserHandler) GetAllUser(ctx *gin.Context) {
	var users []models.User
	h.Db.Preload("Todos").Find(&users)
	ctx.JSON(http.StatusOK, gin.H{
		"total": len(users),
		"users": users,
	})
}

func (h *UserHandler) GetUserDetail(ctx *gin.Context) {
	var user models.User
	id := ctx.Param("id")
	h.Db.Preload("Todos").Where("id = ?", id).Find(&user)

	ctx.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	user := models.User{}
	id := ctx.Param("id")
	h.Db.First(&user, id)
	if err := ctx.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid request")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	h.Db.Save(&user)
	ctx.JSON(http.StatusAccepted, user)
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	user := models.User{}
	id := ctx.Param("id")
	if err := h.Db.First(&user, id).Error; err != nil {
		restErr := errors.NewBadRequestError("invalid request")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	h.Db.Delete(&user)
	ctx.JSON(http.StatusNoContent, gin.H{
		"message": "削除されました",
	})
}

func (h *UserHandler) UploadUserImage(ctx *gin.Context) {
	user := models.User{}
	id := ctx.Param("id")
	if err := h.Db.First(&user, id).Error; err != nil {
		restErr := errors.NewBadRequestError("invalid request")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	user.UploadImage(ctx)
	h.Db.Save(&user)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "画像がアップロードされました",
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

	writer.Write([]string{
		"ID", "LastName", "FirstName", "Email", "Age",
	})

	for _, user := range users {
		data := []string{
			strconv.Itoa(user.ID),
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
