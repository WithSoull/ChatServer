package main

import (
	"context"
	"fmt"
	"log"
	"net"

	desc "github.com/WithSoull/ChatServer/pkg/chat/v1"
	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)
const (
  grpcPort = 50050
)


type server struct {
  desc.UnimplementedChatV1Server
}


func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
  chatId := gofakeit.Int64()
  log.Printf("Chat(%d) creating with this users:", chatId)
  for _, user := range req.GetUsernames() {
    log.Printf("User: %s", user)
  }
  
  // TODO: Creating Chat and adding users

  return &desc.CreateResponse{
    Id: chatId,
  }, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
  log.Printf("Chat deleting {Id: %d}",
    req.GetId(),
  )

  // TODO: Chat delete

  return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
  log.Printf("Message sent from %s with text: \"%s\" at %s",
    req.GetFrom(),
    req.GetText(),
    req.GetSentAt(),
  )

  // TODO: Send message

  return &emptypb.Empty{}, nil
}
func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
