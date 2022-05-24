package main

import (
	"encoding/csv"
	"fmt"
	"github.com/AI1411/golang-admin-api/util"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/AI1411/golang-admin-api/db"
	"github.com/AI1411/golang-admin-api/models"
)

func export(args []string) error {
	app := &cli.App{
		Name:  "注文一覧をCSVにエクスポートする",
		Usage: "本日分の注文一覧をCSVにエクスポートする",
		Action: func(c *cli.Context) error {
			year, month, day := time.Now().Date()
			dateStr := fmt.Sprintf("%d%02d%02d", year, month, day)
			fileName := dateStr + "_orders.csv"
			fileDir := "assets/csv/orders/" + strconv.Itoa(year) +
				"/" + strconv.Itoa(int(month)) + "/" + strconv.Itoa(day) + "/"

			util.CheckDir(fileDir)

			filePath := fileDir + fileName

			if err := createCsvFile(filePath); err != nil {
				return err
			}

			return nil
		},
	}

	err := app.Run(args)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func createCsvFile(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	var orders []models.Order
	dbConn := db.Init()

	if err = dbConn.Find(&orders).Error; err != nil {
		return err
	}

	if err := writer.Write([]string{
		"注文ID",
		"ユーザID",
		"数量",
		"合計金額",
		"注文ステータス",
		"注文備考",
		"作成日時",
		"更新日時",
	}); err != nil {
		return err
	}

	for _, order := range orders {
		if err := writer.Write([]string{
			order.ID,
			strconv.Itoa(int(order.UserID)),
			strconv.Itoa(int(order.Quantity)),
			strconv.Itoa(int(order.TotalPrice)),
			string(order.OrderStatus),
			order.Remarks,
			order.CreatedAt.Format("2006-01-02 15:04:05"),
			order.UpdatedAt.Format("2006-01-02 15:04:05"),
		}); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	fmt.Println("注文一覧CSVエクスポートを開始します。")
	if err := export(os.Args); err != nil {
		log.Fatal(err)
	}
	fmt.Println("注文一覧CSVエクスポート終了しました。")
}
