package decoder

import (
	"fmt"

	"github.com/WithSoull/ChatServer/internal/converter/kafka"
	"github.com/WithSoull/ChatServer/internal/model"
	events_v1 "github.com/WithSoull/platform_common/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

type userCreatedDecoder struct{}

func NewUserCreatedDecoder() kafka.UserCreatedDecoder {
	return &userCreatedDecoder{}
}

func (d *userCreatedDecoder) Decode(data []byte) (model.UserCreatedEvent, error) {
	var pb events_v1.UserCreated
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.UserCreatedEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	createdAt := pb.GetCreatedAt().AsTime()
	return model.UserCreatedEvent{
		UserID:    pb.GetUserId(),
		CreatedAt: &createdAt,
	}, nil
}
