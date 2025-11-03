package main

import (
	"os"

	"github.com/jimmysharp/conseal/pkg/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	config, err := analyzer.ParseConfig(".conseal.yml")
	if err != nil {
		os.Exit(1)
	}

	a := analyzer.NewAnalyzer(config)

	singlechecker.Main(a)
}
