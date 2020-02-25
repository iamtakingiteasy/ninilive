// Package config for all configuration needs
package config

// Loader for configuration
type Loader interface {
	Load() (Values, error)
}
