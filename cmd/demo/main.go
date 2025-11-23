package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	nc "nodechain"

	"github.com/sashabaranov/go-openai"
)

func main() {
	ctx := context.Background()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		panic("Please set OPENAI_API_KEY")
	}

	client := openai.NewClient(apiKey)

	provider := &nc.OpenAIProvider{
		Client: client,
		Model:  "gpt-4o-mini", // fast + cheap for demos
	}

	// A node that inserts a prompt into memory
	init := nc.NewValueNode("prompt",
		"Explain NodeChain in one short sentence.")

	// A node that calls the LLM and stores result in memory["answer"]
	llm := nc.NewLLMNode(provider, "prompt", "answer")

	// A node that prints the answer
	print := &nc.PrintNode{Keys: []string{"answer"}}

	init.On(nc.DefaultAction, llm)
	llm.On(nc.DefaultAction, print)

	flow := nc.NewFlow(init)

	tree, err := flow.Run(ctx, nil)
	if err != nil {
		panic(err)
	}

	bts, _ := json.MarshalIndent(tree, "", "  ")
	fmt.Println("\n--- EXECUTION TREE ---")
	fmt.Println(string(bts))
}
