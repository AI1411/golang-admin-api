package models

import (
	"crypto/sha1"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"path/filepath"
	"time"

	"github.com/olahol/go-imageupload"
)

const DefaultPasswordCost = 14

type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	FirstName string    `json:"first_name" binding:"required,max=16"`
	LastName  string    `json:"last_name" binding:"required,max=16"`
	Image     string    `json:"image" binding:"required,max=16"`
	Age       uint8     `json:"age" binding:"required,max=16"`
	Email     string    `json:"email" binding:"required,max=16"`
	Password  []byte    `json:"password" binding:"required,max=16"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Todos     []Todo    `json:"todos" binding:"omitempty"`
}

func (user *User) SetPassword(password string) {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), DefaultPasswordCost)

	user.Password = hashPassword
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}

func (user *User) UploadImage(ctx *gin.Context) {
	distDir := "./assets/images/users"
	img, err := imageupload.Process(ctx.Request, "image")
	if err != nil {
		panic(err)
	}

	thumb, err := imageupload.ThumbnailJPEG(img, 300, 300, 90)
	if err != nil {
		panic(err)
	}

	h := sha1.Sum(thumb.Data)
	imgName := fmt.Sprintf("%s_%x.jpg", time.Now().Format("20060102150405"), h[:4])
	savePath := filepath.Join(distDir, imgName)
	thumb.Save(savePath)

	user.Image = imgName
}
