module transaction-service

go 1.23.2

require github.com/lib/pq v1.10.9

require github.com/rabbitmq/amqp091-go v1.10.0

require (
	google.golang.org/grpc v1.73.0
	proto v0.0.0
)

require (
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace proto => ../proto
