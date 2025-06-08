package queries

const (
	PlaceHolder = "$"
)

const (
	InsertNewChat = `INSERT INTO chats (created_at) 
			VALUES (NOW())
			RETURNING id;`
	InsertNewParticipant = `INSERT INTO chat_participants (chat_id, username, joined_at)
			VALUES ($1, $2, NOW())`
	DeleteChatById = `DELETE FROM chats WHERE id = $1`
	InsertNewMessage = `INSERT INTO messages (chat_id, from_user, text, timestamp)
			VALUES ($1, $2, $3, NOW())`
)
