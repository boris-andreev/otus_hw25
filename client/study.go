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
	homeworkClient := todo_api.NewStudyServiceClient(conn)

	log.Println("Crreate item")
	_, err := homeworkClient.CreateStudy(
		context.Background(),
		&todo_api.CreateStudyRequest{
			Topic: "Math",
		},
	)

	if err != nil {
		log.Println("Error: ", err.Error())
	}

	log.Println("Item created")

	lr, err := homeworkClient.ListStudy(
		context.Background(),
		&emptypb.Empty{},
	)

	if err != nil {
		log.Fatal("Error: ", err.Error())
	}

	log.Println("List Items:")

	for _, item := range lr.Result {
		log.Printf("Id: %d, Topic: %s", item.Id, item.Topic)
	}

	id := lr.Result[0].Id

	ir, err := homeworkClient.GetStudy(
		context.Background(),
		&todo_api.GetStudyRequest{
			Id: id,
		})

	if err != nil {
		log.Fatal("Error: ", err.Error())
	}

	log.Println("Item:")
	log.Printf("Id: %d, Topic: %s", ir.Id, ir.Topic)

	log.Println("Delete Item")
	_, err = homeworkClient.DeleteStudy(
		context.Background(),
		&todo_api.DeleteStudyRequest{
			Id: id,
		})

	if err != nil {
		log.Fatal("Error: ", err.Error())
	}
	log.Println("Deleted")
}
