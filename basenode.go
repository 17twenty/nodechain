package nodechain

import "sync"

var (
	idMu   sync.Mutex
	nextID int
)

func nextNodeID() int {
	idMu.Lock()
	defer idMu.Unlock()
	id := nextID
	nextID++
	return id
}

type BaseNode struct {
	id         int
	successors map[Action][]Node
}

func NewBaseNode() BaseNode {
	return BaseNode{
		id:         nextNodeID(),
		successors: make(map[Action][]Node),
	}
}

func (b *BaseNode) ID() int { return b.id }

func (b *BaseNode) On(action Action, next Node) Node {
	if action == "" {
		action = DefaultAction
	}
	b.successors[action] = append(b.successors[action], next)
	return next
}

func (b *BaseNode) GetNextNodes(action Action) []Node {
	if action == "" {
		action = DefaultAction
	}
	return b.successors[action]
}
