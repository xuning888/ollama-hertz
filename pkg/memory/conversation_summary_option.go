package memory

import (
	"github.com/pkg/errors"
	"github.com/tmc/langchaingo/memory"
)

var (
	ErrorModelNil = errors.New("llm is nil")
)

func applyBufferOptions(opts ...memory.ConversationBufferOption) *memory.ConversationBuffer {
	m := &memory.ConversationBuffer{
		ReturnMessages: false,
		InputKey:       "",
		OutputKey:      "",
		HumanPrefix:    "Human",
		AIPrefix:       "AI",
		MemoryKey:      "history",
	}

	for _, opt := range opts {
		opt(m)
	}

	if m.ChatHistory == nil {
		m.ChatHistory = memory.NewChatMessageHistory()
	}
	return m
}
