package config

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App        AppConfig       `yaml:"app"`
	HttpServer ServerConfig    `yaml:"server"` // use server to match YAML
	DB         DBConfig        `yaml:"db"`
	Logging    LoggingConfig   `yaml:"logging"`
	Auth       AuthConfig      `yaml:"auth"`
	Resources  ResourcesConfig `yaml:"resources"`
	Roles      *RolesConfig    `yaml:"roles,omitempty"` // Optional: path to roles.yaml file
}

type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Env     string `yaml:"env"`
}

// ServerConfig conatains http server settings
type ServerConfig struct {
	Host         string        `yaml:"host"`
	HttpPort     int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// Address returns the full server address in the format "host:port"
func (c ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.HttpPort)
}

type DBConfig struct {
	Driver   string `yaml:"driver"` // postgres, mysql, sqlite, etc
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"` // postgres-specific
}

// DNS returns a connection string
func (c DBConfig) DSN() string {
	if c.Driver != "postgres" {
		return ""
	}

	if c.SSLMode != "" {
		return fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s",
			c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode,
		)
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		c.User, c.Password, c.Host, c.Port, c.Database,
	)
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level          string  `yaml:"level"`           // debug, info, warn, error
	Format         string  `yaml:"format"`          // text, json
	SourceLocation bool    `yaml:"source_location"` // include file:line in logs
	SamplingRate   float64 `yaml:"sampling_rate"`   // 0.0-1.0, 1.0 = log all, 0.5 = log 50%
	Async          bool    `yaml:"async"`           // use async buffered logging
	Sanitize       bool    `yaml:"sanitize"`        // enable sensitive data sanitization
}

// AuthConfig contains auth mode settings, and auth-mode specific settings
type AuthConfig struct {
	Mode   string           `yaml:"mode"` // basic | oauth2
	Basic  *BasicAuthConfig `yaml:"basic,omitempty"`
	OAuth2 *OAuth2Config    `yaml:"oauth2,omitempty"`
}

// BasicAuthConfig holds static credentials for Basic authentication.
type BasicAuthConfig struct {
	Realm       string          `yaml:"realm"`
	EmailDomain string          `yaml:"email_domain"` // Default email domain for superadmin users (e.g., "@admin.local")
	UserPrefix  string          `yaml:"user_prefix"`  // Default prefix for superadmin usernames (e.g., "admin-")
	UserName    string          `yaml:"user_name"`    // Default display name for superadmin users (e.g., "Super Administrator")
	Users       []BasicAuthUser `yaml:"users"`
}

// BasicAuthUser represents a configured Basic auth principal.
type BasicAuthUser struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// OAuth2Config contains auth server configs
type OAuth2Config struct {
	ClientID     string   `yaml:"client_id"`
	ClientSecret string   `yaml:"client_secret"`
	AuthorizeURL string   `yaml:"authorize_url"`
	TokenURL     string   `yaml:"token_url"`
	RedirectURL  string   `yaml:"redirect_url"`
	Scopes       []string `yaml:"scopes"`
}

// ResourcesConfig contains resource manifest configuration
type ResourcesConfig struct {
	Path string `yaml:"path"` // Path to resources.yaml file
}

// RolesConfig contains role definitions configuration
type RolesConfig struct {
	Path string `yaml:"path"` // Path to roles.yaml file
}

// RoleDefinition defines a role with its permission patterns
type RoleDefinition struct {
	ID          string              `yaml:"id"`                    // Unique role identifier
	Name        string              `yaml:"name"`                  // Display name
	Description string              `yaml:"description"`           // Role description
	Permissions []PermissionPattern `yaml:"permissions,omitempty"` // Permission patterns for this role
}

// PermissionPattern defines a pattern for matching permissions
type PermissionPattern struct {
	Resource string   `yaml:"resource"` // Resource code (e.g., "users", "roles")
	Actions  []string `yaml:"actions"`  // List of action names (e.g., ["view", "create"])
}

// Load loads all configs from the specified config file path.
func Load(configPath string) (*Config, error) {
	cfg, err := loadConfig(configPath)
	if err != nil {
		slog.Error("Error loading configuration", "err", err, "path", configPath)
		panic(err)
	}

	return cfg, nil
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	var cfg Config

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &cfg, nil
}

// RolesFile represents the structure of the roles YAML file
type RolesFile struct {
	DefaultRoles []RoleDefinition `yaml:"default_roles"`
}

// LoadRolesFromFile loads role definitions from a YAML file
func LoadRolesFromFile(filePath string) ([]RoleDefinition, error) {
	if filePath == "" {
		// Default to roles.yaml in the same directory as config
		filePath = "config/roles.yaml"
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read roles file %s: %w", filePath, err)
	}

	var rolesFile RolesFile
	if err := yaml.Unmarshal(data, &rolesFile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal roles YAML: %w", err)
	}

	return rolesFile.DefaultRoles, nil
}
