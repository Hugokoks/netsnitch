package input

import (
	"fmt"
	"netsnitch/internal/domain"
	"strings"
	"time"
)

type FlagSpec struct {
	HasValue bool
	Default  string
	Usage    string
}

var FlagRegistry = map[string]FlagSpec{
	"p": {
		HasValue: true,
		Usage:    "Port range or list (e.g. 1-100 or 80,443)",
	},
	"mode": {
		HasValue: true,
		Default:  "full",
		Usage:    "Scan mode (full | stealth)",
	},
	"render": {
		HasValue: true,
		Default:  "rows",
		Usage:    "Output format (rows | json)",
	},
	"t": {
		HasValue: true,
		Default:  "1s",
		Usage:    "Timeout (e.g. 500ms, 2s)",
	},
	"open": {
		HasValue: false,
		Usage:    "Show only open ports",
	},
	"help": {
		HasValue: false,
		Usage:    "Show help message",
	},
}

type Flags map[string]string

func ExtractFlags(tokens []string) (Flags, []string, error) {
	flags := make(Flags)
	var rest []string

	for i := 0; i < len(tokens); i++ {
		t := tokens[i]

		if strings.HasPrefix(t, "--") {

			key := t[2:]

			spec, exists := FlagRegistry[key]
			if !exists {
				return nil, nil, fmt.Errorf("unknown flag --%s", key)
			}

			// Flag expects value
			if spec.HasValue {

				if i+1 >= len(tokens) {
					return nil, nil, fmt.Errorf("flag --%s requires a value", key)
				}

				if strings.HasPrefix(tokens[i+1], "--") {
					return nil, nil, fmt.Errorf("flag --%s requires a value", key)
				}

				flags[key] = tokens[i+1]
				i++
				continue
			}

			// Boolean flag
			flags[key] = "true"
			continue
		}

		rest = append(rest, t)
	}

	return flags, rest, nil
}

func applyGlobalFlags(cfg *domain.Config, flags Flags) error {

	// Timeout from CLI
	if val, ok := flags["t"]; ok {
		dur, err := time.ParseDuration(val)
		if err != nil {
			return fmt.Errorf("invalid timeout value %s", val)
		}
		cfg.Timeout = dur
	}

	// Render type from CLI
	if val, ok := flags["render"]; ok {
		render, err := domain.ParseRenderType(val)
		if err != nil {
			return fmt.Errorf("invalid render expression %s", val)
		}
		cfg.Render = render
	}

	return nil
}
