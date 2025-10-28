package converter

import (
	"github.com/WithSoull/ChatServer/internal/model"
	dmodel "github.com/WithSoull/ChatServer/internal/model" // domain/service model
	rmodel "github.com/WithSoull/ChatServer/internal/repository/model"
)

// ---------- Role ----------

func FromModelToRepoRole(r dmodel.Role) rmodel.Role {
	switch r {
	case dmodel.ROLE_ADMIN:
		return rmodel.ROLE_ADMIN
	case dmodel.ROLE_OWNER:
		return rmodel.ROLE_OWNER
	default:
		return rmodel.ROLE_USER
	}
}

func FromRepoToModelRole(r rmodel.Role) dmodel.Role {
	switch r {
	case rmodel.ROLE_ADMIN:
		return dmodel.ROLE_ADMIN
	case rmodel.ROLE_OWNER:
		return dmodel.ROLE_OWNER
	default:
		return dmodel.ROLE_USER
	}
}

// ---------- Chat ----------

func FromModelToRepoChat(c dmodel.Chat) rmodel.Chat {
	return rmodel.Chat{
		ID:          c.ChatID,
		OwnerID:     c.OwnerID,
		Name:        c.ChatInfo.Name,
		Description: c.ChatInfo.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func FromRepoToModelChat(r rmodel.Chat, userIDs []int64) dmodel.Chat {
	return dmodel.Chat{
		ChatID:  r.ID,
		OwnerID: r.OwnerID,
		ChatInfo: dmodel.ChatInfo{
			Name:        r.Name,
			Description: r.Description,
			UserIDs:     userIDs,
		},
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

// ---------- ChatUser ----------

func FromModelToRepoChatParticipant(m dmodel.ChatParticipant) rmodel.ChatParticipant {
	return rmodel.ChatParticipant{
		ChatID:    m.ChatID,
		UserID:    m.UserID,
		Role:      FromModelToRepoRole(m.Role),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func FromRepoToModelChatParticipant(r rmodel.ChatParticipant) dmodel.ChatParticipant {
	return dmodel.ChatParticipant{
		ChatID:    r.ChatID,
		UserID:    r.UserID,
		Role:      FromRepoToModelRole(r.Role),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

// ---------- Message ----------

func FromModelToRepoMessage(m dmodel.Message) rmodel.Message {
	return rmodel.Message{
		ID:        m.MessageID,
		ChatID:    m.ChatID,
		SenderID:  m.SenderID,
		Text:      m.Text,
		IsPinned:  m.IsPinned,
		SendAt:    m.SendAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func FromRepoToModelMessage(r rmodel.Message) dmodel.Message {
	return dmodel.Message{
		MessageID: r.ID,
		SenderID:  r.SenderID,
		ChatID:    r.ChatID,
		Text:      r.Text,
		IsPinned:  r.IsPinned,
		SendAt:    r.SendAt,
		UpdatedAt: r.UpdatedAt,
	}
}

// ---------- Messages ----------

func FromRepoToModelMessages(messages []rmodel.Message) []model.Message {
	msgs := make([]model.Message, len(messages))

	for i, message := range messages {
		msgs[i] = FromRepoToModelMessage(message)
	}

	return msgs
}

func FromModelToRepoMessages(messages []model.Message) []rmodel.Message {
	msgs := make([]rmodel.Message, len(messages))

	for i, message := range messages {
		msgs[i] = FromModelToRepoMessage(message)
	}

	return msgs
}
