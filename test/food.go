package test

import (
	"fmt"
	"net/http"

	"github.com/IceySam/network/network"
)

var foods [5]string = [5]string{"yam", "bean", "rice", "oil", "beef"} 

func getOne(index int) string {
	return foods[index]
}

func getAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println(foods)
}

func Setup(h network.NetHandler) {
	// pkg2 := network.Handler{
	// 	Item: "/food",
	// 	Initiate: getAll,
	// }

}