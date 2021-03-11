package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/kelvynsantana/fc2-grpc/pb"
)

// UserService is a struct of user
type UserService struct {
	pb.UnimplementedUserServiceServer
}

// NewUserService return a UserService
func NewUserService() *UserService {
	return &UserService{}
}

//AddUser create a new User
func (*UserService) AddUser(ctx context.Context, req *pb.User) (*pb.User, error) {

	return &pb.User {
		Id: "1234",
		Name: req.GetName(),
		Email: req.GetEmail(),
	}, nil
}

// AddUserVerbose retutn detailed process to inser user in database
func (*UserService) AddUserVerbose(req *pb.User, stream pb.UserService_AddUserVerboseServer) error {
	stream.Send(&pb.UserResultStream{
		Status: "Init",
		User: &pb.User{},
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResultStream{
		Status: "Insert In Database",
		User: &pb.User{},
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResultStream{
		Status: "User has been inserted",
		User: &pb.User{
			Id: "123",
			Name: req.GetName(),
			Email: req.GetEmail(),
		},
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResultStream{
		Status: "Completed",
		User: &pb.User{
			Id: "123",
			Name: req.GetName(),
			Email: req.GetEmail(),
		},
	})

	time.Sleep(time.Second * 3)

	return nil
}

// AddUsers function
func (*UserService) AddUsers(stream pb.UserService_AddUsersServer) error {
	users := []*pb.User{}

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&pb.Users {
				User: users,
			})
		}

		if err != nil {
			log.Fatalf("Erro reciving stream: %v", err)
		}

		users = append(users, &pb.User{
			Id: req.GetId(),
			Name: req.GetName(),
			Email: req.GetEmail(),
		})

		fmt.Println("Adding", req.GetName())
	}
}

func (*UserService) AddUserStreamBoth(stream pb.UserService_AddUserStreamBothServer) error {
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil 
		}

		if err != nil {
			log.Fatalf("Error receiving stream from the client: %v", err)
		}

		err = stream.Send(&pb.UserResultStream{
			Status: "User Added",
			User: req,
		})
	}
}