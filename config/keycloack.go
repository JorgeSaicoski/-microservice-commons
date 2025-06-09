// config/keycloak.go
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/JorgeSaicoski/microservice-commons/utils"
)

// KeycloakConfig holds Keycloak authentication configuration
type KeycloakConfig struct {
	URL                string
	Realm              string
	PublicKeyBase64    string
	RequiredClaims     []string
	SkipPaths          []string
	KeyRefreshInterval time.Duration
	HTTPTimeout        time.Duration
}

// LoadKeycloakConfig loads Keycloak configuration from environment
func LoadKeycloakConfig() KeycloakConfig {
	return KeycloakConfig{
		URL:                utils.GetEnv("KEYCLOAK_URL", ""),
		Realm:              utils.GetEnv("KEYCLOAK_REALM", "master"),
		PublicKeyBase64:    utils.GetEnv("KEYCLOAK_PUBLIC_KEY", ""),
		RequiredClaims:     parseStringSlice(utils.GetEnv("KEYCLOAK_REQUIRED_CLAIMS", "sub,preferred_username")),
		SkipPaths:          parseStringSlice(utils.GetEnv("KEYCLOAK_SKIP_PATHS", "/health,/metrics")),
		KeyRefreshInterval: parseDuration(utils.GetEnv("KEYCLOAK_KEY_REFRESH_INTERVAL", "1h")),
		HTTPTimeout:        parseDuration(utils.GetEnv("KEYCLOAK_HTTP_TIMEOUT", "10s")),
	}
}

// Validate validates the Keycloak configuration
func (kc *KeycloakConfig) Validate() error {
	// Must have either static key or JWKS endpoint
	hasStaticKey := kc.PublicKeyBase64 != ""
	hasJWKS := kc.URL != "" && kc.Realm != ""

	if !hasStaticKey && !hasJWKS {
		return fmt.Errorf("keycloak config must provide either PublicKeyBase64 or both URL and Realm")
	}

	if kc.KeyRefreshInterval <= 0 {
		return fmt.Errorf("key refresh interval must be positive")
	}

	if kc.HTTPTimeout <= 0 {
		return fmt.Errorf("HTTP timeout must be positive")
	}

	return nil
}

// HasStaticKey returns true if using static public key
func (kc *KeycloakConfig) HasStaticKey() bool {
	return kc.PublicKeyBase64 != ""
}

// HasJWKS returns true if using JWKS endpoint
func (kc *KeycloakConfig) HasJWKS() bool {
	return kc.URL != "" && kc.Realm != ""
}

// GetJWKSURL returns the JWKS endpoint URL
func (kc *KeycloakConfig) GetJWKSURL() string {
	if !kc.HasJWKS() {
		return ""
	}
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", kc.URL, kc.Realm)
}

// ShouldSkipPath returns true if the path should skip authentication
func (kc *KeycloakConfig) ShouldSkipPath(path string) bool {
	for _, skipPath := range kc.SkipPaths {
		if path == skipPath || strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// parseStringSlice parses a comma-separated string into a slice
func parseStringSlice(value string) []string {
	if value == "" {
		return []string{}
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// parseDuration parses a duration string with fallback
func parseDuration(value string) time.Duration {
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	return 1 * time.Hour // Default fallback
}
