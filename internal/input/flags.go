package input

import (
	"fmt"
	"netsnitch/internal/domain"
	"strings"
	"time"
)

type Flags map[string]string

// flagsWithValue defines which flags require a value.
// Boolean-only flags (like --open) are not listed here.
var flagsWithValue = map[string]bool{
	"p":      true,
	"mode":   true,
	"render": true,
	"t":      true,
}

func ExtractFlags(tokens []string) (Flags, []string, error) {
	flags := make(Flags)
	var rest []string

	for i := 0; i < len(tokens); i++ {
		t := tokens[i]

		// Check if token is a long flag (--something)
		if strings.HasPrefix(t, "--") {

			key := t[2:]

			// If flag expects a value, consume the next token
			if flagsWithValue[key] {

				if i+1 >= len(tokens) {
					return nil, nil, fmt.Errorf("flag --%s requires a value", key)
				}

				// Prevent another flag from being used as value
				if strings.HasPrefix(tokens[i+1], "--") {
					return nil, nil, fmt.Errorf("flag --%s requires a value", key)
				}

				flags[key] = tokens[i+1]
				i++ // Skip value token
				continue
			}

			// Boolean flag (no value expected)
			flags[key] = "true"
			continue
		}

		// Non-flag tokens (protocol, target, etc.)
		rest = append(rest, t)
	}

	return flags, rest, nil
}

func applyGlobalFlags(cfg *domain.Config, flags Flags) error {

	// Parse timeout flag
	if str, ok := flags["t"]; ok {
		dur, err := time.ParseDuration(str)
		if err != nil {
			return fmt.Errorf("invalid time expression %s", str)
		}
		cfg.Timeout = dur
	}

	// Parse render type (rows/json)
	if str, ok := flags["render"]; ok {
		render, err := domain.ParseRenderType(str)
		if err != nil {
			return fmt.Errorf("invalid render expression %s", str)
		}
		cfg.Render = render
	}

	return nil
}
