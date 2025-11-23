package nodechain

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type Embedder interface {
	EmbedText(ctx context.Context, texts []string) ([][]float32, error)
	Name() string
}

type OpenAIEmbedder struct {
	Client *openai.Client
	Model  openai.EmbeddingModel
}

func NewOpenAIEmbedder(client *openai.Client, model openai.EmbeddingModel) *OpenAIEmbedder {
	return &OpenAIEmbedder{
		Client: client,
		Model:  model,
	}
}

func (e *OpenAIEmbedder) Name() string { return "openai-embedder" }

func (e *OpenAIEmbedder) EmbedText(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("EmbedText: no texts provided")
	}

	resp, err := e.Client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Model: e.Model,
		Input: texts,
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Data) != len(texts) {
		return nil, fmt.Errorf("EmbedText: got %d embeddings for %d texts", len(resp.Data), len(texts))
	}

	out := make([][]float32, len(resp.Data))
	for i, d := range resp.Data {
		out[i] = d.Embedding
	}
	return out, nil
}
