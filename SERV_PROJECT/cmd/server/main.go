package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"database/sql"


	_ "github.com/lib/pq"
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
	deserilizer := Product{}
	vars := mux.Vars(r)
	name := vars["key"]
	for elem := range products{
		if products[elem].NAME == name{
			jsonBytes, err := json.Marshal(&products[elem])
			if err != nil{
				log.Fatal(err)
			}
			w.Write(jsonBytes)
			log.Printf("Структура %v предоставленна в виде битов\n", jsonBytes)
			fmt.Printf("Продукт под номером %d\n", products[elem].ID)
			err = json.Unmarshal(jsonBytes, &deserilizer)
			if err != nil{
				log.Fatal(err)
			}
			log.Printf("Структурап ерезаписана в JSON формат %v\n", deserilizer)
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