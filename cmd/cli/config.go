package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

const Filename = "beago.json"

// Config is the top-level structure of beago.json.
type Config struct {
	Layers []Layer `json:"layers"`
}

// Layer describes an architectural layer and the rules it must follow.
type Layer struct {
	Name  string   `json:"name"`
	Paths []string `json:"paths"` // glob patterns relative to project root
	Rules []string `json:"rules"`
}

// load reads beago.json from root. Returns an empty Config if the file does
// not exist so callers can proceed without a config.
func load(root string) (Config, error) {
	data, err := os.ReadFile(filepath.Join(root, Filename))
	if os.IsNotExist(err) {
		return Config{}, nil
	}
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// MatchLayer returns the first Layer whose path patterns match target.
// target and root should both be absolute paths.
// Returns nil when no layer matches or the config has no layers.
// Patterns support ** to match across directory separators.
func (c Config) MatchLayer(root, target string) *Layer {
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return nil
	}
	for i, layer := range c.Layers {
		for _, pattern := range layer.Paths {
			if matchGlob(pattern, rel) {
				return &c.Layers[i]
			}
		}
	}
	return nil
}

// matchGlob matches path against pattern with ** support.
// ** matches any number of path segments including zero.
func matchGlob(pattern, path string) bool {
	if !strings.Contains(pattern, "**") {
		matched, err := filepath.Match(pattern, path)
		return err == nil && matched
	}

	parts := strings.SplitN(pattern, "**", 2)
	prefix, suffix := parts[0], parts[1]

	path = filepath.ToSlash(path)
	prefix = filepath.ToSlash(prefix)
	suffix = filepath.ToSlash(suffix)

	if prefix != "" && !strings.HasPrefix(path, prefix) {
		return false
	}
	remaining := strings.TrimPrefix(path, prefix)

	suffix = strings.TrimPrefix(suffix, "/")
	if suffix == "" {
		return true
	}

	segments := strings.Split(remaining, "/")
	for i := range segments {
		tail := strings.Join(segments[i:], "/")
		if matchGlob(suffix, tail) {
			return true
		}
	}
	return false
}
