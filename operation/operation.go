package operation

import (
	"fmt"
	"time"

	"gitea.com/gitea/gitea-mcp/operation/issue"
	"gitea.com/gitea/gitea-mcp/operation/label"
	"gitea.com/gitea/gitea-mcp/operation/pull"
	"gitea.com/gitea/gitea-mcp/operation/repo"
	"gitea.com/gitea/gitea-mcp/operation/search"
	"gitea.com/gitea/gitea-mcp/operation/user"
	"gitea.com/gitea/gitea-mcp/operation/version"
	"gitea.com/gitea/gitea-mcp/pkg/flag"
	"gitea.com/gitea/gitea-mcp/pkg/log"

	"github.com/mark3labs/mcp-go/server"
)

var mcpServer *server.MCPServer

func RegisterTool(s *server.MCPServer) {
	// User Tool
	s.AddTools(user.Tool.Tools()...)

	// Repo Tool
	s.AddTools(repo.Tool.Tools()...)

	// Issue Tool
	s.AddTools(issue.Tool.Tools()...)

	// Label Tool
	s.AddTools(label.Tool.Tools()...)

	// Pull Tool
	s.AddTools(pull.Tool.Tools()...)

	// Search Tool
	s.AddTools(search.Tool.Tools()...)

	// Version Tool
	s.AddTools(version.Tool.Tools()...)

	s.DeleteTools("")
}

func Run() error {
	mcpServer = newMCPServer(flag.Version)
	RegisterTool(mcpServer)
	switch flag.Mode {
	case "stdio":
		if err := server.ServeStdio(
			mcpServer,
		); err != nil {
			return err
		}
	case "sse":
		sseServer := server.NewSSEServer(
			mcpServer,
		)
		log.Infof("Gitea MCP SSE server listening on :%d", flag.Port)
		if err := sseServer.Start(fmt.Sprintf(":%d", flag.Port)); err != nil {
			return err
		}
	case "http":
		httpServer := server.NewStreamableHTTPServer(
			mcpServer,
			server.WithLogger(log.New()),
			server.WithHeartbeatInterval(30*time.Second),
			server.WithStateLess(true),
		)
		log.Infof("Gitea MCP HTTP server listening on :%d", flag.Port)
		if err := httpServer.Start(fmt.Sprintf(":%d", flag.Port)); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid transport type: %s. Must be 'stdio', 'sse' or 'http'", flag.Mode)
	}
	return nil
}

func newMCPServer(version string) *server.MCPServer {
	return server.NewMCPServer(
		"Gitea MCP Server",
		version,
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithRecovery(),
	)
}
