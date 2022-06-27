package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func importProductFromCsv() error {
	f, err := os.Open("33OKAYAM 2.CSV")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	r := csv.NewReader(transform.NewReader(f, japanese.ShiftJIS.NewDecoder()))
	for {
		records, err := r.Read()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(records)
	}
}

func main() {
	fmt.Println("商品一覧CSVインポートを開始します。")
	if err := importProductFromCsv(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("商品一覧CSVインポートを開始します。")
}
