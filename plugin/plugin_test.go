package plugin_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/golangci/plugin-module-register/register"
	_ "github.com/jimmysharp/goseal/plugin"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestModulePlugin_Basic(t *testing.T) {
	newPlugin, err := register.GetPlugin("goseal")
	require.NoError(t, err)

	settings := map[string]any{
		"target-packages": []string{"example\\.com/testproject/domain.*"},
		"factory-names":   []string{"^New.*"},
		"init-scope":      "same-package",
		"mutation-scope":  "receiver",
		"ignore-files":    []string{"_test\\.go$"},
	}

	p, err := newPlugin(settings)
	require.NoError(t, err)
	require.Equal(t, register.LoadModeTypesInfo, p.GetLoadMode())

	analyzers, err := p.BuildAnalyzers()
	require.NoError(t, err)
	require.Len(t, analyzers, 1)
	require.Equal(t, "goseal", analyzers[0].Name)

	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), ".."))

	testdataDir := filepath.Join(repoRoot, "testdata", "basic")
	analysistest.Run(t, testdataDir, analyzers[0], "example.com/testproject/...")
}
