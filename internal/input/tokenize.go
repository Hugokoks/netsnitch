package input

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
