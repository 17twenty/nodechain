package nodechain

import (
	"context"
	"fmt"
)

type Flow struct {
	Start       Node
	MaxVisits   int
	visitCounts map[int]int
}

func NewFlow(start Node) *Flow {
	return &Flow{
		Start:       start,
		MaxVisits:   15,
		visitCounts: make(map[int]int),
	}
}

func (f *Flow) Run(ctx context.Context, global map[string]any) (ExecutionTree, error) {
	f.visitCounts = map[int]int{}
	mem := NewMemory(global)
	return f.runNode(ctx, f.Start, mem)
}

func (f *Flow) runNode(ctx context.Context, n Node, mem *Memory) (ExecutionTree, error) {
	id := n.ID()
	f.visitCounts[id]++
	if f.visitCounts[id] > f.MaxVisits {
		return ExecutionTree{}, fmt.Errorf("cycle limit reached for %s#%d", n.TypeName(), id)
	}

	// clone memory for this node
	cloned := mem.Clone(nil)
	triggers, err := n.Run(ctx, cloned)
	if err != nil {
		return ExecutionTree{}, err
	}

	out := ExecutionTree{
		Order: id,
		Type:  n.TypeName(),
	}

	if len(triggers) == 0 {
		return out, nil
	}

	out.Triggered = make(map[Action][]ExecutionTree)

	for _, trig := range triggers {
		nextNodes := n.GetNextNodes(trig.Action)
		if len(nextNodes) == 0 {
			out.Triggered[trig.Action] = []ExecutionTree{}
			continue
		}

		forkMem := cloned.Clone(trig.ForkingData)

		var children []ExecutionTree
		for _, nextNode := range nextNodes {
			childTree, err := f.runNode(ctx, nextNode, forkMem)
			if err != nil {
				return ExecutionTree{}, err
			}
			children = append(children, childTree)
		}
		out.Triggered[trig.Action] = children
	}

	return out, nil
}
