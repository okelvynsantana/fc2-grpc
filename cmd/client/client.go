package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/kelvynsantana/fc2-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not conncet to GRPC server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	AddUsersStreamBoth(client)

}

// AddUser create a new User
func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "123456",
		Name:  "Kelvyn",
		Email: "kelvyn@kelvyn.com",
	}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make GRPC request: %v", err)
	}

	fmt.Println(res)
}

// AddUserVerbose recieve realtime status in User creation
func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "2",
		Name:  "Kelvyn Santana",
		Email: "santanakelvyn@gmail.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not  make Grpc request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not recieve the message: %v", err)
		}
		fmt.Println("Status:", stream.Status)
	}
}

// AddUsers func
func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "1",
			Name:  "Kelvyn 1",
			Email: "email1@teste.com",
		},
		&pb.User{
			Id:    "2",
			Name:  "Kelvyn 2",
			Email: "email2@teste.com",
		},
		&pb.User{
			Id:    "3",
			Name:  "Kelvyn 3",
			Email: "email3@teste.com",
		},
		&pb.User{
			Id:    "4",
			Name:  "Kelvyn 4",
			Email: "email4@teste.com",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error receiving request: %v", err)
	}
	fmt.Println(res)

}

// AddUsersStreamBoth func
func AddUsersStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUserStreamBoth(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "1",
			Name:  "Kelvyn 1",
			Email: "email1@teste.com",
		},
		&pb.User{
			Id:    "2",
			Name:  "Kelvyn 2",
			Email: "email2@teste.com",
		},
		&pb.User{
			Id:    "3",
			Name:  "Kelvyn 3",
			Email: "email3@teste.com",
		},
		&pb.User{
			Id:    "4",
			Name:  "Kelvyn 4",
			Email: "email4@teste.com",
		},
	}

	wait := make(chan int)
	go func() {
		for _, req := range reqs {
			fmt.Println("Sending User: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error recieving data: %v", err)
			}

			fmt.Printf("recebendo User %v com status %v", res.GetUser(), res.GetStatus())
		}
		close(wait)
	}()
	<-wait
}
