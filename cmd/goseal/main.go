package main

import (
	"os"

	"github.com/jimmysharp/goseal"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	config, err := goseal.ParseConfig(".goseal.yml")
	if err != nil {
		os.Exit(1)
	}

	a := goseal.NewAnalyzer(config)

	singlechecker.Main(a)
}
