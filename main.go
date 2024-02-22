package main

import (
	"context"
	"fmt"
	"log"

	"github.com/IceySam/serve-soft/db"
	// "github.com/IceySam/serve-soft/examples"
	"github.com/joho/godotenv"
)

type car struct {
	Id    int
	Brand string
	Model string
	Year int
}

func main() {
	// mux := http.NewServeMux()
	// netHandler := network.NewNetwork(mux)
	// examples.Setup(netHandler)
	// http.ListenAndServe(":3000", mux)

	ENV, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal(err)
	}
	mysqlConStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", ENV["DB_USER"], ENV["DB_PASSWORD"], ENV["DB_HOST"], ENV["DB_PORT"], ENV["DB_NAME"])

	// conn, err := db.New("postgres", "postgres://postgres:S@mmy123@localhost:5432/sam")
	conn, err := db.New("mysql", mysqlConStr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	q := db.Query{Conn: conn}
	// err := q.Create("car", "id INT NOT NULL AUTO_INCREMENT", "brand VARCHAR(255)", "model VARCHAR(255)", "year INT", "PRIMARY KEY (id)")
	// lastId, err := q.Insert(&car{Brand: "Lambda", Model: "owl", Year: 2017})
	// lastId, err := q.InsertCtx(context.Background(), &car{Brand: "Sonata", Model: "brail", Year: 2020})
	// err = q.Update(&car{}).Set(map[string]any{"brand": "Lexus", "model": "lion"}).Where(map[string]any{
	// 	"id": 1,
	// }).Apply()
	// err = q.Update(&car{}).Set(map[string]any{"brand": "Lexus", "model": "lion"}).Where(map[string]any{
	// 	"id": 5,
	// }).ApplyCtx(context.Background())
	// err := q.Update(&car{}).Set(map[string]any{"brand": "Lexus", "model": "elephant"}).In("model", []interface{}{"Tiger", "viper"}).Apply()
	// err := q.Delete(&car{}).Where(map[string]any{"brand": "Toyota"}).Apply()
	// cars, err := q.FindAll(&car{})
	cars, err := q.FindAllCtx(context.Background(), &car{})
	// c := car{}
	// err = q.Find(&car{}).One(&c)
	// err = q.Find(&car{}).OneCtx(context.Background(), &c)
	// cars := make([]car, 0)
	// err = q.Find(&car{}).Where([]map[string]interface{}{{"year": 2020}, {"year": 2023}}).Many(&cars)
	// err = q.Find(&car{}).Where([]map[string]interface{}{{"year": 2020}, {"year": 2023}}).ManyCtx(context.Background(), &cars)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cars)
	// fmt.Println(lastId)
}
