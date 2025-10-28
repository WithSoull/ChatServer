package decoder

import (
	"fmt"

	"github.com/WithSoull/ChatServer/internal/converter/kafka"
	"github.com/WithSoull/ChatServer/internal/model"
	events_v1 "github.com/WithSoull/platform_common/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

type userDeletedDecoder struct{}

func NewUserDeletedDecoder() kafka.UserDeletedDecoder {
	return &userDeletedDecoder{}
}

func (d *userDeletedDecoder) Decode(data []byte) (model.UserDeletedEvent, error) {
	var pb events_v1.UserDeleted
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.UserDeletedEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	deletedAt := pb.GetDeletedAt().AsTime()
	return model.UserDeletedEvent{
		UserID:    pb.GetUserId(),
		DeletedAt: &deletedAt,
	}, nil
}
