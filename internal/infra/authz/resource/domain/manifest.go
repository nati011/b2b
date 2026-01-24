package domain

import (
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

// ManifestFile represents the structure of the external manifest file.
type ManifestFile struct {
	Resources []ManifestResourceEntry `yaml:"resources"`
}

// ManifestResourceEntry represents a resource entry in the YAML file.
type ManifestResourceEntry struct {
	Code        string                `yaml:"code"`
	Service     string                `yaml:"service"`
	Description string                `yaml:"description"`
	Actions     []ManifestActionEntry `yaml:"actions"`
}

// ManifestActionEntry represents an action entry in the YAML file.
type ManifestActionEntry struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// LoadResources loads the resource manifest from the path specified in the config.
// This follows the same pattern as Load() for consistency.
func LoadResources(manifestPath string) ([]ManifestResource, error) {
	if manifestPath == "" {
		// Default to resources.yaml in the same directory as config
		manifestPath = "config/resources.yaml"
	}

	manifest, err := LoadManifestFromFile(manifestPath)
	if err != nil {
		slog.Error("Error loading resource manifest", "err", err, "path", manifestPath)
		return nil, fmt.Errorf("failed to load resource manifest: %w", err)
	}
	return manifest, nil
}

// LoadManifestFromFile loads the resource manifest from a YAML file.
func LoadManifestFromFile(filePath string) ([]ManifestResource, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest file %s: %w", filePath, err)
	}

	var manifestFile ManifestFile
	if err := yaml.Unmarshal(data, &manifestFile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal manifest YAML: %w", err)
	}

	return convertToDomainManifest(manifestFile.Resources), nil
}

// convertToDomainManifest converts file entries to domain manifest resources.
func convertToDomainManifest(entries []ManifestResourceEntry) []ManifestResource {
	result := make([]ManifestResource, len(entries))
	for i, entry := range entries {
		actions := make([]ManifestAction, len(entry.Actions))
		for j, action := range entry.Actions {
			actions[j] = ManifestAction{
				Name:        action.Name,
				Description: action.Description,
			}
		}
		result[i] = ManifestResource{
			Code:        entry.Code,
			Service:     entry.Service,
			Description: entry.Description,
			Actions:     actions,
		}
	}
	return result
}
