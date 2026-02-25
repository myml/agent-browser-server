package main

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"agent-browser-server - Remote agent-browser via mcp",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Add tool
	tool := mcp.NewTool("agent-browser",
		mcp.WithDescription("Execute agent-browser"),
		mcp.WithArray("args",
			mcp.Required(),
			mcp.Description("Command-line arguments to pass to agent-browser on the remote server. First, Use --help to see all available parameters."),
		),
	)

	// Add tool handler
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := request.RequireStringSlice("args")
		if err != nil {
			slog.Error("Missing 'args' parameter")
			return mcp.NewToolResultError("Missing 'args' parameter"), nil
		}
		cmdCtx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()

		session := server.ClientSessionFromContext(ctx)
		if session == nil {
			return mcp.NewToolResultError("No active session"), nil
		}
		slog.Info("session", "session id", session.SessionID())

		var stdoutBuf, stderrBuf bytes.Buffer
		cmd := exec.CommandContext(cmdCtx, "agent-browser", args...)
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf
		cmd.Env = append(os.Environ(), "AGENT_BROWSER_SESSION="+session.SessionID())
		startTime := time.Now()
		err = cmd.Run()

		exitCode := 0
		status := "success"
		if err != nil {
			status = "exit error"
			if exitError, ok := err.(*exec.ExitError); ok {
				exitCode = exitError.ExitCode()
			} else {
				exitCode = -1
				status = err.Error()
			}
		}
		response := map[string]interface{}{
			"status":         status,
			"exit_code":      exitCode,
			"stdout":         stdoutBuf.String(),
			"stderr":         stderrBuf.String(),
			"command":        strings.Join(cmd.Args, " "),
			"execution_time": time.Since(startTime),
		}
		slog.Info("Command execution completed", "command", response["command"], "exit_code", exitCode)
		return mcp.NewToolResultJSON(response)
	})

	slog.Info("MCP server initialized, serving on http://127.0.0.1:8080/mcp")
	// Start the stdio server
	if err := server.NewStreamableHTTPServer(s).Start(":8080"); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
