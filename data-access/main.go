package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	// Capture connection properties.
	cfg := mysql.Config{
		//User:   os.Getenv("DBUSER"),
		//Passwd: os.Getenv("DBPASS"),
		User:   "root",
		Passwd: "mysql123456",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}

	test_string := "connect"

	fmt.Printf("trying to %q\n", test_string)

	// Get a database handle.
	//println(os.Getenv("DBUSER"), ",", os.Getenv("DBPASS"))
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal("unable to open mysql", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal("connection time out\n", pingErr)
	}
	fmt.Println("Connected!")

	var toInsert string
	fmt.Println("Do you want to insert a new row into the table alubm? (y/n)")
	_, err = fmt.Scanln(&toInsert)
	if err != nil {
		log.Fatal(err)
	}
	for {
		if !notValid(toInsert) {
			break
		}
		fmt.Println("Invalid input. Try again.")
		_, err = fmt.Scanln(&toInsert)
		if err != nil {
			log.Fatal(err)
		}
	}
	if toInsert == "y" || toInsert == "Y" {
		albID, err := addAlbum(Album{
			Title:  "The Modern Sound of Betty Carter",
			Artist: "Betty Carter",
			Price:  49.99,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("ID of added album: %v\n:", albID)
	}

	var albums []Album
	var albumsKey string

	fmt.Println("please input a regular expression that you wanna use to find albums")
	_, err = fmt.Scanln(&albumsKey)
	if err != nil {
		log.Fatal(err)
	}

	albums, err = albumsByArtist(albumsKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found : %v\n", albums)

	var albumKey int64
	fmt.Println("please input a integer that you wanna use to find a single album")
	_, err = fmt.Scanln(&albumKey)

	album, err := albumsByID(albumKey)
	if err != nil {
		log.Fatal("unable to find the album with id 2\n", err)
	}
	fmt.Printf("Album found : %v\n", album)
}

func albumsByArtist(name string) ([]Album, error) {
	var albums []Album
	rows, err := db.Query("SELECT * FROM album WHERE artist REGEXP ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: query : %v", name, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("albumsByArtist %q: close rows: %v", name, err)
		}
	}()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: loop through the rows: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

func albumsByID(id int64) (Album, error) {
	var album Album
	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		if err == sql.ErrNoRows {
			return album, fmt.Errorf("albumsByID %d: no such album", id)
		}
		return album, fmt.Errorf("albumsByID %q: %v", id, err)
	}
	return album, nil
}

// addAlbum adds the specified album to the database,
// returning the album ID of the new entry

func addAlbum(alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album(title, artist, price) VALUES(?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}

func notValid(toInsert string) bool {
	if toInsert == "y" || toInsert == "Y" || toInsert == "n" || toInsert == "N" {
		return false
	}
	return true
}
