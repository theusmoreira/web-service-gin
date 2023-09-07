package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type album struct {
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type Album struct {
	gorm.Model
	Title  string
	Artist string
	Price  float64
}

// var albums = []album{
// 	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
// 	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
// 	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
// }

func getAlbums(c *gin.Context, db *gorm.DB) {
	albums := []Album{}

	db.Find(&albums)

	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context, db *gorm.DB) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// verify if album already exists with same title and artist
	var album Album
	db.Where("title = ? AND artist = ?", newAlbum.Title, newAlbum.Artist).First(&album)

	if album.ID != 0 {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "album already exists"})
		return
	}

	db.Create(&Album{Title: newAlbum.Title, Artist: newAlbum.Artist, Price: newAlbum.Price})

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	var album Album = Album{}
	db.First(&album, id)

	if album.ID != 0 {
		c.IndentedJSON(http.StatusOK, &album)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func connectWithDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("albums.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Album{})

	return db
}

func main() {
	db := connectWithDB()

	router := gin.Default()
	router.GET("/albums", (func(c *gin.Context) { getAlbums(c, db) }))
	router.GET("/albums/:id", (func(c *gin.Context) { getAlbumByID(c, db) }))
	router.POST("/albums", (func(c *gin.Context) { postAlbums(c, db) }))

	router.Run("localhost:8080")
}
