package grpcclient

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "lua-monologue-middleend/proto"

	"google.golang.org/grpc"
)

const serverAddress = "localhost:50051"

func SendChatMessage(message string, role string, user_id string) {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewChatServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.ChatRequest{Content: message, Role: role, Id: user_id}
	res, err := client.SendMessage(ctx, req)
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}

	fmt.Printf("Server replied: %s\n", res.Reply)
}
