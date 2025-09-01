package converter

import (
	"github.com/WithSoull/ChatServer/internal/model"
	desc "github.com/WithSoull/ChatServer/pkg/chat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromProtoToModelRole(r desc.Role) model.Role {
	switch r {
	case desc.Role_ROLE_ADMIN:
		return model.ROLE_ADMIN
	case desc.Role_ROLE_OWNER:
		return model.ROLE_OWNER
	default:
		return model.ROLE_USER
	}
}

func FromModelToProtoRole(r model.Role) desc.Role {
	switch r {
	case model.ROLE_ADMIN:
		return desc.Role_ROLE_ADMIN
	case model.ROLE_OWNER:
		return desc.Role_ROLE_OWNER
	default:
		return desc.Role_ROLE_USER
	}
}

// ChatInfo
func FromProtoToModelChatInfo(proto *desc.ChatInfo) model.ChatInfo {
	if proto == nil {
		return model.ChatInfo{}
	}
	return model.ChatInfo{
		Name:        proto.GetName(),
		Description: proto.GetDescription(),
		UserIDs:     proto.GetUserIds(),
	}
}

func FromModelToProtoChatInfo(model model.ChatInfo) *desc.ChatInfo {
	return &desc.ChatInfo{
		Name:        model.Name,
		Description: model.Description,
		UserIds:     model.UserIDs,
	}
}

// Chat
func FromProtoToModelChat(proto *desc.Chat) model.Chat {
	if proto == nil {
		return model.Chat{}
	}
	return model.Chat{
		ChatID:   proto.GetChatId(),
		OwnerID:  proto.GetOwnerId(),
		ChatInfo: FromProtoToModelChatInfo(proto.GetChatInfo()),
	}
}

func FromModelToProtoChat(model model.Chat) *desc.Chat {
	return &desc.Chat{
		ChatId:   model.ChatID,
		OwnerId:  model.OwnerID,
		ChatInfo: FromModelToProtoChatInfo(model.ChatInfo),
	}
}

// Message
func FromProtoToModelMessage(proto *desc.Message) model.Message {
	if proto == nil {
		return model.Message{}
	}

	out := model.Message{
		MessageID: proto.GetMessageId(),
		SenderID:  proto.GetSenderId(),
		ChatID:    proto.GetChatId(),
		Text:      proto.GetText(),
		IsPinned:  proto.GetIsPinned(),
		SendAt:    proto.GetSendAt().AsTime(),
		UpdatedAt: proto.GetUpdatedAt().AsTime(),
	}

	return out
}

func FromModelToProtoMessage(model model.Message) *desc.Message {
	return &desc.Message{
		MessageId: model.MessageID,
		SenderId:  model.SenderID,
		ChatId:    model.ChatID,
		Text:      model.Text,
		IsPinned:  model.IsPinned,
		SendAt:    timestamppb.New(model.SendAt),
		UpdatedAt: timestamppb.New(model.UpdatedAt),
	}
}

// Messages (slice of model.Message)
func FromProtoToModelMessages(messages []*desc.Message) []model.Message {
	out := make([]model.Message, len(messages))

	for _, message := range messages {
		out = append(out, FromProtoToModelMessage(message))
	}

	return out
}

func FromModelToProtoMessages(messages []model.Message) []*desc.Message {
	out := make([]*desc.Message, len(messages))

	for _, message := range messages {
		out = append(out, FromModelToProtoMessage(message))
	}

	return out
}
