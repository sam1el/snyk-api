package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

const (
	configDirName  = "snyk-api"
	configFileName = "config.yaml"
)

// Profile holds user-defined defaults for a CLI profile.
type Profile struct {
	Token      string `yaml:"token,omitempty"`
	APIURL     string `yaml:"api_url,omitempty"`
	RestAPIURL string `yaml:"rest_api_url,omitempty"`
	APIVersion string `yaml:"api_version,omitempty"`
	OrgID      string `yaml:"org_id,omitempty"`
	GroupID    string `yaml:"group_id,omitempty"`
	ProjectID  string `yaml:"project_id,omitempty"`
	Output     string `yaml:"output,omitempty"`
	PageSize   int    `yaml:"page_size,omitempty"`
	Debug      bool   `yaml:"debug,omitempty"`
}

// File represents the persisted config file.
type File struct {
	Current  string             `yaml:"current,omitempty"`
	Profiles map[string]Profile `yaml:"profiles,omitempty"`
}

// FlagOverrides captures CLI-level overrides.
type FlagOverrides struct {
	Profile    string
	Token      string
	APIURL     string
	RestAPIURL string
	APIVersion string
	OrgID      string
	GroupID    string
	ProjectID  string
	Output     string
	PageSize   *int
	Debug      *bool
}

// Resolved holds the merged configuration after applying precedence.
type Resolved struct {
	ProfileName string
	Token       string
	APIURL      string
	RestAPIURL  string
	APIVersion  string
	OrgID       string
	GroupID     string
	ProjectID   string
	Output      string
	PageSize    int
	Debug       bool
}

// DefaultPath returns ~/.config/snyk-api/config.yaml (or OS equivalent).
func DefaultPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("config dir: %w", err)
	}
	return filepath.Join(dir, configDirName, configFileName), nil
}

// LoadFile loads a config file. Missing file returns an empty config.
func LoadFile(path string) (*File, error) {
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &File{Profiles: map[string]Profile{}}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	cfg := &File{}
	if err := yaml.Unmarshal(content, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if cfg.Profiles == nil {
		cfg.Profiles = map[string]Profile{}
	}
	return cfg, nil
}

// SaveFile writes the config file, creating parent directories if needed.
func SaveFile(path string, cfg *File) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

// Resolve merges flags > env > profile > defaults. The profile name itself
// is resolved by flags.profile, then SNYK_API_PROFILE, then file.Current,
// falling back to "default".
func Resolve(fileCfg *File, flags FlagOverrides, env map[string]string) Resolved {
	defaults := Resolved{
		APIURL:     "https://api.snyk.io",
		RestAPIURL: "https://api.snyk.io/rest",
		APIVersion: "2025-11-05",
		Output:     "json",
		PageSize:   100,
		Debug:      false,
	}

	profileName := resolveProfileName(flags.Profile, env, fileCfg)
	profile := Profile{}
	if fileCfg != nil && fileCfg.Profiles != nil {
		if p, ok := fileCfg.Profiles[profileName]; ok {
			profile = p
		}
	}

	// Start with defaults, then profile, then env, then flags
	res := defaults
	res.ProfileName = profileName

	applyProfile(&res, profile)
	applyEnv(&res, env)
	applyFlags(&res, flags)

	return res
}

func resolveProfileName(flagVal string, env map[string]string, fileCfg *File) string {
	if flagVal != "" {
		return flagVal
	}
	if env != nil {
		if v, ok := env["SNYK_API_PROFILE"]; ok && v != "" {
			return v
		}
	}
	if fileCfg != nil && fileCfg.Current != "" {
		return fileCfg.Current
	}
	return "default"
}

func applyProfile(res *Resolved, p Profile) {
	if p.Token != "" {
		res.Token = p.Token
	}
	if p.APIURL != "" {
		res.APIURL = p.APIURL
	}
	if p.RestAPIURL != "" {
		res.RestAPIURL = p.RestAPIURL
	}
	if p.APIVersion != "" {
		res.APIVersion = p.APIVersion
	}
	if p.OrgID != "" {
		res.OrgID = p.OrgID
	}
	if p.GroupID != "" {
		res.GroupID = p.GroupID
	}
	if p.ProjectID != "" {
		res.ProjectID = p.ProjectID
	}
	if p.Output != "" {
		res.Output = p.Output
	}
	if p.PageSize > 0 {
		res.PageSize = p.PageSize
	}
	if p.Debug {
		res.Debug = true
	}
}

func applyEnv(res *Resolved, env map[string]string) {
	if env == nil {
		return
	}

	stringMappings := map[string]func(string){
		"SNYK_TOKEN":       func(v string) { res.Token = v },
		"SNYK_API":         func(v string) { res.APIURL = v },
		"SNYK_REST_API":    func(v string) { res.RestAPIURL = v },
		"SNYK_API_VERSION": func(v string) { res.APIVersion = v },
		"SNYK_ORG_ID":      func(v string) { res.OrgID = v },
		"SNYK_GROUP_ID":    func(v string) { res.GroupID = v },
		"SNYK_PROJECT_ID":  func(v string) { res.ProjectID = v },
		"SNYK_OUTPUT":      func(v string) { res.Output = v },
	}

	for key, apply := range stringMappings {
		if v, ok := env[key]; ok && v != "" {
			apply(v)
		}
	}

	if v, ok := env["SNYK_PAGE_SIZE"]; ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			res.PageSize = n
		}
	}

	if v, ok := env["SNYK_DEBUG"]; ok && v != "" {
		if parsed, err := strconv.ParseBool(v); err == nil {
			res.Debug = parsed
		}
	}
}

func applyFlags(res *Resolved, flags FlagOverrides) {
	if flags.Token != "" {
		res.Token = flags.Token
	}
	if flags.APIURL != "" {
		res.APIURL = flags.APIURL
	}
	if flags.RestAPIURL != "" {
		res.RestAPIURL = flags.RestAPIURL
	}
	if flags.APIVersion != "" {
		res.APIVersion = flags.APIVersion
	}
	if flags.OrgID != "" {
		res.OrgID = flags.OrgID
	}
	if flags.GroupID != "" {
		res.GroupID = flags.GroupID
	}
	if flags.ProjectID != "" {
		res.ProjectID = flags.ProjectID
	}
	if flags.Output != "" {
		res.Output = flags.Output
	}
	if flags.PageSize != nil && *flags.PageSize > 0 {
		res.PageSize = *flags.PageSize
	}
	if flags.Debug != nil {
		res.Debug = *flags.Debug
	}
}
