package nodechain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type AgentNode struct {
	BaseNode
	Provider  LLMProvider
	StateKey  string
	OutputKey string
}

func NewAgentNode(provider LLMProvider, stateKey, outputKey string) *AgentNode {
	return &AgentNode{
		BaseNode:  NewBaseNode(),
		Provider:  provider,
		StateKey:  stateKey,
		OutputKey: outputKey,
	}
}

func (n *AgentNode) TypeName() string { return "AgentNode" }

func (n *AgentNode) Run(ctx context.Context, mem *Memory) ([]Trigger, error) {

	// ----------------------------
	// 1. Load state history
	// ----------------------------
	rawState, ok := mem.Local[n.StateKey]
	if !ok {
		rawState = []string{}
		mem.Local[n.StateKey] = rawState
	}

	state, ok := rawState.([]string)
	if !ok {
		return nil, errors.New("AgentNode: state must be []string")
	}

	// ----------------------------
	// 2. Load the user task
	// ----------------------------
	taskRaw, ok := mem.Local["task"]
	if !ok {
		return nil, fmt.Errorf("AgentNode: no 'task' found in memory")
	}

	taskString, ok := taskRaw.(string)
	if !ok {
		return nil, errors.New("AgentNode: task must be a string")
	}

	// ----------------------------
	// 3. Build improved agent prompt
	// ----------------------------

	prompt := `
You are an autonomous agent with access to two tools:

TOOLS:
1) web_search(query: string) -> returns JSON
2) docker_exec(cmd: string) -> runs commands inside a persistent Ubuntu container at /workspace

REQUIREMENTS:
- Think step-by-step.
- Use tools when necessary.
- Always return ONLY JSON.
- NEVER include explanations outside JSON.
- When the task is complete, return:
  {"action":"final","response":"..."}
- Otherwise use:
  {"action":"tool","tool":"web_search"|"docker_exec","input": ...}

HISTORY:
` + buildHistory(state) + `

USER TASK:
` + taskString + `

Now produce the next action strictly in JSON format.
`

	// ----------------------------
	// 4. Call the LLM
	// ----------------------------
	resp, err := n.Provider.Chat(ctx, []LLMMessage{
		{Role: "system", Content: "Follow the instructions carefully."},
		{Role: "user", Content: prompt},
	})
	if err != nil {
		return nil, err
	}

	// Debug (optional)
	// fmt.Println("[AgentNode] LLM Output:", resp.Text)

	// ----------------------------
	// 5. Parse JSON output
	// ----------------------------
	var parsed struct {
		Action   string `json:"action"`
		Tool     string `json:"tool,omitempty"`
		Input    any    `json:"input,omitempty"`
		Response string `json:"response,omitempty"`
	}

	if err := json.Unmarshal([]byte(resp.Text), &parsed); err != nil {
		return nil, fmt.Errorf("AgentNode: invalid JSON from LLM: %v\nRaw output:\n%s", err, resp.Text)
	}

	defer fmt.Printf("[Agent] LLM decision: %s\n", resp.Text)

	mem.Local[n.OutputKey] = parsed

	// Update state history
	mem.Local[n.StateKey] = append(state, resp.Text)

	// ----------------------------
	// 6. Branch control flow
	// ----------------------------
	switch parsed.Action {
	case "tool":
		return []Trigger{{
			Action: "tool",
			ForkingData: map[string]any{
				"tool":  parsed.Tool,
				"input": parsed.Input,
			},
		}}, nil

	case "final":
		mem.Local["final_answer"] = parsed.Response
		return []Trigger{{Action: "final"}}, nil

	default:
		return nil, fmt.Errorf("AgentNode: unknown action '%s'", parsed.Action)
	}
}

// Helper to format history text
func buildHistory(lines []string) string {
	out := ""
	for _, l := range lines {
		out += "- " + l + "\n"
	}
	return out
}
