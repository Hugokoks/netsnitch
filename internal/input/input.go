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

	// rozdělíme na stage podle &&
	rawStages := splitStages(args)

	var stages []Stage

	for _, tokens := range rawStages {
		if len(tokens) == 0 {
			return Query{}, fmt.Errorf("empty stage")
		}

		// první token = protokol (arp / tcp / ...)
		proto, err := ParseProtocol(tokens[0])
		if err != nil {
			return Query{}, err
		}

		parser, err := getParser(proto)
		if err != nil {
			return Query{}, err
		}

		stage, err := parser.Parse(tokens)
		if err != nil {
			return Query{}, err
		}

		stages = append(stages, stage)
	}

	return Query{Stages: stages}, nil
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
