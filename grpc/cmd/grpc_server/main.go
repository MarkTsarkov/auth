package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/marktsarkov/auth/grpc/cmd/internal/config/env"
	desc "github.com/marktsarkov/auth/grpc/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	"honnef.co/go/tools/config"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedUserV1Server
	pool *pgxpool.Pool
}

func (s* server) Create (ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("User id: %v", req.Info)

	return &desc.CreateResponse{
		Id: req.Info.Id,
	}, nil
}

func (s* server) Get (ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("User id: %v", req.Id)

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

func (s* server) Update (ctx context.Context, req *desc.UpdateRequest) (*desc.UpdateResponse, error) {
	log.Printf("User id: %v", req.Id)

	return &desc.UpdateResponse{
	}, nil
}

func (s* server) Delete (ctx context.Context, req *desc.DeleteRequest) (*desc.DeleteResponse, error) {
	log.Printf("User id: %v", req.Id)

	return &desc.DeleteResponse{
	}, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	_, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{pool: pool})

	log.Printf("server lisening at #{lis.Addr()}")

	if err = s.Serve(lis); err != nil {
		log.Fatal("failed to serve: #{err}")
	}
}