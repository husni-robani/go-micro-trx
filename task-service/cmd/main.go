package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"task-service/cmd/api"
	"task-service/cmd/config"
	"task-service/cmd/grpc_server"

	// task_grpc_server "task-service/cmd/grpc_server"
	"task-service/cmd/rpc_server"
	"task-service/data"
	"time"

	taskpb "proto/task"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	models := data.New(connectDB())
	rpcClientTransaction := connectRPCClient(os.Getenv("TRANSACTION_RPC_ADDRESS"))
	
	app := config.Config{
		Models: models,
	}

				// ============RPC============
	rpc_server := rpc_server.RPCServer{
		App: &app,
		RPCClientTransaction: rpcClientTransaction,
	}
	// register RPC object
	if err := rpc.Register(rpc_server); err != nil {
		log.Fatal("Failed to register RPC: ", err)
	}
	// serve RPC
	go serveRPC()

	
				// ============gRPC============
	g_server := grpc.NewServer()
	
	// register gRPC object
	taskpb.RegisterTaskServiceServer(g_server, grpc_server.TaskGRPCServer{
		Models: models,
		RPCClientTransaction: rpcClientTransaction,
	})

	// serve gRPC
	go serveGRPC(g_server)

				// ============API============
	api := api.APIHandler{
		App: &app,
	}
	serveAPI(api.Routes())
}

func serveGRPC(grpcServer *grpc.Server) {
	port := os.Getenv("GRPC_PORT")
	counter := 0
	log.Printf("Starting gRPC on Port:%s ....", port)

	var err error
	for {
		if counter >= 4 {
				log.Fatal("gRPC connection failed: ", err)
		}
		
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
		if err != nil {
			counter ++
			log.Println("gRPC not connected yet....")
			time.Sleep(time.Second * 2)
			continue
		}
		if err := grpcServer.Serve(lis); err != nil {
			counter ++
			log.Println("gRPC not connected yet....")
			time.Sleep(time.Second * 2)
			continue
		}
		
	}
}

func serveRPC() error {
	port := os.Getenv("RPC_PORT")

	log.Printf("Starting RPC on port %v....\n", port)

	listener, err := net.Listen("tcp", ":" + port)
	if err != nil {
		log.Panic("RPC | TCP Connection failed: ", err)
		return err
	}
	defer listener.Close()

	for {
		rpcConn, err := listener.Accept()
		if err != nil {
			log.Println("RPC | Failed to accept request TCP: ", err)
			return err
		}

		go rpc.ServeConn(rpcConn)
	}
}

func serveAPI(handlers http.Handler) {
	port := os.Getenv("HTTP_PORT")

	log.Printf("Starting HTTP on port %v....", port)

	server := http.Server{
		Addr: ":" + port,
		Handler: handlers,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("HTTP | Failed to run server on port ", port)
		return
	}
}

func connectDB() *sql.DB {
	dsn := os.Getenv("POSTGRE_DSN")
	counts := 1
	
	for {
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Println("Postgre connection not ready yet ...")
			counts++
		}else {
			err = db.Ping()
			if err == nil {
				log.Println("Database Connected!")
				return db
			}

			log.Println("Postgre connection not ready yet ...")
			counts++
		}

		time.Sleep(3 * time.Second)

		if counts > 5 {
			log.Panic("Cant connect to Postgre database: ", err)
			return nil
		}
	}
}

func connectRPCClient(address string) *rpc.Client {
	counter := 0
	
	for {
		client, err := rpc.Dial("tcp", address)

		if counter >= 5 {
			log.Fatalf("%s connection failed: %v", address, err)
			return nil
		}

		if err != nil {
			log.Printf("%s not connected yet....", address)
			time.Sleep(time.Second * 3)
			counter++
			continue
		}

		log.Printf("%v connected!", address)

		return client
	}
}