package main

import (
	//"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type album struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Artist string  `json:"artist"`
    Price  float64 `json:"price"`
}

var albums = []album{
    {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
    {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
    {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main(){
	r := gin.Default()
	r.GET("/albums", getAlbums)
	r.POST("/postalbum", postAlbums)
	r.GET("/albums/:id", getAlbumsId)
	r.DELETE("/deletealbum/:id", deleteAlbumsId)
	r.Run("localhost:8080")
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
	log.Println("GET ALL ALBUMS")
}

func postAlbums(c *gin.Context) {
    var newAlbum album

	// Вызов BindJSON чтобы привязать полученный JSON 
	// к newAlbum.
    if err := c.BindJSON(&newAlbum); err != nil {
        return
    }
    albums = append(albums, newAlbum)
    c.IndentedJSON(http.StatusCreated, newAlbum)
	log.Println("POST ALBUM")
}

func getAlbumsId(c *gin.Context){
	id := c.Param("id")

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

	for index, elem := range albums{
		if elem.ID == id{
			albums = append(albums[:index], albums[index+1:]...)
			c.IndentedJSON(http.StatusOK, elem)
			log.Printf("DELETED ALBUM %s", elem.ID)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}