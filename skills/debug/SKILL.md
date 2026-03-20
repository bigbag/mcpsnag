---
name: debug
description: Debug and test MCP servers over HTTP - send methods like tools/list or tools/call, inspect responses, manage sessions, and troubleshoot connectivity
user_invocable: true
---

# mcpsnag: MCP Server Debugging

You are debugging an MCP (Model Context Protocol) server using the `mcpsnag` CLI tool.

## Prerequisites

- `mcpsnag` must be installed and on PATH (or available at `./bin/mcpsnag`)
- The target MCP server must be reachable over HTTP

## How to use

### List available tools

```bash
mcpsnag http://localhost:8000/mcp -d '{"method":"tools/list"}'
```

### Call a tool with arguments

```bash
mcpsnag http://localhost:8000/mcp -d '{"method":"tools/call","params":{"name":"search","arguments":{"query":"test"}}}'
```

### List resources or prompts

```bash
mcpsnag http://localhost:8000/mcp -d '{"method":"resources/list"}'
mcpsnag http://localhost:8000/mcp -d '{"method":"prompts/list"}'
```

### With authentication

```bash
mcpsnag http://localhost:8000/mcp -H "Authorization: Bearer $TOKEN" -d '{"method":"tools/list"}'
```

### Verbose mode for debugging

Use `-v` to see HTTP request/response details and session negotiation:

```bash
mcpsnag http://localhost:8000/mcp -v -d '{"method":"tools/list"}'
```

### Compact output for piping

Use `-c` for single-line JSON suitable for piping to `jq`:

```bash
mcpsnag http://localhost:8000/mcp -c -d '{"method":"tools/list"}' | jq '.tools[].name'
```

### Initialize only (get session ID)

```bash
mcpsnag http://localhost:8000/mcp --init-only
```

### Reuse an existing session

```bash
mcpsnag http://localhost:8000/mcp --session <session-id> -d '{"method":"tools/list"}'
```

### Raw mode (skip auto-initialization)

Send a raw JSON-RPC request without the MCP handshake:

```bash
mcpsnag http://localhost:8000/mcp --raw -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}'
```

### Disable streaming

Wait for the full response instead of streaming SSE events:

```bash
mcpsnag http://localhost:8000/mcp --no-stream -d '{"method":"tools/list"}'
```

## Common MCP methods

- `initialize` - Handshake with server (done automatically unless `--raw`)
- `tools/list` - List available tools
- `tools/call` - Call a tool: `{"name":"<tool>","arguments":{...}}`
- `resources/list` - List available resources
- `resources/read` - Read a resource: `{"uri":"<resource-uri>"}`
- `prompts/list` - List available prompts
- `prompts/get` - Get a prompt: `{"name":"<prompt>","arguments":{...}}`
- `completion/complete` - Request completions

## CLI flags reference

- `-d`, `--data` - JSON body with method + params (required unless `--init-only`)
- `-H`, `--header` - HTTP header, repeatable (format: `Key: Value`)
- `-c`, `--compact` - Compact single-line JSON output
- `-v`, `--verbose` - Show HTTP request/response details
- `--raw` - Skip auto-initialization, send raw JSON-RPC
- `--session` - Reuse existing session ID
- `--init-only` - Only initialize and print session ID
- `--no-stream` - Wait for full response instead of streaming
- `--timeout` - Request timeout (default: 30s)

## Debugging workflow

1. Start with `--init-only` to verify the server is reachable and responds to the MCP handshake
2. Use `tools/list` to discover available tools and their input schemas
3. Call specific tools with `tools/call` providing required arguments
4. If something fails, add `-v` to inspect the full HTTP exchange
5. Use `--raw` to bypass initialization when testing non-standard endpoints
6. Use compact mode (`-c`) with `jq` for programmatic analysis
