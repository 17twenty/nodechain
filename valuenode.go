package nodechain

import "context"

type ValueNode struct {
	BaseNode
	Key   string
	Value any
}

func NewValueNode(key string, value any) *ValueNode {
	return &ValueNode{
		BaseNode: NewBaseNode(),
		Key:      key,
		Value:    value,
	}
}

func (n *ValueNode) TypeName() string { return "ValueNode" }

func (n *ValueNode) Run(ctx context.Context, mem *Memory) ([]Trigger, error) {
	mem.Local[n.Key] = n.Value
	return []Trigger{{Action: DefaultAction}}, nil
}
