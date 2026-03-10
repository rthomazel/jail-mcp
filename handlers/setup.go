package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

var orderedRules = []struct {
	file    string
	command string
}{
	{".tool-versions", "mise install"},
	{"go.mod", "go mod download"},
	{"tools.go", "go install -tags tools ./..."},
	{"yarn.lock", "yarn install"},
	{"package.json", "npm install"},
	{"requirements.txt", "pip install -r requirements.txt"},
	{"pyproject.toml", "pip install ."},
	{"Gemfile", "bundle install"},
	{"Cargo.toml", "cargo fetch"},
	{"mix.exs", "mix deps.get"},
}

func (h *Handler) HandleSetup(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	raw, _ := req.Params.Arguments["paths"].([]any)
	if len(raw) == 0 {
		return mcp.NewToolResultError("missing required parameter: paths"), nil
	}

	paths := make([]string, 0, len(raw))
	for _, v := range raw {
		str, ok := v.(string)
		if !ok {
			return mcp.NewToolResultError(fmt.Sprintf("paths must be strings, got %T", v)), nil
		}
		paths = append(paths, str)
	}

	result := map[string]any{}

	for _, mountPath := range paths {
		command := buildSetupCommand(mountPath)
		if command == "" {
			result[mountPath] = map[string]any{
				"error": "no supported manifests found; project may use an unsupported language or package manager",
			}
			continue
		}

		j := h.startJob(command, mountPath)
		result[mountPath] = map[string]any{"job_id": j.id}
	}

	b, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError("failed to encode result"), nil
	}

	return mcp.NewToolResultText(string(b)), nil
}

func buildSetupCommand(projectPath string) string {
	var commands []string
	for _, rule := range orderedRules {
		_, statErr := os.Stat(filepath.Join(projectPath, rule.file))
		if statErr == nil {
			commands = append(commands, rule.command)
		}
	}
	if len(commands) == 0 {
		return ""
	}
	return strings.Join(commands, " && ")
}
