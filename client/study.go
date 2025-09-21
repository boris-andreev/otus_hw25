package main

import (
	"context"
	"log"
	"proto_api/pkg/grpc/v1/todo_api"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func testStudykApi(conn *grpc.ClientConn) {
	log.Println("Study items")
	studyClient := todo_api.NewStudyServiceClient(conn)

	log.Println("Crreate item")
	_, err := studyClient.CreateStudy(
		context.Background(),
		&todo_api.CreateStudyRequest{
			Topic: "Math",
		},
	)

	if err != nil {
		log.Println("Error: ", err.Error())
	}

	log.Println("Item created")

	lr, err := studyClient.ListStudy(
		context.Background(),
		&emptypb.Empty{},
	)

	if err != nil {
		log.Fatal("Error: ", err.Error())
	}

	log.Println("List Items:")

	for _, item := range lr.Result {
		log.Printf("Id: %s, Topic: %s", item.Id, item.Topic)
	}

	id := lr.Result[0].Id

	ir, err := studyClient.GetStudy(
		context.Background(),
		&todo_api.GetStudyRequest{
			Id: id,
		})

	if err != nil {
		log.Fatal("Error: ", err.Error())
	}

	log.Println("Item:")
	log.Printf("Id: %s, Topic: %s", ir.Id, ir.Topic)

	log.Println("Delete Item")
	_, err = studyClient.DeleteStudy(
		context.Background(),
		&todo_api.DeleteStudyRequest{
			Id: id,
		})

	if err != nil {
		log.Fatal("Error: ", err.Error())
	}
	log.Println("Deleted")
}
