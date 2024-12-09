package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Product struct {
	ID int
	NAME string
}

var products = []Product{{1, "CocaCola"}, {2, "Pepsi"}}

func main(){
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		fmt.Println("HELLO")
	})
	r.HandleFunc("/ProductsName/{key}", ProductHandlerName)
	r.HandleFunc("/ProductsID/{key}", ProductHandlerId)
	http.ListenAndServe(":8080", r)
}


func ProductHandlerName(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	name := vars["key"]
	for elem := range products{
		if products[elem].NAME == name{
			fmt.Printf("Продукт под номером %d\n", products[elem].ID)
		}
	}
}

func ProductHandlerId(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["key"])
	if err != nil{
		fmt.Println(err)
	}
	for elem := range products{
		if products[elem].ID == id{
			fmt.Printf("Продукт называется %s\n", products[elem].NAME)
		}
	}
}