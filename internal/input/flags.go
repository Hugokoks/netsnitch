package input

import "strings"

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
