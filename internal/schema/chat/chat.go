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
	MaxWindows       int    `json:"maxWindows"`
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
	Role    Role   `json:"role"`
	Message string `json:"message"`
}

func NewContext(role Role, message string) *Content {
	return &Content{
		Role:    role,
		Message: message,
	}
}
