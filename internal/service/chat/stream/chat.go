package stream

import (
	"context"
	"sync"

	"github.com/WithSoull/ChatServer/internal/model"
	"github.com/WithSoull/platform_common/pkg/logger"
	"go.uber.org/zap"
)

// general sheme
// chatID -> userID -> message stream

type chatStream struct {
	// userID -> message stream
	stream map[int64]chan *model.Message
	mu     sync.RWMutex

	bufferLength int64
}

func newChatStream(bufferLength int64) *chatStream {
	return &chatStream{
		stream:       make(map[int64]chan *model.Message),
		bufferLength: bufferLength,
	}
}

// data race safly
func (cs *chatStream) msgStream(userID int64) chan *model.Message {
	cs.mu.Lock()
	msgStream, ok := cs.stream[userID]
	if !ok {
		msgStream = make(chan *model.Message, cs.bufferLength)
		cs.stream[userID] = msgStream
	}
	cs.mu.Unlock()

	return msgStream
}

// data race safly
func (cs *chatStream) addMsg(msg *model.Message) {
	// Make copy
	cs.mu.RLock()
	channels := make([]chan *model.Message, 0, len(cs.stream))
	for _, msgStream := range cs.stream {
		channels = append(channels, msgStream)
	}
	cs.mu.RUnlock()

	for _, ch := range channels {
		select {
		case ch <- msg:
		default:
		}
	}
}

// data race safly
func (cs *chatStream) removeMsgStream(userID int64) {
	cs.mu.Lock()
	// this action will be a signal to grpc/http handler, that
	// stream is closed
	cs.stream[userID] <- nil
	close(cs.stream[userID])
	delete(cs.stream, userID)
	cs.mu.Unlock()
}

type ChatStreams struct {
	// chatID -> ChatStreams
	streams map[int64]*chatStream
	mu      sync.RWMutex

	bufferLength int64
}

func NewChatStreams(bufferLength int64) *ChatStreams {
	return &ChatStreams{
		streams:      make(map[int64]*chatStream),
		bufferLength: bufferLength,
	}
}

// data race safly
func (cs *ChatStreams) setChatStream(chatID int64) *chatStream {
	cs.mu.Lock()
	chatStream, ok := cs.streams[chatID]
	if !ok {
		chatStream = newChatStream(cs.bufferLength)
		cs.streams[chatID] = chatStream
	}
	cs.mu.Unlock()

	return chatStream
}

func (cs *ChatStreams) getChatStream(chatID int64) (*chatStream, bool) {
	cs.mu.Lock()
	chatStream, ok := cs.streams[chatID]
	if !ok {
		return nil, false
	}
	cs.mu.Unlock()

	return chatStream, true
}

func (cs *ChatStreams) RemoveMsgStream(chatID, userID int64) {
	chatStream, ok := cs.getChatStream(chatID)
	if ok {
		chatStream.removeMsgStream(userID)
	}
}

// data race safly
func (cs *ChatStreams) RemoveChatStream(chatID int64) {
	cs.mu.Lock()
	chatStream, ok := cs.streams[chatID]
	if ok {
		delete(cs.streams, chatID)
	}
	cs.mu.Unlock()

	if !ok {
		return
	}

	chatStream.mu.Lock()
	userIDs := make([]int64, 0, len(chatStream.stream))
	for userID := range chatStream.stream {
		userIDs = append(userIDs, userID)
	}
	chatStream.mu.Unlock()

	logger.Debug(context.Background(), "start closing/deleting msg streams(may be caused because of deleting whole chat, removing user from the chat)", zap.Int64("chatID", chatID))
	for _, userID := range userIDs {
		chatStream.removeMsgStream(userID)
	}
}

// data race safly
func (cs *ChatStreams) MsgStream(chatID, userID int64) chan *model.Message {
	return cs.setChatStream(chatID).msgStream(userID)
}

// data race safly
func (cs *ChatStreams) AddMsgToChatStream(chatID int64, msg *model.Message) {
	cs.setChatStream(chatID).addMsg(msg)
}
