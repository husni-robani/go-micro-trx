package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"
	"transaction-service/data"

	_ "github.com/lib/pq"
)

type Config struct{
	Models data.Models
}

const webPort = "80"

func main(){
	db := connectDB()
	db.Query("SELECT * FROM transactions;")

	app := Config{
		Models: data.New(db),
	}

	log.Printf("Starting web server on port %s ...\n", webPort)

	srv := http.Server{
		Handler: app.routes(),
		Addr: ":" + webPort,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic("Failed to run web server")
	}

}

func connectDB() *sql.DB {
	dsn := os.Getenv("DSN")
	log.Println("Starting Database Connection ...")

	db, err := sql.Open("postgres" , dsn)
	
	count := 1

	for {
		if err != nil {
			log.Println("Database connection not ready yet ...")
			count++
		}else {
			if err := db.Ping(); err != nil {
				log.Println("Database connection not ready yet ...")
				count++
			}else {
				log.Println("Database Connected!")
				return db
			}
		}

		if count >= 5 {
			log.Panic("Connection Databasae Failed: ", err)
			return nil
		}

		time.Sleep(3 * time.Second)
	}
}