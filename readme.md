# NodeChain
Lightweight, composable flow engine for building LLM agents, workflows, and RAG systems — in Go.

NodeChain is a minimal alternative to heavy agent frameworks like LangChain or LangGraph.

It provides:
- Composable nodes
- Directed execution flows
- Stateful memory
- Tool integration
- LLM-driven agents
- Docker-isolated tool execution
- Real web search (by Serper, get a free API key)
- Deterministic execution trees
- NodeChain is designed for clarity, small surface area, and maximum control, while enabling sophisticated LLM-powered systems.

## Features

### Composable computational graph

Define Nodes and connect them using `.On(action, nextNode)`.

### Memory system

A Memory instance gives each step access to:

- Global shared state
- Local isolated state per branch
- automatic cloning for branching flows

### LLM-powered Agents

AgentNode implements an autonomous reasoning loop using:

- ReAct-style tool calls
- strict JSON actions
- looping control with cycle safety
- memory of previous steps

### ToolNode for real tool execution

Add tools like:

- docker_exec → run inside an Ubuntu sandbox
- web_search → Serper (Google search)
- (extendable) filesystem tools, Python, HTTP, embeddings, etc.

### Stateful Docker environment

A persistent Ubuntu container acts as a safe, isolated computation sandbox.

Agents can:

- install packages
- run curl/wget
- write Python scripts
- download files
- inspect output
- iterate

Fully sandboxed autonomous behavior!

NodeChain is intentionally small and easy to understand — ideal for:

- Backend services
- Dev agents
- Evaluation agents
- Workflow orchestration
- Autonomous tool users
- Embedding/RAG pipelines
