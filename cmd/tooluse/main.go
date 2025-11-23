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

	// OpenAI
	client := openai.NewClient(apiKey)
	provider := &nc.OpenAIProvider{
		Client: client,
		Model:  "gpt-4o-mini",
	}

	// Docker
	docker := nc.NewDockerManager("ubuntu:latest", "./workspace")
	if err := docker.Start(); err != nil {
		panic(err)
	}
	defer docker.Stop()

	// Optional bootstrap
	docker.Run(`apt-get update -y`)
	docker.Run(`apt-get install -y curl wget python3 python3-pip`)

	// Tools
	tools := map[string]nc.Tool{
		"docker_exec": &nc.DockerExecTool{Docker: docker},
		"web_search":  &nc.SerperSearchTool{},
	}

	agent := nc.NewAgentNode(provider, "state", "agent_output")
	toolNode := nc.NewToolNode(tools, "tool", "input", "tool_result")
	finalLLM := nc.NewLLMNode(provider, "final_answer", "final_response")
	printNode := &nc.PrintNode{Keys: []string{"final_response"}}

	// Wiring
	agent.On("tool", toolNode)
	toolNode.On(nc.DefaultAction, agent)
	agent.On("final", finalLLM)
	finalLLM.On(nc.DefaultAction, printNode)

	// Seed task
	startTask := nc.NewValueNode("task",
		"Find a website about koalas and download the largest image to /workspace.",
	)

	startTask.On(nc.DefaultAction, agent)
	flow := nc.NewFlow(startTask)
	flow.MaxVisits = 50

	// Run it (this line FIXES your error)
	tree, err := flow.Run(ctx, nil)
	if err != nil {
		panic(err)
	}

	// Print execution tree
	bts, _ := json.MarshalIndent(tree, "", "  ")
	fmt.Println(string(bts))
}
