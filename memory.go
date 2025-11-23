package nodechain

type Memory struct {
	Global map[string]any
	Local  map[string]any
}

func NewMemory(global map[string]any) *Memory {
	if global == nil {
		global = make(map[string]any)
	}
	return &Memory{
		Global: global,
		Local:  make(map[string]any),
	}
}

func (m *Memory) Clone(fork map[string]any) *Memory {
	newLocal := make(map[string]any, len(m.Local)+len(fork))
	for k, v := range m.Local {
		newLocal[k] = v
	}
	for k, v := range fork {
		newLocal[k] = v
	}
	return &Memory{
		Global: m.Global, // shared
		Local:  newLocal,
	}
}
