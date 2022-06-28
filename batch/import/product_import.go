package main

import (
	"encoding/csv"
	"fmt"
	"github.com/AI1411/golang-admin-api/db"
	"github.com/AI1411/golang-admin-api/models"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"log"
	"os"
	"strconv"
)

func importProductFromCsv() error {
	f, err := os.Open("product.CSV")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	r := csv.NewReader(transform.NewReader(f, japanese.ShiftJIS.NewDecoder()))
	records, err := r.ReadAll()
	for _, record := range records {
		product := newProduct(record[0], record[1], record[2], record[3], record[4])
		dbConn := db.Init()
		dbConn.Create(&product)
	}
	return nil
}

func newProduct(id, name, price, remarks, quantity string) error {
	_price, _ := strconv.Atoi(price)
	_quantity, _ := strconv.Atoi(quantity)
	product := models.Product{
		ID:          id,
		ProductName: name,
		Price:       uint(_price),
		Remarks:     remarks,
		Quantity:    _quantity,
	}
	dbConn := db.Init()
	if err := dbConn.Create(&product).Error; err != nil {
		return err
	}
	return nil
}

func main() {
	fmt.Println("商品一覧CSVインポートを開始します。")
	if err := importProductFromCsv(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("商品一覧CSVインポートを開始します。")
}
