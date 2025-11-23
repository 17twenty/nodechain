package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	nc "nodechain"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	ctx := context.Background()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		panic("Please set OPENAI_API_KEY")
	}

	client := openai.NewClient(apiKey)

	embedder := nc.NewOpenAIEmbedder(client, "text-embedding-3-small")
	provider := &nc.OpenAIProvider{
		Client: client,
		Model:  "gpt-4o-mini",
	}

	store := nc.NewInMemoryVectorStore()
	if err := indexDemoDocs(ctx, embedder, store); err != nil {
		panic(err)
	}

	// Flow:
	// ValueNode("query") -> EmbedQueryNode -> RetrieveNode -> RAGPromptNode
	// -> LLMNode -> PrintNode

	queryNode := nc.NewValueNode("query",
		"How does NodeChain help build RAG systems?")

	embedNode := nc.NewEmbedQueryNode(embedder, "query", "query_embedding")
	retrieveNode := nc.NewRetrieveNode(store, "query_embedding", "contexts", 3)
	ragPromptNode := nc.NewRAGPromptNode("query", "contexts", "prompt")
	llmNode := nc.NewLLMNode(provider, "prompt", "answer")
	printNode := &nc.PrintNode{Keys: []string{"answer"}}

	queryNode.On(nc.DefaultAction, embedNode)
	embedNode.On(nc.DefaultAction, retrieveNode)
	retrieveNode.On(nc.DefaultAction, ragPromptNode)
	ragPromptNode.On(nc.DefaultAction, llmNode)
	llmNode.On(nc.DefaultAction, printNode)

	flow := nc.NewFlow(queryNode)

	tree, err := flow.Run(ctx, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n--- EXECUTION TREE ---")
	bts, _ := json.MarshalIndent(tree, "", "  ")
	fmt.Println(string(bts))
}

func indexDemoDocs(ctx context.Context, embedder nc.Embedder, store nc.VectorStore) error {
	texts := []string{
		"NodeChain is a Go library for building typed, async workflows from simple nodes.",
		"RAG stands for Retrieval Augmented Generation, a technique where relevant documents are retrieved and used as context for an LLM.",
		"NodeChain supports nodes for embedding, retrieval, and LLM calls, making it easy to build RAG systems.",
	}

	embs, err := embedder.EmbedText(ctx, texts)
	if err != nil {
		return err
	}

	docs := make([]nc.Document, len(texts))
	for i, t := range texts {
		docs[i] = nc.Document{
			ID:        fmt.Sprintf("doc-%d", i+1),
			Text:      t,
			Metadata:  map[string]any{"source": "demo"},
			Embedding: embs[i],
		}
	}

	return store.Add(ctx, docs)
}
