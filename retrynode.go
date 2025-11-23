package nodechain

import (
	"context"
	"errors"
	"time"
)

type RetryNode struct {
	BaseNode
	Inner      Node
	MaxRetries int
	RetryDelay time.Duration
}

func NewRetryNode(inner Node, maxRetries int, retryDelay time.Duration) *RetryNode {
	return &RetryNode{
		BaseNode:   NewBaseNode(),
		Inner:      inner,
		MaxRetries: maxRetries,
		RetryDelay: retryDelay,
	}
}

func (n *RetryNode) TypeName() string {
	return "RetryNode(" + n.Inner.TypeName() + ")"
}

// RetryNode simply wraps Inner.Run with retry logic.
func (n *RetryNode) Run(ctx context.Context, mem *Memory) ([]Trigger, error) {
	var lastErr error

	for attempt := 0; attempt < n.MaxRetries; attempt++ {
		triggers, err := n.Inner.Run(ctx, mem)
		if err == nil {
			return triggers, nil
		}

		lastErr = err

		if attempt < n.MaxRetries-1 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(n.RetryDelay):
			}
		}
	}
	return nil, errors.New("RetryNode: all attempts failed: " + lastErr.Error())
}

// successors propagate to inner node's successors
func (n *RetryNode) GetNextNodes(action Action) []Node {
	return n.Inner.GetNextNodes(action)
}

// On attaches successors to inner node
func (n *RetryNode) On(action Action, next Node) Node {
	return n.Inner.On(action, next)
}
