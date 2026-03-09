package input

import (
	"fmt"
	"netsnitch/internal/domain"
	"strings"
	"time"
)

// ExtractFlags separates flags from positional arguments (rest)
func ExtractFlags(tokens []string) (Flags, []string, error) {
	flags := make(Flags)
	var rest []string

	for i := 0; i < len(tokens); i++ {
		t := tokens[i]

		// Check if token starts with at least one dash
		if strings.HasPrefix(t, "-") {
			// Handle both -flag and --flag by trimming all leading dashes
			key := strings.TrimLeft(t, "-")

			spec, exists := FlagRegistry[key]
			if !exists {
				return nil, nil, fmt.Errorf("unknown flag -%s", key)
			}

			// Flag requires an associated value
			if spec.HasValue {
				// Check if there is a next token and it's not another flag
				if i+1 >= len(tokens) || strings.HasPrefix(tokens[i+1], "-") {
					return nil, nil, fmt.Errorf("flag -%s requires a value", key)
				}

				flags[key] = tokens[i+1]
				i++ // Skip the value token in the next iteration
				continue
			}

			// Boolean (switch) flag
			flags[key] = "true"
			continue
		}

		// Positional argument (protocol, IP, etc.)
		rest = append(rest, t)
	}

	return flags, rest, nil
}

// applyGlobalFlags maps extracted flags to the domain Config object
func applyGlobalFlags(cfg *domain.Config, flags Flags) error {
	// Parse timeout duration if provided
	if val, ok := flags["t"]; ok {
		dur, err := time.ParseDuration(val)
		if err != nil {
			return fmt.Errorf("invalid timeout value %s", val)
		}
		cfg.Timeout = dur
	}

	// Parse output render type if provided
	if val, ok := flags["r"]; ok {
		render, err := domain.ParseRenderType(val)
		if err != nil {
			return fmt.Errorf("invalid render expression %s", val)
		}
		cfg.Render = render
	}

	// Set OpenOnly filter if flag is present
	if _, ok := flags["o"]; ok {
		cfg.OpenOnly = true
	}

	return nil
}
