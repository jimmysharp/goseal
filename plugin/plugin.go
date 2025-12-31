package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"github.com/jimmysharp/goseal"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("goseal", New)
}

type Plugin struct {
	config *goseal.Config
}

func New(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[goseal.Config](settings)
	if err != nil {
		return nil, err
	}

	return &Plugin{config: &s}, nil
}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	analyzer := goseal.NewAnalyzer(p.config)

	return []*analysis.Analyzer{analyzer}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
