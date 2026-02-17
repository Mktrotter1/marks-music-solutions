package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds all server configuration.
type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Music   MusicConfig   `yaml:"music"`
	Database DatabaseConfig `yaml:"database"`
	Transcode TranscodeConfig `yaml:"transcode"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// MusicConfig holds music library settings.
type MusicConfig struct {
	Directories []string `yaml:"directories"`
	WatchForChanges bool `yaml:"watch_for_changes"`
}

// DatabaseConfig holds database settings.
type DatabaseConfig struct {
	Path string `yaml:"path"`
}

// TranscodeConfig holds transcoding settings.
type TranscodeConfig struct {
	CacheDir   string `yaml:"cache_dir"`
	FFmpegPath string `yaml:"ffmpeg_path"`
}

// Addr returns the listen address string.
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// Load reads configuration from a YAML file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Database: DatabaseConfig{
			Path: "data/mms.db",
		},
		Transcode: TranscodeConfig{
			CacheDir:   "data/cache/transcode",
			FFmpegPath: "ffmpeg",
		},
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if len(cfg.Music.Directories) == 0 {
		return nil, fmt.Errorf("at least one music directory must be configured")
	}

	return cfg, nil
}

// LoadOrDefault tries to load config from path, falls back to env/defaults.
func LoadOrDefault(path string) (*Config, error) {
	if path != "" {
		return Load(path)
	}

	// Fall back to environment variable
	if envPath := os.Getenv("MMS_CONFIG"); envPath != "" {
		return Load(envPath)
	}

	// Try default locations
	for _, p := range []string{"config.yaml", "config.yml"} {
		if _, err := os.Stat(p); err == nil {
			return Load(p)
		}
	}

	return nil, fmt.Errorf("no config file found (tried config.yaml, config.yml, MMS_CONFIG env)")
}
