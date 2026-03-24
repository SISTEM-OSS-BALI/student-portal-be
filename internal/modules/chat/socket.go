package chat

import (
	"errors"
	"net/http"
	"strings"
	"time"

	socketio "github.com/googollee/go-socket.io"

	"github.com/username/gin-gorm-api/internal/modules/auth"
)

type SocketServer struct {
	server  *socketio.Server
	service *Service
}

type SocketAck struct {
	OK      bool                `json:"ok"`
	Error   string              `json:"error,omitempty"`
	Message *MessageResponseDTO `json:"message,omitempty"`
}

type SocketJoinPayload struct {
	ConversationID string `json:"conversation_id"`
}

type SocketSendPayload struct {
	ConversationID string   `json:"conversation_id"`
	Type           string   `json:"type"`
	Text           *string  `json:"text"`
	ReplyToID      *string  `json:"reply_to_id"`
	MentionUserIDs []string `json:"mention_user_ids"`
}

type SocketReadPayload struct {
	ConversationID string     `json:"conversation_id"`
	At             *time.Time `json:"at"`
}

func NewSocketServer(service *Service) (*SocketServer, error) {
	server := socketio.NewServer(nil)
	socket := &SocketServer{
		server:  server,
		service: service,
	}
	socket.register()
	return socket, nil
}

func (s *SocketServer) Start() {
	go s.server.Serve()
}

func (s *SocketServer) Close() error {
	return s.server.Close()
}

func (s *SocketServer) Handler() http.Handler {
	return s.server
}

func (s *SocketServer) BroadcastMessage(conversationID string, message MessageResponseDTO) {
	if s == nil {
		return
	}
	if strings.TrimSpace(conversationID) == "" {
		return
	}
	s.server.BroadcastToRoom("/", conversationID, "message:new", message)
}

func (s *SocketServer) register() {
	s.server.OnConnect("/", func(conn socketio.Conn) error {
		token, err := extractToken(conn)
		if err != nil {
			return err
		}
		claims, err := auth.ParseToken(token)
		if err != nil {
			return err
		}
		conn.SetContext(claims)
		return nil
	})

	s.server.OnEvent("/", "join", func(conn socketio.Conn, payload SocketJoinPayload) SocketAck {
		claims, ok := conn.Context().(*auth.Claims)
		if !ok {
			return SocketAck{OK: false, Error: "unauthorized"}
		}
		if strings.TrimSpace(payload.ConversationID) == "" {
			return SocketAck{OK: false, Error: "conversation_id is required"}
		}
		ok, err := s.service.IsMember(payload.ConversationID, claims.UserID)
		if err != nil {
			return SocketAck{OK: false, Error: err.Error()}
		}
		if !ok {
			return SocketAck{OK: false, Error: "not a member of this conversation"}
		}
		conn.Join(payload.ConversationID)
		return SocketAck{OK: true}
	})

	s.server.OnEvent("/", "leave", func(conn socketio.Conn, payload SocketJoinPayload) SocketAck {
		if strings.TrimSpace(payload.ConversationID) == "" {
			return SocketAck{OK: false, Error: "conversation_id is required"}
		}
		conn.Leave(payload.ConversationID)
		return SocketAck{OK: true}
	})

	s.server.OnEvent("/", "message:send", func(conn socketio.Conn, payload SocketSendPayload) SocketAck {
		claims, ok := conn.Context().(*auth.Claims)
		if !ok {
			return SocketAck{OK: false, Error: "unauthorized"}
		}
		if strings.TrimSpace(payload.ConversationID) == "" {
			return SocketAck{OK: false, Error: "conversation_id is required"}
		}

		message, err := s.service.SendMessage(payload.ConversationID, claims.UserID, SendMessageDTO{
			Type:           payload.Type,
			Text:           payload.Text,
			ReplyToID:      payload.ReplyToID,
			MentionUserIDs: payload.MentionUserIDs,
		})
		if err != nil {
			return SocketAck{OK: false, Error: err.Error()}
		}

		dto := NewMessageResponseDTO(message)
		s.server.BroadcastToRoom("/", payload.ConversationID, "message:new", dto)
		return SocketAck{OK: true, Message: &dto}
	})

	s.server.OnEvent("/", "message:read", func(conn socketio.Conn, payload SocketReadPayload) SocketAck {
		claims, ok := conn.Context().(*auth.Claims)
		if !ok {
			return SocketAck{OK: false, Error: "unauthorized"}
		}
		if strings.TrimSpace(payload.ConversationID) == "" {
			return SocketAck{OK: false, Error: "conversation_id is required"}
		}
		at := time.Now()
		if payload.At != nil {
			at = *payload.At
		}
		if err := s.service.MarkRead(payload.ConversationID, claims.UserID, at); err != nil {
			return SocketAck{OK: false, Error: err.Error()}
		}
		return SocketAck{OK: true}
	})

	s.server.OnError("/", func(conn socketio.Conn, err error) {
		_ = err
	})
}

func extractToken(conn socketio.Conn) (string, error) {
	if conn == nil {
		return "", errors.New("missing connection")
	}

	// Best practice: prefer HttpOnly cookie auth.
	req := http.Request{Header: conn.RemoteHeader()}
	if c, err := req.Cookie("sp_access_token"); err == nil {
		if token := strings.TrimSpace(c.Value); token != "" {
			return token, nil
		}
	}

	// Optional dev fallback: token query.
	url := conn.URL()
	if token := strings.TrimSpace((&url).Query().Get("token")); token != "" {
		return token, nil
	}

	authHeader := conn.RemoteHeader().Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing auth token")
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid Authorization header")
	}
	return strings.TrimSpace(parts[1]), nil
}
