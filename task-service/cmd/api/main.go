package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"task-service/data"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Models data.Models
}

func main() {
	db := connectDB()

	app := Config{
		Models: data.New(db),
	}
	

	// Running RPC
	rpcServer := NewRPCServer(app.Models)
	if err := rpc.Register(&rpcServer); err != nil{
		log.Panic("Failed to register RPC: ", err)
	}

	go rpcListen()

	log.Printf("Running task service on port %s (HTTP) and %s (RPC) ...\n", os.Getenv("HTTP_PORT"), os.Getenv("RPC_PORT"))

	server := http.Server{
		Addr: ":" + os.Getenv("HTTP_PORT"),
		Handler: app.routes(),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Panic("Failed to run server on port ", os.Getenv("HTTP_PORT"))
		return
	}
}

func rpcListen() error {
	listener, err := net.Listen("tcp", ":" + os.Getenv("RPC_PORT"))
	if err != nil {
		log.Panic("TCP Connection failed: ", err)
		return err
	}
	defer listener.Close()

	for {
		rpcConn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept request TCP: ", err)
			return err
		}

		go rpc.ServeConn(rpcConn)
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