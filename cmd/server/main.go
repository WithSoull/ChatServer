package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"strings"

	"github.com/WithSoull/ChatServer/internal/config"
	"github.com/WithSoull/ChatServer/internal/config/env"
	"github.com/WithSoull/ChatServer/internal/queries"
	desc "github.com/WithSoull/ChatServer/pkg/chat/v1"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)


var configPath string

func init() {
  flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
  desc.UnimplementedChatV1Server
	db *pgxpool.Pool
}

func NewServer(dbPool *pgxpool.Pool) *server {
	return &server{
		db: dbPool,
	}
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	if len(req.GetUsernames()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "no users in chat")
	}

	// Begin transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
    log.Printf("failed to begin transaction: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Insert new chat and users
	var chatId int64
	err = tx.QueryRow(ctx, queries.InsertNewChat).Scan(&chatId)
	if err != nil {
		log.Printf("failed to create new chat with error: %v", err)
		return nil, status.Error(codes.Internal, "failed to create new chat")
	}
	
  for _, username := range req.GetUsernames() {
		_, err = tx.Exec(ctx, queries.InsertNewParticipant, chatId, username)
		if err != nil {
			log.Printf("failed to create add participant %s with error: %v", username, err)
			return nil, status.Error(codes.Internal, "failed to create new chat")
		}
  }

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
    log.Printf("failed to commit transaction: %v", err)
    return nil, status.Errorf(codes.Internal, "failed to commit transaction")
	}

  return &desc.CreateResponse{
    Id: chatId,
  }, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	// Begin transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
    log.Printf("failed to begin transaction: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Deleting chat
	_, err = tx.Exec(ctx, queries.DeleteChatById, req.GetId())
	if err != nil {
    log.Printf("failed to delete chat: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete chat")
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
    log.Printf("failed to commit transaction: %v", err)
    return nil, status.Errorf(codes.Internal, "failed to commit transaction")
	}

  return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	// Begin transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
    log.Printf("failed to begin transaction: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Insert message
	_, err = tx.Exec(ctx, queries.InsertNewMessage, req.GetChatId(), req.GetFrom(), req.GetText())
	if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
					switch pgErr.Code {
					case "23503": // Foreign key violation
							if strings.Contains(pgErr.Message, "messages_chat_id_fkey") {
									return nil, status.Error(codes.NotFound, "chat not found")
							}
					}
			}
    
    log.Printf("failed to send message: %v", err)
    return nil, status.Errorf(codes.Internal, "failed to send message: %v", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
    log.Printf("failed to commit transaction: %v", err)
    return nil, status.Errorf(codes.Internal, "failed to commit transaction")
	}
  return &emptypb.Empty{}, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	//Load config
	if err := config.Load(configPath); err != nil {
    log.Printf("configPath=%s", configPath)
    log.Fatalf("failed load config: %s", err)
	}

  grpcConfig, err := env.NewGRPCConfig()
  if err != nil {
    log.Fatalf("failed load grpc config: %s", err)
  }

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed load pg congig: %s", err)
	}

	// Create connection pool
	dbPool, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to create connection pool: %s", err)
	}
	defer dbPool.Close()

	// Ping db
  if err := dbPool.Ping(ctx); err != nil {
    log.Fatalf("failed to connect to database: %v", err)
  }

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	desc.RegisterChatV1Server(s, NewServer(dbPool))

	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
