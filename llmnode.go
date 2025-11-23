package nodechain

import (
	"context"
	"errors"
	"fmt"
)

type LLMNode struct {
	BaseNode
	Provider LLMProvider
	InputKey string
	StoreKey string
	System   string
}

func NewLLMNode(provider LLMProvider, inputKey, storeKey string) *LLMNode {
	return &LLMNode{
		BaseNode: NewBaseNode(),
		Provider: provider,
		InputKey: inputKey,
		StoreKey: storeKey,
		System:   "You are a helpful agent inside NodeChain.",
	}
}

func (n *LLMNode) TypeName() string { return "LLMNode" }

func (n *LLMNode) Run(ctx context.Context, mem *Memory) ([]Trigger, error) {
	raw, ok := mem.Local[n.InputKey]
	if !ok {
		raw, ok = mem.Global[n.InputKey]
	}
	if !ok {
		return nil, fmt.Errorf("LLMNode: no prompt found at key '%s'", n.InputKey)
	}

	prompt, ok := raw.(string)
	if !ok {
		return nil, errors.New("LLMNode: prompt must be a string")
	}

	resp, err := n.Provider.Chat(ctx, []LLMMessage{
		{Role: "system", Content: n.System},
		{Role: "user", Content: prompt},
	})
	if err != nil {
		return nil, err
	}

	mem.Local[n.StoreKey] = resp.Text

	return []Trigger{
		{Action: DefaultAction, ForkingData: map[string]any{}},
	}, nil
}
