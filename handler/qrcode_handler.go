package handler

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/skip2/go-qrcode"

	"github.com/AI1411/golang-admin-api/util/errors"
)

type QrcodeHandler struct {
	Db *gorm.DB
}

func NewQrcodeHandler(db *gorm.DB) *QrcodeHandler {
	return &QrcodeHandler{Db: db}
}

func (h *QrcodeHandler) GenerateQrcode(ctx *gin.Context) {
	filePath := "./assets/qrcode/"
	if err := checkDir(filePath); err != nil {
		log.Println(err)
		return
	}
	fileName := "qrcode.png"
	if err := qrcode.WriteFile("https://emaple.org", qrcode.Medium, 256, filePath+fileName); err != nil {
		log.Fatal(err)
	}
	return
}

func checkDir(fileDir string) error {
	_, err := os.Stat(fileDir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(fileDir, 0o777); err != nil {
			return errors.NewInternalServerError("failed to create directory", err)
		}
		return errors.New("ディレクトリを作成しました。")
	}

	if err != nil {
		return err
	}
	return nil
}
