package nodechain

import "context"

type LLMMessage struct {
	Role    string
	Content string
}

type LLMResponse struct {
	Text string
}

type LLMProvider interface {
	Chat(ctx context.Context, messages []LLMMessage) (LLMResponse, error)
	Name() string
}
