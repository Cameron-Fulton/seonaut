package config

import (
	"path/filepath"

	"github.com/spf13/viper"
)

// CrawlerConfig stores the configuration for the crawler.
type CrawlerConfig struct {
	Agent string `mapstructure:"agent"`
}

// HTTPServerConfig stores the configuration for the HTTP server.
type HTTPServerConfig struct {
	Server string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	URL    string `mapstructure:"url"`
}

// DBConfig stores the configuration for the database store.
type DBConfig struct {
	Server string `mapstructure:"server"`
	Port   int    `mapstructure:"port"`
	User   string `mapstructure:"user"`
	Pass   string `mapstructure:"password"`
	Name   string `mapstructure:"database"`
}

// AIProviderConfig stores configuration for AI providers.
type AIProviderConfig struct {
	ClaudeAPIKey  string `mapstructure:"claude_api_key"`
	ClaudeModel   string `mapstructure:"claude_model"`
	LMStudioURL   string `mapstructure:"lmstudio_url"`
	LMStudioModel string `mapstructure:"lmstudio_model"`
	Enabled       bool   `mapstructure:"enabled"`
}

// IntegrationsProviderConfig stores configuration for external API integrations.
type IntegrationsProviderConfig struct {
	PageSpeedAPIKey string `mapstructure:"pagespeed_api_key"`
	AhrefsAPIKey    string `mapstructure:"ahrefs_api_key"`
}

// Config stores the configuration for the application.
type Config struct {
	Crawler      *CrawlerConfig              `mapstructure:"crawler"`
	HTTPServer   *HTTPServerConfig           `mapstructure:"server"`
	DB           *DBConfig                   `mapstructure:"database"`
	AI           *AIProviderConfig           `mapstructure:"ai"`
	Integrations *IntegrationsProviderConfig `mapstructure:"integrations"`
}

// NewConfig loads the configuration from the specified file and path.
func NewConfig(configFile string) (*Config, error) {
	viper.AddConfigPath(filepath.Dir(configFile))
	viper.SetConfigName(filepath.Base(configFile))
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
