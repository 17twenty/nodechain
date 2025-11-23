package nodechain

import (
	"context"
	"fmt"
)

type Tool interface {
	Name() string
	Run(input any) (any, error)
}

type ToolNode struct {
	BaseNode
	Tools        map[string]Tool
	ToolNameKey  string // memory key that holds which tool to call
	ToolInputKey string // memory key for tool input
	ResultKey    string // where to store tool output
}

func NewToolNode(tools map[string]Tool, toolNameKey, toolInputKey, resultKey string) *ToolNode {
	return &ToolNode{
		BaseNode:     NewBaseNode(),
		Tools:        tools,
		ToolNameKey:  toolNameKey,
		ToolInputKey: toolInputKey,
		ResultKey:    resultKey,
	}
}

func (n *ToolNode) TypeName() string { return "ToolNode" }

func (n *ToolNode) Run(ctx context.Context, mem *Memory) ([]Trigger, error) {
	toolNameRaw, ok := mem.Local[n.ToolNameKey]
	if !ok {
		toolNameRaw, ok = mem.Global[n.ToolNameKey]
	}
	if !ok {
		return nil, fmt.Errorf("ToolNode: tool name key '%s' missing", n.ToolNameKey)
	}

	toolName, ok := toolNameRaw.(string)
	if !ok {
		return nil, fmt.Errorf("ToolNode: tool name must be string")
	}

	tool, ok := n.Tools[toolName]
	if !ok {
		return nil, fmt.Errorf("ToolNode: tool '%s' not found", toolName)
	}

	input := mem.Local[n.ToolInputKey]
	result, err := tool.Run(input)
	if err != nil {
		// still record the error as an observation
		result = map[string]any{"error": err.Error()}
	}

	mem.Local[n.ResultKey] = result

	// 🔑 NEW: write a summary into "state" so AgentNode sees it
	rawState, _ := mem.Local["state"]
	state, _ := rawState.([]string)

	summary := fmt.Sprintf("TOOL %s INPUT=%v OUTPUT=%v", toolName, input, result)
	state = append(state, summary)
	mem.Local["state"] = state

	fmt.Printf("[Tool %s] input=%v output=%v\n", toolName, input, result)

	return []Trigger{
		{Action: DefaultAction, ForkingData: map[string]any{}},
	}, nil
}
