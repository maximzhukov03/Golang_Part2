package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)
var db *sql.DB
type album struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Artist string  `json:"artist"`
    Price  float64 `json:"price"`
}

var albums []album

func initDB(){
	var err error
	connection := "host=127.0.0.1 port=5432 user=postgres dbname=movie_data sslmode=disable password=goLANG"
	db, err = sql.Open("postgres", connection)
	if err != nil{
		log.Println("Не удалось подключиться к базе данных")
	}
}

func main(){
	initDB()
	defer db.Close()

	r := gin.Default()
	
	r.GET("/albums", getAlbums)
	r.POST("/postalbum", postAlbums)
	r.GET("/albums/:id", getAlbumsId)
	r.DELETE("/deletealbum/:id", deleteAlbumsId)
	r.Run("localhost:8080")
}

func getAlbums(c *gin.Context) {
	albums, err := getAlbumsDb()
	if err != nil{
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "albums not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, albums)
	log.Println("GET ALL ALBUMS")
}

func getAlbumsDb() ([]album, error){
	rows, err := db.Query("SELECT * FROM movie_table")
	if err != nil{
		log.Println("Ошибка в запросе")
		return nil, err
	}

	albums := []album{}
	for rows.Next(){
		var a album
		if err := rows.Scan(&a.ID, &a.Title, &a.Artist, &a.Price); err != nil{
			return nil, err
		}
		albums = append(albums, a)
	}
	return albums, nil
}

func postAlbums(c *gin.Context) {
    var newAlbum album

	// Вызов BindJSON чтобы привязать полученный JSON 
	// к newAlbum.
    if err := c.BindJSON(&newAlbum); err != nil {
        return
    }
    postAlbumDb(newAlbum)
    c.IndentedJSON(http.StatusCreated, newAlbum)
	log.Println("POST ALBUM")
}

func postAlbumDb(newAlbum album){
	_, err := db.Exec("INSERT INTO movie_table (id, title, artist, price) VALUES ($1, $2, $3, $4)", newAlbum.ID, newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	if err != nil{
		log.Println("Ошибка в запросе")
		return
	}
}

func getAlbumsId(c *gin.Context){
	id := c.Param("id")
	albums, err := getAlbumsDb()
	if err != nil{
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "albums not found"})
		return
	}
	for _, elem := range albums{
		if elem.ID == id{
			log.Printf("GET ALBUM %s", elem.ID)
			c.IndentedJSON(http.StatusOK, elem)
			return
		}	
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}


func deleteAlbumsId(c *gin.Context){
	id := c.Param("id")
	deleteAlbumIdBD(db, id)
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func deleteAlbumIdBD(db *sql.DB, id string){
	_, err := db.Exec("DELETE FROM movie_table WHERE id = $1", id)
	if err != nil{
		log.Println("Ошибка в запросе")
	}
}