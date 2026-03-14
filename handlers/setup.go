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

// any file that matches will start a setup command
var orderedRules = []struct {
	file    string
	command string
}{
	{".tool-versions", "mise install"},
	{"go.mod", "go mod download && go install tool"},
	{"yarn.lock", "yarn install"},
	{"package.json", "npm install"},
	{"requirements.txt", "pip install -r requirements.txt"},
	{"pyproject.toml", "pip install ."},
	{"Gemfile", "bundle install"},
	{"Cargo.toml", "cargo fetch"},
	{"mix.exs", "mix deps.get"},
}

// first match only
var setupScriptCandidates = []string{
	"setup.sh",
	"setup",
	"bin/setup",
	"script/setup",
	"scripts/setup",
	"scripts/setup.sh",
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
		pathResult := map[string]any{}

		manifest := buildManifestCommand(mountPath)
		script, err := findSetupScript(mountPath)

		var command string
		switch {
		case err == nil && manifest != "":
			// setup language first
			command = ". " + script + " && " + manifest
			pathResult["setup_script"] = script
		case manifest != "" && err != nil:
			command = manifest
		case err == nil:
			command = ". " + script
			pathResult["setup_script"] = script
		}

		if command == "" {
			pathResult["error"] = "no supported rule found; project may use an unsupported language or package manager"
		} else {
			j := h.startJob(command, mountPath)
			pathResult["job_id"] = j.id
		}

		result[mountPath] = pathResult
	}

	b, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError("failed to encode result"), nil
	}

	return mcp.NewToolResultText(string(b)), nil
}

func buildManifestCommand(projectPath string) string {
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

// findSetupScript checks known candidate paths under projectPath and returns
// the first regular file found.
func findSetupScript(projectPath string) (string, error) {
	for _, candidate := range setupScriptCandidates {
		full := filepath.Join(projectPath, candidate)
		info, err := os.Stat(full)
		if err == nil && info.Mode().IsRegular() {
			return full, nil
		}
	}
	return "", fmt.Errorf("not found")
}
