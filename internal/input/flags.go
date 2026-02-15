package input

import (
	"fmt"
	"netsnitch/internal/domain"
	"strings"
	"time"
)

type Flags map[string]string

func ExtractFlags(tokens []string) (Flags, []string, error) {
	flags := make(Flags)
	var rest []string

	for i := 0; i < len(tokens); i++ {
		t := tokens[i]

		if strings.HasPrefix(t, "--") {
			parts := strings.SplitN(t[2:], ":", 2)

			key := parts[0]

			if len(parts) == 2 {
				// --key:value
				flags[key] = parts[1]
				continue
			}

			// --key value
			if i+1 < len(tokens) && !strings.HasPrefix(tokens[i+1], "--") {
				flags[key] = tokens[i+1]
				i++
				continue
			}

			// --flag without value
			flags[key] = "true"
			continue
		}

		rest = append(rest, t)
	}

	return flags, rest, nil
}

func applyGlobalFlags(cfg *domain.Config, flags Flags) error {
	// Timeout
	if t, ok := flags["t"]; ok {
		d, err := time.ParseDuration(t)
		if err != nil {
			return fmt.Errorf("wrong time expression %s", t)
		}
		cfg.Timeout = d
	}

	//render rows/json
	if r, ok := flags["render"]; ok {

		r, err := domain.ParseRenderType(r)
		if err != nil {
			return fmt.Errorf("wrong render expression %s", r)
		}

	}

	return nil

}
