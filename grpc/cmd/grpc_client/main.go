package main

import (
	"context"
	"log"
	"time"

	"github.com/fatih/color"
	desc "github.com/marktsarkov/auth/grpc/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	adress = "localhost:50052"
	userID = 1
	)

func main(){
	conn, err := grpc.Dial(adress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: #{err}")
	}

	defer conn.Close()

	c := desc.NewUserV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &desc.GetRequest{Id: userID})
	if err != nil {
		log.Fatalf("failder to get user by id: #{err}")
	}

	log.Printf(color.RedString("User info:\n"), color.GreenString("%v", r.GetUser()))
}