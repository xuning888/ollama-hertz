package chat

import "github.com/tmc/langchaingo/llms"

type ChatReq struct {
	Content          string `json:"content"`
	LlmTimeoutSecond int    `json:"llmTimeoutSecond"`
}

type ChatWithSessonReq struct {
	Content          string `json:"content"`
	LlmTimeoutSecond int    `json:"llmTimeoutSecond"`
	UserId           string `json:"userId"`
	SessionId        string `json:"sessionId"`
	MaxWindows       int    `json:"maxWindows"`
	LlmModel         string `json:"llmModel"`
}

type Role string

const (
	User      Role = "user"
	Assistant Role = "assistant"
	System    Role = "system"
	Function  Role = "function"
)

func (r Role) LlmsRole() llms.ChatMessageType {
	switch r {
	case User:
		return llms.ChatMessageTypeHuman
	case Assistant:
		return llms.ChatMessageTypeAI
	case System:
		return llms.ChatMessageTypeSystem
	case Function:
		return llms.ChatMessageTypeFunction
	default:
		return llms.ChatMessageTypeGeneric
	}
}

type Content struct {
	Role      Role   `json:"role"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func NewContext(role Role, message string, timestamp int64) *Content {
	return &Content{
		Role:      role,
		Message:   message,
		Timestamp: timestamp,
	}
}
