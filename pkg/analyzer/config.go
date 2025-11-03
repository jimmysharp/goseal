package analyzer

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Packages         []*regexp.Regexp // Regex patterns for target packages
	Constructors     []*regexp.Regexp // Regex patterns for constructor functions
	IgnoreFiles      []*regexp.Regexp // Regex patterns for files to ignore
	AllowSamePackage bool             // Whether to allow construction within the same package
}

func NewConfig(packages, constructors, ignoreFiles []*regexp.Regexp, allowSamePackage bool) (*Config, error) {
	if packages == nil {
		packages = []*regexp.Regexp{}
	}

	if len(constructors) == 0 {
		constructors = []*regexp.Regexp{regexp.MustCompile("^New.*")}
	}

	config := &Config{
		Packages:         packages,
		Constructors:     constructors,
		IgnoreFiles:      ignoreFiles,
		AllowSamePackage: allowSamePackage,
	}

	return config, nil
}

type yamlConfig struct {
	Packages         []string `yaml:"packages"`
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
	var packages []*regexp.Regexp
	var constructors []*regexp.Regexp
	var ignoreFiles []*regexp.Regexp

	// packages
	for _, pattern := range y.Packages {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid packages pattern '%s': %w", pattern, err)
		}
		packages = append(packages, re)
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

	return NewConfig(packages, constructors, ignoreFiles, y.AllowSamePackage)
}
