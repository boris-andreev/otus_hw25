package main

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	testHomeworkApi(conn)
	log.Println("!!!!!!--------------------------!!!!!!")
	testStudykApi(conn)
	log.Println("!!!!!!--------------------------!!!!!!")
	testWorkoutkApi(conn)
}
