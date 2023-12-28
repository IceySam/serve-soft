package main

import (
	"context"
	"fmt"

	"github.com/IceySam/serve-soft/db"
	"github.com/IceySam/serve-soft/test"
)

type car struct {
	Brand string
	Model string
	Year  int
}

func main() {
	// mux := http.NewServeMux()
	// netHandler := network.NewNetwork(mux)
	// test.Setup(netHandler)
	// http.ListenAndServe(":3000", mux)

	conn := db.New("postgres://postgres:S@mmy123@localhost:5432/sam")
	defer conn.Close(context.Background())

	q := test.Query{Conn: conn}
	// err := q.Create("car", "brand VARCHAR(255)", "model VARCHAR(255)", "year INT")
	// err := q.Insert(&car{Brand: "Toyota", Model: "jaguar", Year: 2023})
	// err := q.Update(&car{}).Set(map[string]any{"brand": "Lexus", "model": "antelope"}).Where(map[string]any{
	// 	"model": "antelope",
	// }).Apply()
	// err := q.Update(&car{}).Set(map[string]any{"brand": "Lexus", "model": "elephant"}).In("model", []interface{}{"Tiger", "antelope"}).Apply()
	// err := q.Delete(&car{}).Where(map[string]any{"brand": "Toyota"}).Apply()
	// cars, err := q.FindAll(&car{})
	cars, err := q.Find(&car{}).Where([]map[string]interface{}{{"year": 2020}, {"year": 2023}}).Many()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cars)
}
