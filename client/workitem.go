package main

import (
	"context"
	"log"
	"proto_api/pkg/grpc/v1/todo_api"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func testWorkoutkApi(conn *grpc.ClientConn) {
	log.Println("Workout items")
	workoutClient := todo_api.NewWorkoutServiceClient(conn)

	log.Println("Crreate item")
	_, err := workoutClient.CreateWorkout(
		context.Background(),
		&todo_api.CreateWorkoutRequest{
			Target: "Target",
		},
	)

	if err != nil {
		log.Println("Error: ", err.Error())
	}

	log.Println("Item created")

	lr, err := workoutClient.ListWorkout(
		context.Background(),
		&emptypb.Empty{},
	)

	if err != nil {
		log.Fatal("Error: ", err.Error())
	}

	log.Println("List Items:")

	for _, item := range lr.Result {
		log.Printf("Id: %s, Target: %s", item.Id, item.Target)
	}

	id := lr.Result[0].Id

	ir, err := workoutClient.GetWorkout(
		context.Background(),
		&todo_api.GetWorkoutRequest{
			Id: id,
		})

	if err != nil {
		log.Fatal("Error: ", err.Error())
	}

	log.Println("Item:")
	log.Printf("Id: %s, Target: %s", ir.Id, ir.Target)

	log.Println("Delete Item")
	_, err = workoutClient.DeleteWorkout(
		context.Background(),
		&todo_api.DeleteWorkoutRequest{
			Id: id,
		})

	if err != nil {
		log.Fatal("Error: ", err.Error())
	}
	log.Println("Deleted")
}
