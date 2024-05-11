package memory

import (
	"context"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/prompts"
)

var defaultSummaryTemplate = `
Progressively summarize the lines of conversation provided, adding onto the previous summary returning a new summary.

EXAMPLE
Current summary:
The human asks what the AI thinks of artificial intelligence. The AI thinks artificial intelligence is a force for good.

New lines of conversation:
Human: Why do you think artificial intelligence is a force for good?
AI: Because artificial intelligence will help humans reach their full potential.

New summary:
The human asks what the AI thinks of artificial intelligence. The AI thinks artificial intelligence is a force for good because it will help humans reach their full potential.
END OF EXAMPLE

Current summary:
{{.summary}}

New lines of conversation:
{{.new_lines}}

New summary:
`

var defaultSummaryPrompt = prompts.NewPromptTemplate(defaultSummaryTemplate, []string{"summary", "new_lines"})

type ConversationSummary struct {
	memory.ConversationBuffer
	llm           llms.Model
	summaryPrompt prompts.FormatPrompter
}

func (c *ConversationSummary) MemoryVariables(ctx context.Context) []string {
	return c.ConversationBuffer.MemoryVariables(ctx)
}

func (c *ConversationSummary) LoadMemoryVariables(ctx context.Context, inputs map[string]any) (map[string]any, error) {
	messages, err := c.ChatHistory.Messages(ctx)

	if err != nil {
		return nil, err
	}

	if c.ReturnMessages {
		return map[string]any{
			c.MemoryKey: messages,
		}, nil
	}

	return nil, nil
}

func (c *ConversationSummary) predictNewSummary(ctx context.Context,
	messages []llms.ChatMessage, existingSummary string) (summary string, err error) {

	chain := chains.NewLLMChain(c.llm, c.summaryPrompt)

	bufferString, err2 := llms.GetBufferString(messages, c.HumanPrefix, c.AIPrefix)
	if err != nil {
		return "", err2
	}

	summary, err = chains.Predict(ctx, chain, map[string]any{
		"summary":   existingSummary,
		"new_lines": bufferString,
	})
	return
}

func NewConversationSummary(llm llms.Model, options ...memory.ConversationBufferOption) (*ConversationSummary, error) {

	css := &ConversationSummary{
		llm:                llm,
		summaryPrompt:      defaultSummaryPrompt,
		ConversationBuffer: *applyBufferOptions(options...),
	}

	if css == nil {
		return nil, ErrorModelNil
	}

	return css, nil
}
