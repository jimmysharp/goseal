package conseal

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type InitScope string

const (
	InitScopeAny              InitScope = "any"
	InitScopeInTargetPackages InitScope = "in-target-packages"
	InitScopeSamePackage      InitScope = "same-package"
)

type MutationScope string

const (
	MutationScopeAny         MutationScope = "any"
	MutationScopeReceiver    MutationScope = "receiver"
	MutationScopeSamePackage MutationScope = "same-package"
	MutationScopeNever       MutationScope = "never"
)

type Config struct {
	TargetPackages []*regexp.Regexp // Regex patterns for packages containing target structs (if empty, all packages are targeted)
	ExcludeStructs []*regexp.Regexp // Regex patterns for struct names to exclude from protection
	FactoryNames   []*regexp.Regexp // Regex patterns for factory function names (if empty, all function names are allowed)
	InitScope      InitScope        // Scope for struct initialization
	MutationScope  MutationScope    // Scope for field mutation
	IgnoreFiles    []*regexp.Regexp // Regex patterns for files to ignore
}

func NewConfig(
	targetPackages []*regexp.Regexp,
	excludeStructs []*regexp.Regexp,
	factoryNames []*regexp.Regexp,
	initScope InitScope,
	mutationScope MutationScope,
	ignoreFiles []*regexp.Regexp,
) (*Config, error) {
	// Set default values
	if targetPackages == nil {
		targetPackages = []*regexp.Regexp{}
	}
	if excludeStructs == nil {
		excludeStructs = []*regexp.Regexp{}
	}
	if factoryNames == nil {
		factoryNames = []*regexp.Regexp{}
	}
	if initScope == "" {
		initScope = InitScopeSamePackage
	}
	if mutationScope == "" {
		mutationScope = MutationScopeReceiver
	}
	if ignoreFiles == nil {
		ignoreFiles = []*regexp.Regexp{}
	}

	// Validate scopes
	if err := validateInitScope(initScope); err != nil {
		return nil, err
	}
	if err := validateMutationScope(mutationScope); err != nil {
		return nil, err
	}

	return &Config{
		TargetPackages: targetPackages,
		ExcludeStructs: excludeStructs,
		FactoryNames:   factoryNames,
		InitScope:      initScope,
		MutationScope:  mutationScope,
		IgnoreFiles:    ignoreFiles,
	}, nil
}

func validateInitScope(scope InitScope) error {
	switch scope {
	case InitScopeAny, InitScopeInTargetPackages, InitScopeSamePackage:
		return nil
	default:
		return fmt.Errorf("invalid init-scope: %s (must be 'any', 'in-target-packages', or 'same-package')", scope)
	}
}

func validateMutationScope(scope MutationScope) error {
	switch scope {
	case MutationScopeAny, MutationScopeReceiver, MutationScopeSamePackage, MutationScopeNever:
		return nil
	default:
		return fmt.Errorf("invalid mutation-scope: %s (must be 'any', 'receiver', 'same-package', or 'never')", scope)
	}
}

type yamlConfig struct {
	TargetPackages []string `yaml:"target-packages"`
	ExcludeStructs []string `yaml:"exclude-structs"`
	FactoryNames   []string `yaml:"factory-names"`
	InitScope      string   `yaml:"init-scope"`
	MutationScope  string   `yaml:"mutation-scope"`
	IgnoreFiles    []string `yaml:"ignore-files"`
}

func ParseConfig(path string) (*Config, error) {
	if path == "" {
		path = ".conseal.yml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewConfig(nil, nil, nil, "", "", nil)
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
	var targetPackages []*regexp.Regexp
	var excludeStructs []*regexp.Regexp
	var factoryNames []*regexp.Regexp
	var ignoreFiles []*regexp.Regexp

	// target-packages
	for _, pattern := range y.TargetPackages {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid target-packages pattern '%s': %w", pattern, err)
		}
		targetPackages = append(targetPackages, re)
	}

	// exclude-structs
	for _, pattern := range y.ExcludeStructs {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid exclude-structs pattern '%s': %w", pattern, err)
		}
		excludeStructs = append(excludeStructs, re)
	}

	// factory-names
	for _, pattern := range y.FactoryNames {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid factory-names pattern '%s': %w", pattern, err)
		}
		factoryNames = append(factoryNames, re)
	}

	// ignore-files
	for _, pattern := range y.IgnoreFiles {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid ignore-files pattern '%s': %w", pattern, err)
		}
		ignoreFiles = append(ignoreFiles, re)
	}

	return NewConfig(
		targetPackages,
		excludeStructs,
		factoryNames,
		InitScope(y.InitScope),
		MutationScope(y.MutationScope),
		ignoreFiles,
	)
}
