package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/AI1411/golang-admin-api/db"
	"github.com/AI1411/golang-admin-api/models"
)

func run(args []string) error {
	app := &cli.App{
		Name:  "注文キャンセルバッチ",
		Usage: "注文をキャンセルする",
		Action: func(c *cli.Context) error {
			dbConn := db.Init()
			dbConn.Table("orders").
				Select("order_status").
				Where("order_status = ?", models.OrderStatusNew).
				Updates(models.Order{
					OrderStatus: "cancelled",
				})
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
	run(os.Args)
}
