package input

import (
	"fmt"
	"netsnitch/internal/domain"
	"strings"
)

func Parse(args []string) (Query, error) {
	if len(args) == 0 {
		return Query{}, fmt.Errorf("empty input")
	}

	// split on stages
	rawStages := splitStages(args)

	var stages []Stage

	for _, tokens := range rawStages {
		if len(tokens) == 0 {
			return Query{}, fmt.Errorf("empty stage")
		}

		// first token = protocol (arp / tcp / ...)
		proto, err := ParseProtocol(tokens[0])
		if err != nil {
			return Query{}, err
		}
		////get unique parser for differente protocols like arp, tcp...
		////scans has own self register sturcts for Parser according to Parser interface
		parser, err := getParser(proto)
		if err != nil {
			return Query{}, err
		}
		////then use this parser to parse tokens
		////According to the protocol input blueprint
		////Parse method will return input.Stage struct with data
		stage, err := parser.Parse(tokens)
		if err != nil {
			return Query{}, err
		}

		stages = append(stages, stage)
	}

	return Query{Stages: stages}, nil
}

// // return with multiple commands && [[arp, 192.168.0.0/24], [tcp, 22,444,420, 192.168.0.3]]
// // return single commands [[arp, 192.168.0.0/24]]

func splitStages(args []string) [][]string {
	var stages [][]string
	var current []string

	for _, arg := range args {
		if arg == "&&" {
			stages = append(stages, current)
			current = nil
			continue
		}
		current = append(current, arg)
	}

	if len(current) > 0 {
		stages = append(stages, current)
	}

	return stages
}

func ParseProtocol(s string) (domain.Protocol, error) {
	switch strings.ToLower(s) {
	case "tcp":
		return domain.TCP, nil
	case "arp":
		return domain.ARP_ACTIVE, nil
	default:
		return "", fmt.Errorf("unknown protocol: %s", s)
	}
}
