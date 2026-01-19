# mcpsnag

[![Go Version](https://img.shields.io/github/go-mod/go-version/bigbag/mcpsnag)](https://github.com/bigbag/mcpsnag)
[![Build](https://img.shields.io/github/actions/workflow/status/bigbag/mcpsnag/build.yaml?branch=master)](https://github.com/bigbag/mcpsnag/actions/workflows/build.yaml)
[![Release](https://img.shields.io/github/v/release/bigbag/mcpsnag)](https://github.com/bigbag/mcpsnag/releases/latest)
[![license](https://img.shields.io/github/license/bigbag/mcpsnag.svg)](https://github.com/bigbag/mcpsnag/blob/master/LICENSE)

A curl-like CLI tool for testing and debugging [MCP (Model Context Protocol)](https://modelcontextprotocol.io/docs/getting-started/intro) servers over HTTP.

## Features

- **Auto-initialization** - Handles MCP handshake automatically
- **Session management** - Reuse sessions across requests
- **SSE streaming** - Print events as they arrive
- **Pretty output** - Formatted JSON by default
- **Raw mode** - Skip initialization for custom flows
- **Verbose mode** - Show request/response details

## Quick Start

```bash
# Build
make build

# Basic request
./bin/mcpsnag http://localhost:3000/mcp -d '{"method":"tools/list"}'
```

## Installation

```bash
# Clone the repository
git clone https://github.com/bigbag/mcpsnag.git
cd mcpsnag

# Build
make build

# Or install to GOPATH/bin
make install
```

## CLI Flags

- `-d, --data` - JSON body (method + params) *required*
- `-H, --header` - HTTP header (repeatable)
- `--raw` - Skip auto-initialization
- `--session` - Use existing session ID
- `--init-only` - Only initialize, print session
- `-c, --compact` - Compact JSON output
- `--no-stream` - Wait for full response
- `-v, --verbose` - Show request/response details
- `--timeout` - Request timeout (default: 30s)

## MCP Protocol Flow

By default, mcpsnag handles the MCP initialization handshake:

```
1. POST -> initialize request
   <- capabilities + Mcp-Session-Id
2. POST -> notifications/initialized
3. POST -> user's request (with session header)
   <- response
```

Use `--raw` to skip this and send requests directly.

## Examples

### Tools

List available tools:
```bash
mcpsnag http://localhost:3000/mcp -d '{"method":"tools/list"}'
```

Call a tool with arguments:
```bash
mcpsnag http://localhost:3000/mcp -d '{"method":"tools/call","params":{"name":"search","arguments":{"query":"hello world"}}}'
```

Call a tool with complex arguments:
```bash
mcpsnag http://localhost:3000/mcp -d '{
  "method": "tools/call",
  "params": {
    "name": "create_file",
    "arguments": {
      "path": "/tmp/test.txt",
      "content": "Hello, World!"
    }
  }
}'
```

### Resources

List available resources:
```bash
mcpsnag http://localhost:3000/mcp -d '{"method":"resources/list"}'
```

Read a specific resource:
```bash
mcpsnag http://localhost:3000/mcp -d '{"method":"resources/read","params":{"uri":"file:///path/to/file.txt"}}'
```

List resource templates:
```bash
mcpsnag http://localhost:3000/mcp -d '{"method":"resources/templates/list"}'
```

Subscribe to resource changes:
```bash
mcpsnag http://localhost:3000/mcp -d '{"method":"resources/subscribe","params":{"uri":"file:///path/to/watch"}}'
```

### Prompts

List available prompts:
```bash
mcpsnag http://localhost:3000/mcp -d '{"method":"prompts/list"}'
```

Get a prompt with arguments:
```bash
mcpsnag http://localhost:3000/mcp -d '{"method":"prompts/get","params":{"name":"code_review","arguments":{"language":"go"}}}'
```

### Authentication

With Bearer token:
```bash
mcpsnag http://localhost:3000/mcp \
  -H "Authorization: Bearer your-token-here" \
  -d '{"method":"tools/list"}'
```

With API key:
```bash
mcpsnag http://localhost:3000/mcp \
  -H "X-API-Key: your-api-key" \
  -d '{"method":"tools/list"}'
```

With multiple headers:
```bash
mcpsnag http://localhost:3000/mcp \
  -H "Authorization: Bearer token" \
  -H "X-Request-ID: req-123" \
  -H "X-Tenant-ID: tenant-456" \
  -d '{"method":"tools/list"}'
```

### Session Management

Initialize and capture session:
```bash
export MCP_SESSION=$(mcpsnag http://localhost:3000/mcp --init-only | jq -r '.sessionId')
echo "Session: $MCP_SESSION"
```

Reuse session for multiple requests:
```bash
mcpsnag http://localhost:3000/mcp --session "$MCP_SESSION" -d '{"method":"tools/list"}'
mcpsnag http://localhost:3000/mcp --session "$MCP_SESSION" -d '{"method":"resources/list"}'
mcpsnag http://localhost:3000/mcp --session "$MCP_SESSION" -d '{"method":"prompts/list"}'
```

### Debugging

Verbose mode (show request/response details):
```bash
mcpsnag http://localhost:3000/mcp -v -d '{"method":"tools/list"}'
```

Raw mode (skip initialization, send custom request):
```bash
mcpsnag http://localhost:3000/mcp --raw -d '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2025-03-26",
    "capabilities": {},
    "clientInfo": {"name": "test", "version": "1.0"}
  }
}'
```

With custom timeout:
```bash
mcpsnag http://localhost:3000/mcp --timeout 60s -d '{"method":"tools/call","params":{"name":"slow_operation"}}'
```

### Output Formatting

Pretty print (default):
```bash
mcpsnag http://localhost:3000/mcp -d '{"method":"tools/list"}'
```

Compact JSON (for scripting):
```bash
mcpsnag http://localhost:3000/mcp -c -d '{"method":"tools/list"}'
```

Disable streaming (wait for complete response):
```bash
mcpsnag http://localhost:3000/mcp --no-stream -d '{"method":"tools/list"}'
```

### Scripting Patterns

Extract tool names with jq:
```bash
mcpsnag http://localhost:3000/mcp -c -d '{"method":"tools/list"}' | jq -r '.tools[].name'
```

Check if a specific tool exists:
```bash
mcpsnag http://localhost:3000/mcp -c -d '{"method":"tools/list"}' | jq -e '.tools[] | select(.name == "search")' > /dev/null && echo "Found"
```

Loop through and call multiple tools:
```bash
for tool in search fetch analyze; do
  echo "Calling $tool..."
  mcpsnag http://localhost:3000/mcp -c -d "{\"method\":\"tools/call\",\"params\":{\"name\":\"$tool\",\"arguments\":{}}}"
done
```

Save response to file:
```bash
mcpsnag http://localhost:3000/mcp -d '{"method":"tools/list"}' > tools.json
```

Pipe JSON from file:
```bash
cat request.json | xargs -0 mcpsnag http://localhost:3000/mcp -d
```

### Error Handling

Check exit code:
```bash
mcpsnag http://localhost:3000/mcp -d '{"method":"tools/call","params":{"name":"invalid"}}' || echo "Request failed"
```

Capture errors:
```bash
result=$(mcpsnag http://localhost:3000/mcp -c -d '{"method":"tools/list"}' 2>&1)
if echo "$result" | jq -e '.error' > /dev/null 2>&1; then
  echo "Error: $(echo "$result" | jq -r '.error.message')"
else
  echo "Success"
fi
```

### Testing Different Servers

Test local development server:
```bash
mcpsnag http://localhost:3000/mcp -d '{"method":"tools/list"}'
```

Test remote server:
```bash
mcpsnag https://api.example.com/mcp -H "Authorization: Bearer $TOKEN" -d '{"method":"tools/list"}'
```

Compare two servers:
```bash
diff <(mcpsnag http://localhost:3000/mcp -c -d '{"method":"tools/list"}' | jq -S .) \
     <(mcpsnag http://localhost:3001/mcp -c -d '{"method":"tools/list"}' | jq -S .)
```

## Make Commands

```bash
make build         # Build binary to bin/mcpsnag
make run           # Build and run
make run/quick     # Run without rebuild
make test          # Run tests
make test-race     # Run tests with race detection
make coverage      # Run tests with coverage report
make coverage-html # Generate HTML coverage report
make fmt           # Format code
make vet           # Run go vet
make lint          # Run fmt and vet
make tidy          # Tidy Go modules
make clean         # Remove build artifacts
make install       # Install to GOPATH/bin
make build-all     # Build for linux/darwin/windows amd64/arm64
```

## Testing

```bash
# Run all tests
make test

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
make coverage

# Run tests with race detection
make test-race
```

## License

MIT License - see [LICENSE](LICENSE) file.
