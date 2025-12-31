package goseal

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/goccy/go-yaml"
)

type InitScope string

const (
	InitScopeAny              InitScope = "any"
	InitScopeInTargetPackages InitScope = "in-target-packages"
	InitScopeSamePackage      InitScope = "same-package"
)

type MutationScope string

const (
	MutationScopeAny              MutationScope = "any"
	MutationScopeInTargetPackages MutationScope = "in-target-packages"
	MutationScopeReceiver         MutationScope = "receiver"
	MutationScopeSamePackage      MutationScope = "same-package"
	MutationScopeNever            MutationScope = "never"
)

type Config struct {
	TargetPackages []*regexp.Regexp // Regex patterns for packages containing target structs (if empty, all packages are targeted)
	ExcludeStructs []*regexp.Regexp // Regex patterns for struct names to exclude from protection
	FactoryNames   []*regexp.Regexp // Regex patterns for factory function names (if empty, all function names are allowed)
	InitScope      InitScope        // Scope for struct initialization
	MutationScope  MutationScope    // Scope for field mutation
	IgnoreFiles    []*regexp.Regexp // Regex patterns for files to ignore
}

func (c *Config) UnmarshalJSON(data []byte) error {
	type rawConfig struct {
		TargetPackages []string `json:"target-packages"`
		ExcludeStructs []string `json:"exclude-structs"`
		FactoryNames   []string `json:"factory-names"`
		InitScope      string   `json:"init-scope"`
		MutationScope  string   `json:"mutation-scope"`
		IgnoreFiles    []string `json:"ignore-files"`
	}

	var raw rawConfig
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var targetPackages []*regexp.Regexp
	var excludeStructs []*regexp.Regexp
	var factoryNames []*regexp.Regexp
	var ignoreFiles []*regexp.Regexp

	// target-packages
	if raw.TargetPackages != nil {
		targetPackages = make([]*regexp.Regexp, len(raw.TargetPackages))
		for i, pattern := range raw.TargetPackages {
			re, err := regexp.Compile(pattern)
			if err != nil {
				return fmt.Errorf("invalid target-packages pattern '%s': %w", pattern, err)
			}
			targetPackages[i] = re
		}
	} else {
		targetPackages = []*regexp.Regexp{}
	}

	// exclude-structs
	if raw.ExcludeStructs != nil {
		excludeStructs = make([]*regexp.Regexp, len(raw.ExcludeStructs))
		for i, pattern := range raw.ExcludeStructs {
			re, err := regexp.Compile(pattern)
			if err != nil {
				return fmt.Errorf("invalid exclude-structs pattern '%s': %w", pattern, err)
			}
			excludeStructs[i] = re
		}
	} else {
		excludeStructs = []*regexp.Regexp{}
	}

	// factory-names
	if raw.FactoryNames != nil {
		factoryNames = make([]*regexp.Regexp, len(raw.FactoryNames))
		for i, pattern := range raw.FactoryNames {
			re, err := regexp.Compile(pattern)
			if err != nil {
				return fmt.Errorf("invalid factory-names pattern '%s': %w", pattern, err)
			}
			factoryNames[i] = re
		}
	} else {
		factoryNames = []*regexp.Regexp{}
	}

	// ignore-files
	if raw.IgnoreFiles != nil {
		ignoreFiles = make([]*regexp.Regexp, len(raw.IgnoreFiles))
		for i, pattern := range raw.IgnoreFiles {
			re, err := regexp.Compile(pattern)
			if err != nil {
				return fmt.Errorf("invalid ignore-files pattern '%s': %w", pattern, err)
			}
			ignoreFiles[i] = re
		}
	} else {
		ignoreFiles = []*regexp.Regexp{}
	}

	cfg, err := NewConfig(
		targetPackages,
		excludeStructs,
		factoryNames,
		InitScope(raw.InitScope),
		MutationScope(raw.MutationScope),
		ignoreFiles,
	)
	if err != nil {
		return err
	}

	*c = *cfg
	return nil
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
	case MutationScopeAny, MutationScopeInTargetPackages, MutationScopeReceiver, MutationScopeSamePackage, MutationScopeNever:
		return nil
	default:
		return fmt.Errorf("invalid mutation-scope: %s (must be 'any', 'in-target-packages', 'receiver', 'same-package', or 'never')", scope)
	}
}

func ParseFromYAML(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.UnmarshalWithOptions(data, &cfg, yaml.UseJSONUnmarshaler()); err != nil {
		return nil, fmt.Errorf("failed to parse config data: %w", err)
	}

	// Recover defaults for fully empty YAML
	finalCfg, err := NewConfig(
		cfg.TargetPackages,
		cfg.ExcludeStructs,
		cfg.FactoryNames,
		cfg.InitScope,
		cfg.MutationScope,
		cfg.IgnoreFiles,
	)
	if err != nil {
		return nil, err
	}

	return finalCfg, nil
}

func ParseConfig(path string) (*Config, error) {
	if path == "" {
		path = ".goseal.yml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewConfig(nil, nil, nil, "", "", nil)
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return ParseFromYAML(data)
}
