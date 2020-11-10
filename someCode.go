package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	_ "net/http"
	"os"
	"os/exec"
	"path/filepath"
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
	FileName    string
	FilePath    string
	FileMD5     string
	FileContent []byte
	FileLoaded  time.Time
}

func main() {
	db := connect()
	files := getFilesFromDB(db)
	download(files)
	addFileInDB(db, "e:/GoLangProjects/awesomeProject/PlanLogController.go")
	// запуск внешней программы
	/*
		path, err := exec.LookPath("updater.exe")
		if err != nil {
			fmt.Println("Файл не найден")
		}
		cmd := exec.Command(path)
	*/
	cmd := exec.Command("updater.exe")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Unable to Run Updater file:", err)
		os.Exit(1)
	}
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

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func getFilesFromDB(db *sql.DB) []File {
	files := []File{}
	rows, err := db.Query("SELECT file_name, file_path, md5, file_content, time_loaded FROM updater r WHERE time_loaded = (select max(time_loaded) from updater u where r.file_name = u.file_name ) order by r.file_name")
	errPanic(err)
	defer rows.Close()
	for rows.Next() {
		post := File{}
		err = rows.Scan(&post.FileName, &post.FilePath, &post.FileMD5, &post.FileContent, &post.FileLoaded)
		errPanic(err)
		files = append(files, post)
	}
	return files
}

func addFileInDB(db *sql.DB, path string) {
	// выход если файла нет
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return
	}
	file := File{}
	file.FileMD5 = FileMD5(path)
	file.FileName = filepath.Base(path)

	// путь к исполняемому файлу
	cwd, err := os.Executable()
	errPanic(err)
	rootDir := filepath.Dir(cwd)
	cwd, err = filepath.Rel(rootDir, path)
	errPanic(err)
	file.FilePath = filepath.Dir(cwd)

	//fmt.Println(file.FilePath)

	// вытащить из пути название файла и папок
	/*
		f, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer f.Close()

		// чтение данных из файла (ест пробелы)

			wr := bytes.Buffer{}
			sc := bufio.NewScanner(f)
			for sc.Scan() {
				wr.WriteString(sc.Text())
			}
			fmt.Println(wr.String())
	*/

	// можно и так тож норм
	/*	data := make([]byte, 64)
		var str string
			for{
				n, err := f.Read(data)
				if err == io.EOF{   // если конец файла
					break           // выходим из цикла
				}
				str += string(data[:n])

			}
			fmt.Print(str)

		file.FileContent = []byte(str)*/

	fContent, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	file.FileContent = fContent

	file.FileLoaded = time.Now()

	defer db.Close()
	res, err := db.Exec("INSERT INTO updater (file_name, file_path, md5, file_content, time_loaded) "+
		"VALUES ($1, $2, $3, $4, $5)", file.FileName, file.FilePath, file.FileMD5, file.FileContent, file.FileLoaded)
	errPanic(err)
	if n, _ := res.RowsAffected(); n > 0 {
		fmt.Println("Файл успешно добавлен в базу")

	}
}

func connect() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	errPanic(err)
	//	defer db.Close()

	err = db.Ping()
	errPanic(err)

	fmt.Println("Successfully connected!")
	return db
}

func download(files []File) {
	//fmt.Println(filenames(db))
	var (
		h  string
		h2 string
	)
	for _, file := range files {
		// если папок нет, то создаем
		if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
			err := os.MkdirAll(file.FilePath, 0777)
			errPanic(err)
		}

		path := file.FilePath + "\\" + file.FileName
		h2 = FileMD5(path)
		if h = fmt.Sprintf("%x", md5.Sum(file.FileContent)); h != h2 {
			WriteContentToFile(path, file.FileContent)
			fmt.Println(file.FileName, "был перезаписан", h, h2)
		}
	}
	fmt.Println("Done.")
}

func FileMD5(path string) string {
	h := md5.New()
	f, err := os.OpenFile(path, os.O_CREATE, 0666)
	errPanic(err)
	defer f.Close()
	_, err = io.Copy(h, f)
	errPanic(err)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func WriteContentToFile(path string, content []byte) {
	f, err := os.OpenFile(path, os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer f.Close()
	_, err = f.Write(content)
	errPanic(err)

}
