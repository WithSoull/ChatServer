package queries

const (
	PlaceHolder = "$"
)

const (
	InsertNewChat = `INSERT INTO chats (created_at) 
			VALUES (NOW())
			RETURNING id;`
	InsertNewParticipant = `INSERT INTO chat_participants (chat_id, username, joined_at)
			VALUES ($1, $2, CURRENT_TIMESTAMP)`
)
