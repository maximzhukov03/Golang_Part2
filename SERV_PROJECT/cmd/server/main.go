package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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

var prod Product

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
	r.HandleFunc("/ProductsID/{key}", ProductHandleId)
	r.HandleFunc("/ProductDel/{key}", ProductDeleteId)
	r.HandleFunc("/ProductPost/", ProductHendlerPost)
	http.ListenAndServe(":8080", r)

}

func ProductHandle(w http.ResponseWriter, r *http.Request){

		connection := "host=127.0.0.1 port=5432 user=postgres dbname=product_data sslmode=disable password=goLANG"
		db, err := sql.Open("postgres", connection)
		if err != nil{
			log.Fatal(err)
		}
		defer db.Close()
		if err := db.Ping(); err != nil{
			log.Fatal(err)
		}
		
		switch r.Method{
			case http.MethodGet: ProductHandlerGet(w, r, db)
			default: w.WriteHeader(http.StatusBadRequest)
	}
}

func ProductHandleId(w http.ResponseWriter, r *http.Request){
	connection := "host=127.0.0.1 port=5432 user=postgres dbname=product_data sslmode=disable password=goLANG"
	db, err := sql.Open("postgres", connection)
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil{
		log.Fatal(err)
	}
	switch r.Method{
		case http.MethodGet: ProductHandlerGetId(w, r, db)
		default: w.WriteHeader(http.StatusBadRequest)
	}
}

func ProductDeleteId(w http.ResponseWriter, r *http.Request){
	connection := "host=127.0.0.1 port=5432 user=postgres dbname=product_data sslmode=disable password=goLANG"
	db, err := sql.Open("postgres", connection)
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil{
		log.Fatal(err)
	}
	switch r.Method{
	case http.MethodDelete: ProductDelete(w, r, db)
	default: w.WriteHeader(http.StatusBadRequest)
	}
}

func ProductHendlerPost(w http.ResponseWriter, r *http.Request){
	connection := "host=127.0.0.1 port=5432 user=postgres dbname=product_data sslmode=disable password=goLANG"
	db, err := sql.Open("postgres", connection)
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil{
		log.Fatal(err)
	}
	switch r.Method{
	case http.MethodPost: ProductPost(w, r, db)
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

func ProductHandlerGetId(w http.ResponseWriter, r *http.Request, db *sql.DB){
	product, err := ProductHandlerId(w, r, db)
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

func ProductHandlerId(w http.ResponseWriter, r *http.Request, db *sql.DB) (Product, error){
	vars := mux.Vars(r)
	key_id, err := strconv.Atoi(vars["key"])
	if err != nil{
		fmt.Println(err)
	}
	rows, err := db.Query("SELECT * FROM product_data WHERE ID = $1", key_id)
	if err != nil{
		log.Fatal(err)
	}
	defer rows.Close()
	p := Product{}
	for rows.Next(){
		err = rows.Scan(&p.ID, &p.NAME, &p.PRICE)
		if err != nil{
			log.Fatal(err)
		}
	}
	return p, nil
}

func ProductDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	vars := mux.Vars(r)
	key_id, err := strconv.Atoi(vars["key"])
	if err != nil{
		log.Fatal(err)
	}
	_, err = db.Exec("DELETE FROM product_data WHERE ID = $1", key_id)
	if err != nil{
		log.Fatal(err)
	}
	log.Println("PRODUCT DELETE")
	return err
}

func ProductPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(err)
	}
	if err = json.Unmarshal(reqBytes, &prod); err != nil{
		w.WriteHeader(http.StatusBadRequest)
	}
	
	_, err = db.Exec("INSERT INTO product_data (id, name, price) VALUES ($1, $2, $3)", prod.ID, prod.NAME, prod.PRICE)
	if err != nil{
		log.Fatal("NOT POST")
	}
	log.Println("PRODUCT POST")
}