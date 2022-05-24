package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/skip2/go-qrcode"
	"log"
	"os"
)

type QrcodeHandler struct {
	Db *gorm.DB
}

func NewQrcodeHandler(db *gorm.DB) *QrcodeHandler {
	return &QrcodeHandler{Db: db}
}

func (h *QrcodeHandler) GenerateQrcode(ctx *gin.Context) {
	filePath := "./assets/qrcode/"
	checkDir(filePath)
	fileName := "qrcode.png"
	if err := qrcode.WriteFile("https://emaple.org", qrcode.Medium, 256, filePath+fileName); err != nil {
		log.Fatal(err)
	}
	return
}

func checkDir(fileDir string) error {
	_, err := os.Stat(fileDir)
	if os.IsNotExist(err) {
		os.MkdirAll(fileDir, 0777)
		return errors.New("ディレクトリを作成しました。")
	}

	if err != nil {
		return err
	}
	return nil
}
