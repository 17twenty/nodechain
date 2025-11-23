package nodechain

type ExecutionTree struct {
	Order     int                        `json:"order"`
	Type      string                     `json:"type"`
	Triggered map[Action][]ExecutionTree `json:"triggered,omitempty"`
}
