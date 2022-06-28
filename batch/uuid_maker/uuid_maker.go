package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

func makeUUID() error {
	app := &cli.App{
		Name: "UUIDを引数分生成する",
		Action: func(c *cli.Context) error {
			n, _ := strconv.Atoi(c.Args().Get(0))
			for i := 0; i < n; i++ {
				newUuid, _ := uuid.NewRandom()
				fmt.Println(newUuid.String())
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func main() {
	fmt.Println("uuid_makerを開始します。")
	if err := makeUUID(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("uuid_makerを終了します。")
}
