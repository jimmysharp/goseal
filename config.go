package conseal

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type Config struct {
	StructPackages   []*regexp.Regexp // Regex patterns for packages containing target structs (if empty, all packages are targeted)
	Constructors     []*regexp.Regexp // Regex patterns for constructor functions
	IgnoreFiles      []*regexp.Regexp // Regex patterns for files to ignore
	AllowSamePackage bool             // Whether to allow construction within the same package
}

func NewConfig(structPackages, constructors, ignoreFiles []*regexp.Regexp, allowSamePackage bool) (*Config, error) {
	if structPackages == nil {
		structPackages = []*regexp.Regexp{}
	}

	if len(constructors) == 0 {
		constructors = []*regexp.Regexp{regexp.MustCompile("^New.*")}
	}

	config := &Config{
		StructPackages:   structPackages,
		Constructors:     constructors,
		IgnoreFiles:      ignoreFiles,
		AllowSamePackage: allowSamePackage,
	}

	return config, nil
}

type yamlConfig struct {
	StructPackages   []string `yaml:"struct-packages"`
	Constructors     []string `yaml:"constructors"`
	IgnoreFiles      []string `yaml:"ignore-files"`
	AllowSamePackage bool     `yaml:"allow-same-package"`
}

func ParseConfig(path string) (*Config, error) {
	if path == "" {
		path = ".conseal.yml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewConfig(nil, nil, nil, false)
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var yamlCfg yamlConfig
	if err := yaml.Unmarshal(data, &yamlCfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return yamlCfg.compile()
}

func (y *yamlConfig) compile() (*Config, error) {
	var structPackages []*regexp.Regexp
	var constructors []*regexp.Regexp
	var ignoreFiles []*regexp.Regexp

	// struct-packages
	for _, pattern := range y.StructPackages {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid struct-packages pattern '%s': %w", pattern, err)
		}
		structPackages = append(structPackages, re)
	}

	// constroctors
	for _, pattern := range y.Constructors {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid constructors pattern '%s': %w", pattern, err)
		}
		constructors = append(constructors, re)
	}

	// ignore-files
	for _, pattern := range y.IgnoreFiles {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid ignore-files pattern '%s': %w", pattern, err)
		}
		ignoreFiles = append(ignoreFiles, re)
	}

	return NewConfig(structPackages, constructors, ignoreFiles, y.AllowSamePackage)
}
