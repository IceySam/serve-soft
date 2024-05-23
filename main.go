package main

import (
	"log"
	"net/http"

	"github.com/IceySam/serve-soft/examples"
	"github.com/IceySam/serve-soft/network"
)

type car struct {
	Id    int64
	Brand string
	Model string
	Year  int64
}

type container struct {
	Id        int64
	Reference string
	Name      string
}

func main() {
	mux := http.NewServeMux()
	netHandler := network.NewNetwork(mux)
	examples.Setup(netHandler)
	log.Println("Server started\nListening on port:3000...")
	http.ListenAndServe(":3000", mux)

	// ENV, err := godotenv.Read(".env")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// mysqlConStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", ENV["DB_USER"], ENV["DB_PASSWORD"], ENV["DB_HOST"], ENV["DB_PORT"], ENV["DB_NAME"])

	// // // conn, err := db.New("postgres", "postgres://postgres:S@mmy123@localhost:5432/sam")
	// conn, err := db.New("mysql", mysqlConStr)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer conn.Close()

	// q := db.Query{Conn: conn}
	// // err = q.Create("container", "id INT NOT NULL AUTO_INCREMENT", "reference VARCHAR(255)", "name VARCHAR(255) NULL", "PRIMARY KEY (id)")
	// // lastId, err := q.Insert(&container{Reference: "genral78"})
	// // lastId, err := q.Insert(&car{Brand: "Lambda", Model: "owl", Year: 2017})
	// // lastId, err := q.InsertCtx(context.Background(), &car{Brand: "Sonata", Model: "brail", Year: 2020})
	// // err = q.Update(&car{}).Set(map[string]any{"brand": "Lexus", "model": "lion"}).Where(map[string]any{
	// // 	"id": 1,
	// // }).Apply()
	// // err = q.Update(&car{}).Set(map[string]any{"brand": "Lexus", "model": "lion"}).Where(map[string]any{
	// // 	"id": 5,
	// // }).ApplyCtx(context.Background())
	// // tx, err := q.Conn.BeginTx(context.Background(), nil)
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }
	// // defer tx.Rollback()
	// // err = q.Update(&car{}).Set(map[string]interface{}{"brand": "Dune", "model": "lion"}).Where(map[string]any{
	// // 	"id": 5,
	// // }).TxApplyCtx(context.Background(), tx)
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }
	// // if err = tx.Commit(); err != nil {
	// // 	log.Fatal(err)
	// // }
	// // err := q.Update(&car{}).Set(map[string]any{"brand": "Lexus", "model": "elephant"}).In("model", []interface{}{"Tiger", "viper"}).Apply()
	// // err := q.Delete(&car{}).Where(map[string]any{"brand": "Toyota"}).Apply()
	// // cars, err := q.FindAll(&car{})
	// // cars, err := q.FindAllCtx(context.Background(), &car{})
	// // c := car{}
	// // err = q.Find(&car{}).One(&c)
	// // err = q.Find(&car{}).OneCtx(context.Background(), &c)
	// // cars := make([]car, 0)
	// // err = q.Find(&car{}).Where([]map[string]interface{}{{"year": 2020}, {"year": 2023}}).Many(&cars)
	// containers := make([]container, 0)
	// err = q.Find(&container{}).Many(&containers)
	// // err = q.Find(&car{}).Where(map[string]interface{}{"year": 2020}).ManyCtx(context.Background(), &cars)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // fmt.Println(containers)
	// fmt.Println(containers)
}
