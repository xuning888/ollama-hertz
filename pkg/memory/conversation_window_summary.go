package memory

import (
	"context"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/schema"
)

var _ schema.Memory = (*ConversationWindowSummary)(nil)

var (
	defaultMessageSize = 2
)

type ConversationWindowSummary struct {
	memory.ConversationWindowBuffer
	llm        llms.Model
	summaryKey string
}

func (cws *ConversationWindowSummary) MemoryVariables(ctx context.Context) []string {
	return cws.ConversationBuffer.MemoryVariables(ctx)
}
func (cws *ConversationWindowSummary) LoadMemoryVariables(ctx context.Context, _ map[string]any) (map[string]any, error) {
	messages, err := cws.ChatHistory.Messages(ctx)
	if err != nil {
		return nil, err
	}

	needCutMessage, _ := cws.needCutMessage(messages)
	if needCutMessage {
		// summary and cut message
	}

	if cws.ReturnMessages {
		return map[string]any{
			cws.MemoryKey: messages,
		}, nil
	}

	bufferString, err := llms.GetBufferString(messages, cws.HumanPrefix, cws.AIPrefix)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		cws.MemoryKey: bufferString,
	}, nil
}

func (cws *ConversationWindowSummary) needCutMessage(messages []llms.ChatMessage) (need bool, cutIndex int) {
	need = len(messages) > cws.ConversationWindowSize*defaultMessageSize
	if need {
		cutIndex = len(messages) - cws.ConversationWindowSize*defaultMessageSize
	}
	return
}

func (cws *ConversationWindowSummary) summary(messages []llms.ChatMessage) {

}

func NewConversationWindowSummary(maxWindow int, llm llms.Model, summaryKey string,
	opts ...memory.ConversationBufferOption) *ConversationWindowSummary {
	cws := &ConversationWindowSummary{
		ConversationWindowBuffer: *memory.NewConversationWindowBuffer(maxWindow, opts...),
		llm:                      llm,
		summaryKey:               summaryKey,
	}
	return cws
}
