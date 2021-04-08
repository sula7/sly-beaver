package main

import (
	"log"
	"os"
	"path/filepath"

	"sly-beaver/storage"
)

func main() {
	workDir := filepath.Dir(os.Args[0])
	dbFilePath := filepath.Join(filepath.Join(workDir, "sly-beaver.db"))

	db, err := storage.OpenDB(dbFilePath)
	if err != nil {
		log.Fatalln("open db connection:", err)
	}

	defer func() {
		err = db.Close()
		if err != nil {
			log.Println("close db connection:", err)
		}
	}()

	err = db.RunMigrations()
	if err != nil {
		log.Fatalln("run migrations:", err)
	}
}
