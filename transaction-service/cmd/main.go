package main

import (
	"database/sql"
	"log"
	"net"
	"net/rpc"
	"os"
	"time"
	"transaction-service/data"

	_ "github.com/lib/pq"
)

type Config struct{
	Models data.Models
}

func main(){
	db := connectDB()

	app := Config{
		Models: data.New(db),
	}

	rpcServer := NewTransactionRPCServer(app.Models)

	if err := rpc.Register(rpcServer); err != nil {
		log.Panic("Failed to register RPC object: ", err)
	}
	
	listenRPC()
}

func listenRPC() error {
	port := os.Getenv("RPC_PORT")
	log.Printf("Starting rpc server on port %s ...\n", port)

	listener, err := net.Listen("tcp", os.Getenv("RPC_HOST") + ":" + port)
	if err != nil {
		log.Fatal("Failed to listen RPC: ", err)
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