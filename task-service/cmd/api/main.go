package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"task-service/data"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Models data.Models
}

const webPort = "80"

func main() {
	db := connectDB()

	app := Config{
		Models: data.New(db),
	}
	
	log.Printf("Running task service on port %s ...\n", webPort)

	server := http.Server{
		Addr: ":" + webPort,
		Handler: app.routes(),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Panic("Failed to run server on port ", webPort)
		return
	}
}

func connectDB() *sql.DB {
	dsn := os.Getenv("DSN")
	counts := 1
	
	for {
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Println("Connection not ready yet ...")
			counts++
		}else {
			err = db.Ping()
			if err == nil {
				log.Println("Database Connected!")
				return db
			}

			log.Println("Connection not ready yet ...")
			counts++
		}

		time.Sleep(3 * time.Second)

		if counts > 5 {
			log.Panic("Cant connect to database: ", err)
			return nil
		}
	}
}