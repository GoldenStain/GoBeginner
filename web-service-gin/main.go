package main

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var db *sql.DB

func main() {
	connectDB()

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	var albums []album
	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		log.Println(err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			log.Println(err)
		}
	}()
	for rows.Next() {
		var alb album
		if err = rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			log.Println(err)
		}
		albums = append(albums, alb)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
	}
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album
	var err error
	if err = c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid post request error",
			"message": err.Error()})
		return
	}
	// Add the new album to the slice.
	// check for errors
	if err = validateAlbum(newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newAlbum.ID != "" {
		id, err := strconv.Atoi(newAlbum.ID)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		_, err = db.Exec("INSERT INTO album (id, title, artist, price) VALUES (?, ?, ?, ?)", id, newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	} else {
		_, err = db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	}
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func validateAlbum(album album) error {
	if album.ID == "" {
		return nil
	}
	if album.Title == "" {
		return errors.New("Invalid album title")
	}
	if album.Artist == "" {
		return errors.New("Invalid album artist")
	}
	if album.Price <= 0 {
		return errors.New("Invalid album price")
	}
	var exists bool
	id, err := strconv.Atoi(album.ID)
	if err != nil {
		return errors.New("Invalid album id")
	}
	db.QueryRow("SELECT EXISTS(SELECT 1 FROM album WHERE id=?)", id).Scan(&exists)
	if exists {
		return errors.New("Album already exists")
	}
	return nil
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	var album album
	id := c.Param("id")
	row := db.QueryRow("SELECT * FROM album WHERE id=?", id)
	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.IndentedJSON(http.StatusOK, album)
}

func connectDB() {
	cfg := mysql.Config{
		//User:   os.Getenv("DBUSER"),
		//Passwd: os.Getenv("DBPASS"),
		User:   "root",
		Passwd: "mysql123456",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}

	var err error

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal("Unalbe to open mysql: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Unalbe to ping mysql: ", err)
	}
}
