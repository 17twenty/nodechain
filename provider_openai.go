package nodechain

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type OpenAIProvider struct {
	Client *openai.Client
	Model  string
}

func (p *OpenAIProvider) Name() string { return "openai" }

func (p *OpenAIProvider) Chat(ctx context.Context, msgs []LLMMessage) (LLMResponse, error) {
	oaMsgs := make([]openai.ChatCompletionMessage, len(msgs))
	for i, m := range msgs {
		oaMsgs[i] = openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}

	resp, err := p.Client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    p.Model,
		Messages: oaMsgs,
	})
	if err != nil {
		return LLMResponse{}, err
	}

	return LLMResponse{
		Text: resp.Choices[0].Message.Content,
	}, nil
}
