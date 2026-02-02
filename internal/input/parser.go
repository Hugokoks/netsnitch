package input

import "netsnitch/internal/domain"

type Parser interface {
	Protocol() domain.Protocol
	Parse(tokens []string) (Stage, error)
}
