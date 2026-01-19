package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bigbag/mcpsnag/internal/client"
	"github.com/bigbag/mcpsnag/internal/output"
	"github.com/bigbag/mcpsnag/internal/protocol"
)

type headerFlags []string

func (h *headerFlags) String() string {
	return strings.Join(*h, ", ")
}

func (h *headerFlags) Set(value string) error {
	*h = append(*h, value)
	return nil
}

var flagsWithValues = map[string]bool{
	"-d": true, "--data": true, "-data": true,
	"-H": true, "--header": true, "-header": true,
	"--session": true, "-session": true,
	"--timeout": true, "-timeout": true,
}

func reorderArgs(args []string) []string {
	if len(args) <= 1 {
		return args
	}

	var flags []string
	var positional []string

	i := 1
	for i < len(args) {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			if strings.Contains(arg, "=") {
				flags = append(flags, arg)
				i++
			} else if flagsWithValues[arg] && i+1 < len(args) {
				flags = append(flags, arg, args[i+1])
				i += 2
			} else {
				flags = append(flags, arg)
				i++
			}
		} else {
			positional = append(positional, arg)
			i++
		}
	}

	result := []string{args[0]}
	result = append(result, flags...)
	result = append(result, positional...)
	return result
}

func main() {
	var (
		data     string
		headers  headerFlags
		raw      bool
		session  string
		initOnly bool
		compact  bool
		noStream bool
		verbose  bool
		timeout  time.Duration
	)

	flag.StringVar(&data, "d", "", "JSON body (method + params)")
	flag.StringVar(&data, "data", "", "JSON body (method + params)")
	flag.Var(&headers, "H", "HTTP header (repeatable)")
	flag.Var(&headers, "header", "HTTP header (repeatable)")
	flag.BoolVar(&raw, "raw", false, "Skip auto-initialization")
	flag.StringVar(&session, "session", "", "Use existing session ID")
	flag.BoolVar(&initOnly, "init-only", false, "Only initialize, print session")
	flag.BoolVar(&compact, "c", false, "Compact JSON output")
	flag.BoolVar(&compact, "compact", false, "Compact JSON output")
	flag.BoolVar(&noStream, "no-stream", false, "Wait for full response")
	flag.BoolVar(&verbose, "v", false, "Show request/response details")
	flag.BoolVar(&verbose, "verbose", false, "Show request/response details")
	flag.DurationVar(&timeout, "timeout", 30*time.Second, "Request timeout")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: mcpsnag [options] <url>\n\n")
		fmt.Fprintf(os.Stderr, "A curl-like CLI for testing MCP servers over HTTP.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  mcpsnag http://localhost:3000/mcp -d '{\"method\":\"tools/list\"}'\n")
		fmt.Fprintf(os.Stderr, "  mcpsnag http://localhost:3000/mcp -H \"Authorization: Bearer token\" -d '{\"method\":\"tools/list\"}'\n")
		fmt.Fprintf(os.Stderr, "  mcpsnag http://localhost:3000/mcp --init-only\n")
	}

	os.Args = reorderArgs(os.Args)
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "error: URL is required")
		flag.Usage()
		os.Exit(1)
	}

	url := flag.Arg(0)
	printer := output.NewPrinter(os.Stdout, os.Stderr, compact, verbose)

	if !initOnly && data == "" {
		fmt.Fprintln(os.Stderr, "error: -d/--data is required (or use --init-only)")
		flag.Usage()
		os.Exit(1)
	}

	headerMap := make(map[string]string)
	for _, h := range headers {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "warning: invalid header format %q (expected 'Key: Value')\n", h)
			continue
		}
		headerMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	c := client.New(client.Options{
		Endpoint:  url,
		Headers:   headerMap,
		SessionID: session,
		Timeout:   timeout,
		Stream:    !noStream,
	})

	if raw {
		runRaw(c, printer, data)
		return
	}

	if session == "" {
		printer.PrintVerbose("* Initializing MCP session...")
		result, err := c.Initialize()
		if err != nil {
			printer.PrintError(fmt.Errorf("initialization failed: %w", err))
			os.Exit(1)
		}
		printer.PrintVerbose("* Connected to %s %s", result.ServerInfo.Name, result.ServerInfo.Version)
		printer.PrintVerbose("* Session ID: %s", c.Session().ID)
	}

	if initOnly {
		printer.PrintSessionInfo(c.Session().ID)
		return
	}

	runRequest(c, printer, data)
}

func runRaw(c *client.Client, printer *output.Printer, data string) {
	resp, sessionID, err := c.RawRequest([]byte(data), func(r protocol.Response) error {
		return printer.PrintRawJSON(r.Result)
	})
	if err != nil {
		printer.PrintError(err)
		os.Exit(1)
	}

	if sessionID != "" {
		printer.PrintVerbose("* Session ID: %s", sessionID)
	}

	if resp != nil {
		if resp.Error != nil {
			printer.PrintJSON(resp.Error)
			os.Exit(1)
		}
		if resp.Result != nil {
			printer.PrintRawJSON(resp.Result)
		}
	}
}

func runRequest(c *client.Client, printer *output.Printer, data string) {
	var userReq protocol.UserRequest
	if err := json.Unmarshal([]byte(data), &userReq); err != nil {
		printer.PrintError(fmt.Errorf("invalid JSON: %w", err))
		os.Exit(1)
	}

	if userReq.Method == "" {
		printer.PrintError(fmt.Errorf("missing 'method' field in request"))
		os.Exit(1)
	}

	resp, err := c.Request(userReq.Method, userReq.Params, func(r protocol.Response) error {
		if r.Result != nil {
			return printer.PrintRawJSON(r.Result)
		}
		return nil
	})
	if err != nil {
		if resp != nil && resp.Error != nil {
			printer.PrintJSON(resp.Error)
		} else {
			printer.PrintError(err)
		}
		os.Exit(1)
	}

	if resp != nil && resp.Result != nil {
		printer.PrintRawJSON(resp.Result)
	}
}
