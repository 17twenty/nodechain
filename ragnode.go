package nodechain

import (
	"context"
	"fmt"
	"strings"
)

// EmbedQueryNode: reads a query string from memory and stores its embedding.
type EmbedQueryNode struct {
	BaseNode
	Embedder     Embedder
	QueryKey     string // memory key for the query text
	EmbeddingKey string // memory key for the query embedding
}

func NewEmbedQueryNode(embedder Embedder, queryKey, embeddingKey string) *EmbedQueryNode {
	return &EmbedQueryNode{
		BaseNode:     NewBaseNode(),
		Embedder:     embedder,
		QueryKey:     queryKey,
		EmbeddingKey: embeddingKey,
	}
}

func (n *EmbedQueryNode) TypeName() string { return "EmbedQueryNode" }

func (n *EmbedQueryNode) Run(ctx context.Context, mem *Memory) ([]Trigger, error) {
	raw, ok := mem.Local[n.QueryKey]
	if !ok {
		raw, ok = mem.Global[n.QueryKey]
	}
	if !ok {
		return nil, fmt.Errorf("EmbedQueryNode: query not found at key '%s'", n.QueryKey)
	}

	query, ok := raw.(string)
	if !ok {
		return nil, fmt.Errorf("EmbedQueryNode: query at key '%s' is not a string", n.QueryKey)
	}

	embs, err := n.Embedder.EmbedText(ctx, []string{query})
	if err != nil {
		return nil, err
	}
	if len(embs) != 1 {
		return nil, fmt.Errorf("EmbedQueryNode: expected 1 embedding, got %d", len(embs))
	}

	mem.Local[n.EmbeddingKey] = embs[0]

	return []Trigger{
		{Action: DefaultAction, ForkingData: map[string]any{}},
	}, nil
}

// RetrieveNode: uses a vector store and a query embedding to get top-k documents.
type RetrieveNode struct {
	BaseNode
	Store        VectorStore
	EmbeddingKey string // key where the query embedding is stored
	K            int
	ResultKey    string // key where retrieved docs will be stored ([]Document)
}

func NewRetrieveNode(store VectorStore, embeddingKey, resultKey string, k int) *RetrieveNode {
	return &RetrieveNode{
		BaseNode:     NewBaseNode(),
		Store:        store,
		EmbeddingKey: embeddingKey,
		K:            k,
		ResultKey:    resultKey,
	}
}

func (n *RetrieveNode) TypeName() string { return "RetrieveNode" }

func (n *RetrieveNode) Run(ctx context.Context, mem *Memory) ([]Trigger, error) {
	raw, ok := mem.Local[n.EmbeddingKey]
	if !ok {
		raw, ok = mem.Global[n.EmbeddingKey]
	}
	if !ok {
		return nil, fmt.Errorf("RetrieveNode: embedding not found at key '%s'", n.EmbeddingKey)
	}

	emb, ok := raw.([]float32)
	if !ok {
		return nil, fmt.Errorf("RetrieveNode: value at key '%s' is not []float32", n.EmbeddingKey)
	}

	docs, err := n.Store.Search(ctx, emb, n.K)
	if err != nil {
		return nil, err
	}

	mem.Local[n.ResultKey] = docs

	return []Trigger{
		{Action: DefaultAction, ForkingData: map[string]any{}},
	}, nil
}

// RAGPromptNode: builds a prompt using the query and retrieved documents.
type RAGPromptNode struct {
	BaseNode
	QueryKey   string // query text
	ContextKey string // []Document
	PromptKey  string // where to store final prompt
}

func NewRAGPromptNode(queryKey, contextKey, promptKey string) *RAGPromptNode {
	return &RAGPromptNode{
		BaseNode:   NewBaseNode(),
		QueryKey:   queryKey,
		ContextKey: contextKey,
		PromptKey:  promptKey,
	}
}

func (n *RAGPromptNode) TypeName() string { return "RAGPromptNode" }

func (n *RAGPromptNode) Run(ctx context.Context, mem *Memory) ([]Trigger, error) {
	rawQ, ok := mem.Local[n.QueryKey]
	if !ok {
		rawQ, ok = mem.Global[n.QueryKey]
	}
	if !ok {
		return nil, fmt.Errorf("RAGPromptNode: query not found at key '%s'", n.QueryKey)
	}

	query, ok := rawQ.(string)
	if !ok {
		return nil, fmt.Errorf("RAGPromptNode: query at key '%s' is not a string", n.QueryKey)
	}

	rawCtx, ok := mem.Local[n.ContextKey]
	if !ok {
		rawCtx, ok = mem.Global[n.ContextKey]
	}
	if !ok {
		return nil, fmt.Errorf("RAGPromptNode: contexts not found at key '%s'", n.ContextKey)
	}

	docs, ok := rawCtx.([]Document)
	if !ok {
		return nil, fmt.Errorf("RAGPromptNode: value at key '%s' is not []Document", n.ContextKey)
	}

	var b strings.Builder
	b.WriteString("You are a helpful assistant. Use ONLY the following context to answer the question.\n\n")
	b.WriteString("Context:\n")
	for i, d := range docs {
		b.WriteString(fmt.Sprintf("[%d] %s\n", i+1, d.Text))
	}
	b.WriteString("\nQuestion:\n")
	b.WriteString(query)
	b.WriteString("\n\nAnswer:")

	mem.Local[n.PromptKey] = b.String()

	return []Trigger{
		{Action: DefaultAction, ForkingData: map[string]any{}},
	}, nil
}
