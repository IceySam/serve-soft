package main

import (
	"fmt"
	"log"

	"github.com/IceySam/serve-soft/db"
	"github.com/IceySam/serve-soft/test"
	"github.com/IceySam/serve-soft/utility"
	"github.com/joho/godotenv"
)

type car struct {
	Id    int
	Brand string
	Model string
	Year  int
}

func main() {
	// mux := http.NewServeMux()
	// netHandler := network.NewNetwork(mux)
	// test.Setup(netHandler)
	// http.ListenAndServe(":3000", mux)

	ENV, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal(err)
	}
	mysqlConStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", ENV["DB_USER"],ENV["DB_PASSWORD"],ENV["DB_HOST"],ENV["DB_PORT"], ENV["DB_NAME"])

	// conn := db.New("postgres", "postgres://postgres:S@mmy123@localhost:5432/sam")
	conn := db.New("mysql", mysqlConStr)
	defer conn.Close()

	q := test.Query{Conn: conn}
	// err := q.Create("car", "id INT NOT NULL AUTO_INCREMENT", "brand VARCHAR(255)", "model VARCHAR(255)", "year INT", "PRIMARY KEY (id)")
	// err := q.Insert(&car{Brand: "Lambda", Model: "owl", Year: 2017})
	// err := q.Update(&car{}).Set(map[string]any{"brand": "Lexus", "model": "lion"}).Where(map[string]any{
	// 	"id": 1,
	// }).Apply()
	// err := q.Update(&car{}).Set(map[string]any{"brand": "Lexus", "model": "elephant"}).In("model", []interface{}{"Tiger", "viper"}).Apply()
	// err := q.Delete(&car{}).Where(map[string]any{"brand": "Toyota"}).Apply()
	// cars, err := q.FindAll(&car{})
	// c := car{}
	// err := q.Find(&car{}).One(&c)
	items, err := q.Find(&car{}).Where([]map[string]interface{}{{"year": 2020}, {"year": 2023}}).Many()
	if err != nil {
		log.Fatal(err)
	}
	cars := make([]car, len(items))
	for i, v := range items {
		utility.ToStuct(v, &cars[i])
	}
	fmt.Println(cars)
}
