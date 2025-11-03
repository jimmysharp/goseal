package analyzer_test

import (
	"path/filepath"
	"testing"

	"github.com/jimmysharp/conseal/pkg/analyzer"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "basic",
		},
		{
			name: "default_config",
		},
		{
			name: "allow_same_package",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testdataDir := filepath.Join(analysistest.TestData(), tt.name)

			configPath := filepath.Join(testdataDir, ".conseal.yml")
			config, err := analyzer.ParseConfig(configPath)
			require.NoError(t, err)

			a := analyzer.NewAnalyzer(config)

			// All test cases are under module "example.com/testproject"
			analysistest.Run(t, testdataDir, a, "example.com/testproject/...")
		})
	}
}
