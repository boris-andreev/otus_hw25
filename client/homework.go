package main

import (
	"context"
	"log"
	"proto_api/pkg/grpc/v1/todo_api"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func testHomeworkApi(conn *grpc.ClientConn) {
	log.Println("Homework items")
	homeworkClient := todo_api.NewHomeworkServiceClient(conn)

	log.Println("Crreate item")
	_, err := homeworkClient.CreateHomework(
		context.Background(),
		&todo_api.CreateHomeworkRequest{
			Description: "Description",
		},
	)

	if err != nil {
		log.Println("Error: ", err.Error())
	}

	log.Println("Item created")

	lr, err := homeworkClient.ListHomework(
		context.Background(),
		&emptypb.Empty{},
	)

	if err != nil {
		log.Fatal("Error: ", err.Error())
	}

	log.Println("List Items:")

	for _, homeworkItem := range lr.Result {
		log.Printf("Id: %s, Description: %s", homeworkItem.Id, homeworkItem.Description)
	}

	id := lr.Result[0].Id

	ir, err := homeworkClient.GetHomework(
		context.Background(),
		&todo_api.GetHomeworkRequest{
			Id: id,
		})

	if err != nil {
		log.Fatal("Error: ", err.Error())
	}

	log.Println("Item:")
	log.Printf("Id: %s, Description: %s", ir.Id, ir.Description)

	log.Println("Delete Item")
	_, err = homeworkClient.DeleteHomework(
		context.Background(),
		&todo_api.DeleteHomeworkRequest{
			Id: id,
		})

	if err != nil {
		log.Fatal("Error: ", err.Error())
	}
	log.Println("Deleted")
}
