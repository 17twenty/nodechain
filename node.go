package nodechain

import "context"

type Action string

const DefaultAction Action = "default"

type Trigger struct {
	Action      Action
	ForkingData map[string]any
}

type Node interface {
	ID() int
	TypeName() string

	Run(ctx context.Context, mem *Memory) ([]Trigger, error)

	On(action Action, next Node) Node
	GetNextNodes(action Action) []Node
}
