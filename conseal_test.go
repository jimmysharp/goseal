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
			name: "config_default",
		},
		{
			name: "config_allow_same_package",
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
