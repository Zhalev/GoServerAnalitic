package main

import (
	"database/sql"
	"fmt"
	_ "net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123"
	dbname   = "devices_management"
)

type File struct {
	FileContent []byte
	FileName    string
	FilePath    string
	FileMD5     string
	FileLoaded  time.Time
}

func main() {

	download(conn())

	return

	/* пробуем на одном файле */
	/*	var (
			filename string
			name     []byte
		)
		row := db.QueryRow("select file_name, file_content from updater where id=$1", 2)
		if err != nil {
			log.Fatal(err)
		}
		//defer rows.Close()
		//for rows.Next() {
		err = row.Scan(&filename, &name)
		if err != nil {
			log.Fatal(err)
		}
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("Unable to create file:", err)
			os.Exit(1)
		}
		defer file.Close()
		file.Write(name)

		fmt.Println("Done.")
		//log.Println( name)
		//}
		err = row.Err()
		if err != nil {
			log.Fatal(err)
		}*/
}

func filenames(db *sql.DB) []File {
	files := []File{}
	rows, err := db.Query("SELECT file_name, file_path, md5, file_content, time_loaded FROM updater ")
	__err_panic(err)
	for rows.Next() {
		post := File{}
		err = rows.Scan(&post.FileName, &post.FilePath, &post.FileMD5, &post.FileContent, &post.FileLoaded)
		__err_panic(err)
		files = append(files, post)
	}
	return files
}

func __err_panic(err error) {
	if err != nil {
		panic(err)
	}
}
func conn() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	__err_panic(err)
	//	defer db.Close()

	err = db.Ping()
	__err_panic(err)

	fmt.Println("Successfully connected!")
	return db
}
func download(db *sql.DB) {
	//fmt.Println(filenames(db))
	files := filenames(db)
	for _, file := range files {
		err := os.MkdirAll(file.FilePath, 0777)
		__err_panic(err)
		f, err := os.Create(file.FilePath + file.FileName)
		if err != nil {
			fmt.Println("Unable to create file:", err)
			os.Exit(1)
		}
		f.Write(file.FileContent)
		f.Close()
	}
	fmt.Println("Done.")

}
