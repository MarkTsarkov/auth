package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/brianvoe/gofakeit"
	desc "github.com/marktsarkov/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50052

type server struct {
	desc.UnimplementedUserV1Server
}

func (s* server) Get (ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("User id: #{req.GetId()}")

	return &desc.GetResponse{
		User: &desc.User{
			Id: req.GetId(),
			Info: &desc.UserInfo{
				Name: gofakeit.BeerName(),
				Email: gofakeit.Email(),
				Role: 1,
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{})

	log.Printf("server lisening at #{lis.Addr()}")

	if err = s.Serve(lis); err != nil {
		log.Fatal("failed to serve: #{err}")
	}
}