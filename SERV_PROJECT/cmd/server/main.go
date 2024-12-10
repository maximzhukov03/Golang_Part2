package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Product struct {
	ID int
	NAME string
	PRICE int
}

func main(){
	connection := "host=127.0.0.1 port=5432 user=postgres dbname=product_data sslmode=disable password=goLANG"
	db, err := sql.Open("postgres", connection)
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil{
		log.Fatal(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/Products/", ProductHandle)
	r.HandleFunc("/ProductsID/{key}", ProductHandlerId)
	http.ListenAndServe(":8080", r)

}

func ProductHandle(w http.ResponseWriter, r *http.Request){
	switch r.Method{
		case http.MethodGet: 
		connection := "host=127.0.0.1 port=5432 user=postgres dbname=product_data sslmode=disable password=goLANG"
		db, err := sql.Open("postgres", connection)
		if err != nil{
			log.Fatal(err)
		}
		defer db.Close()
		if err := db.Ping(); err != nil{
			log.Fatal(err)
		}
		ProductHandlerGet(w, r, db)
		default: w.WriteHeader(http.StatusBadRequest)
	}
}


func ProductHandlerGet(w http.ResponseWriter, r *http.Request, db *sql.DB){
	product, err := ProductHandler(db)
	if err != nil{
		log.Fatal(err)
		return
	}
	jsonBytes, err := json.Marshal(product)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	w.Write(jsonBytes)
}

func ProductHandler(db *sql.DB) ([]Product, error){
	rows, err := db.Query("SELECT * FROM product_data")
	if err != nil{
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	products := make([]Product, 0)
	for rows.Next(){
		u := Product{}
		err := rows.Scan(&u.ID, &u.NAME, &u.PRICE)
		if err != nil{
			log.Fatal(err)
			return nil, err
		}
		products = append(products, u)
	}
	return products, nil
}

func ProductHandlerId(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["key"])
	if err != nil{
		fmt.Println(err)
	}
}