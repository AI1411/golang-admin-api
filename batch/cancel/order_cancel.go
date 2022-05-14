package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/urfave/cli/v2"

	"github.com/AI1411/golang-admin-api/db"
	"github.com/AI1411/golang-admin-api/models"
)

func run(args []string) error {
	app := &cli.App{
		Name:  "注文キャンセルバッチ",
		Usage: "注文をキャンセルする",
		Action: func(c *cli.Context) error {
			var orders []models.Order
			dbConn := db.Init()
			if err := dbConn.Transaction(func(tx *gorm.DB) error {
				tx.Where("order_status = ?", models.OrderStatusNew).
					Preload("OrderDetails").
					Find(&orders).
					Updates(models.Order{
						OrderStatus: models.OrderStatusCanceled,
					})

				fmt.Printf("注文を %d 件キャンセルしました\n", len(orders))

				for _, order := range orders {
					for _, detail := range order.OrderDetails {
						tx.Where("order_id = ?", order.ID).
							Find(&detail).
							Update("order_detail_status", models.OrderDetailStatusCanceled)
					}
					fmt.Printf("注文ID%sの注文詳細を %d 件キャンセルしました\n", order.ID, len(order.OrderDetails))
				}
				return nil
			}); err != nil {
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

func main() {
	fmt.Println("注文キャンセルバッチを開始します。")
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
	fmt.Println("注文キャンセルバッチを終了します。")
}
