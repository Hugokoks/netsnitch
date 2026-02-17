package input

import (
	"fmt"
	"netsnitch/internal/domain"
)

type Query struct {
	Configs []domain.Config
}

func Parse(args []string) (Query, error) {
	if len(args) == 0 {
		return Query{}, fmt.Errorf("empty input")
	}

	// split on stages
	rawStages := splitStages(args)

	var configs []domain.Config

	for _, tokens := range rawStages {

		flags, rest, err := ExtractFlags(tokens)

		if err != nil {
			return Query{}, err
		}

		if len(rest) == 0 {
			return Query{}, fmt.Errorf("missing protocol")
		}

		////Create config object
		config := domain.NewDefaultConfig()

		// first token = protocol (arp / tcp / ...)
		proto, err := domain.ParseProtocol(rest[0])
		if err != nil {
			return Query{}, err
		}
		////get unique parser for differente protocols like arp, tcp...
		////scans has own self register sturcts for Parser according to Parser interface
		parser, err := getParser(proto)

		if err != nil {
			return Query{}, err
		}

		if err := parser.Parse(&config, rest, flags); err != nil {
			return Query{}, err
		}

		if err = applyGlobalFlags(&config, flags); err != nil {
			return Query{}, err
		}

		////Set defualt values of empty parameters
		parser.ApplyDefaults(&config)
		configs = append(configs, config)
	}

	return Query{Configs: configs}, nil
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
