package conseal_test

import (
	"path/filepath"
	"testing"

	"github.com/jimmysharp/conseal"
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
			name: "config/default",
		},
		{
			name: "config/same_package",
		},
		{
			name: "config/in_target_packages",
		},
		{
			name: "config/only_factory_names",
		},
		{
			name: "config/mutation_scope_any",
		},
		{
			name: "config/mutation_scope_never",
		},
		{
			name: "config/exclude_structs",
		},
		{
			name: "unsupported",
		},
		{
			name: "generated",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testdataDir := filepath.Join(analysistest.TestData(), tt.name)

			configPath := filepath.Join(testdataDir, ".conseal.yml")
			config, err := conseal.ParseConfig(configPath)
			require.NoError(t, err)

			a := conseal.NewAnalyzer(config)

			// All test cases are under module "example.com/testproject"
			analysistest.Run(t, testdataDir, a, "example.com/testproject/...")
		})
	}
}
