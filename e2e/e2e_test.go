package e2e

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

// golangciLintOutput represents the JSON output from golangci-lint.
type golangciLintOutput struct {
	Issues []issue `json:"Issues"`
}

type issue struct {
	FromLinter string   `json:"FromLinter"`
	Text       string   `json:"Text"`
	Pos        position `json:"Pos"`
}

type position struct {
	Filename string `json:"Filename"`
	Line     int    `json:"Line"`
	Column   int    `json:"Column"`
}

func TestGolangciLintPlugin(t *testing.T) {
	// Check that golangci-lint is installed.
	golangciLint, err := exec.LookPath("golangci-lint")
	require.NoError(t, err, "golangci-lint is not installed")

	// Resolve directories.
	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "failed to get caller information")
	e2eDir, err := filepath.Abs(filepath.Dir(thisFile))
	require.NoError(t, err)
	testProjectDir := filepath.Join(e2eDir, "testproject")

	// Build custom golangci-lint binary.
	// golangci-lint custom reads .custom-gcl.yml from the working directory (e2e/testproject/).
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "custom-gcl")
	buildCmd := exec.Command(golangciLint, "custom", "--destination", tmpDir)
	buildCmd.Dir = testProjectDir
	buildCmd.Env = append(os.Environ(), "GOFLAGS=")
	out, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "failed to build custom golangci-lint binary:\n%s", string(out))

	// Run lint on the test project.
	// Write JSON output to a file to avoid mixing with text output on stdout.
	jsonOutputPath := filepath.Join(tmpDir, "output.json")
	runCmd := exec.Command(binaryPath, "run", "--output.json.path", jsonOutputPath, "./...")
	runCmd.Dir = testProjectDir
	// golangci-lint returns exit code 1 when issues are found, so we don't check err.
	runOut, _ := runCmd.CombinedOutput()

	// Parse JSON output into struct with only stable fields.
	// Unstable fields (Offset, SourceLines, etc.) and the Report section
	// are silently dropped by json.Unmarshal.
	rawOutput, err := os.ReadFile(jsonOutputPath)
	require.NoError(t, err, "failed to read JSON output file; run output:\n%s", string(runOut))
	var actual golangciLintOutput
	require.NoError(t, json.Unmarshal(rawOutput, &actual), "failed to parse golangci-lint JSON output:\n%s", string(rawOutput))

	// Load golden file and compare.
	goldenPath := filepath.Join(testProjectDir, "expected.json")
	goldenBytes, err := os.ReadFile(goldenPath)
	require.NoError(t, err, "failed to read golden file: %s", goldenPath)
	var expected golangciLintOutput
	require.NoError(t, json.Unmarshal(goldenBytes, &expected))

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("issues mismatch (-expected +actual):\n%s", diff)
	}
}
