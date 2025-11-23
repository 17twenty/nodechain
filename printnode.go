package nodechain

import (
	"context"
	"fmt"
)

type PrintNode struct {
	BaseNode
	Keys []string
}

func (n *PrintNode) TypeName() string { return "PrintNode" }

func (n *PrintNode) Run(ctx context.Context, mem *Memory) ([]Trigger, error) {
	for _, k := range n.Keys {
		fmt.Printf("%s = %v\n", k, mem.Local[k])
	}
	return nil, nil
}
