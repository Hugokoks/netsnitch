package output

import (
	"fmt"
	"netsnitch/internal/domain"
)

// Formatter defines how a scan result should be displayed.
type Formatter interface {
	Protocol() domain.Protocol
	FormatRows(res domain.Result) string
	FormatJson(res domain.Result) string
}

// internal registry of formatters by protocol

var formatters = make(map[domain.Protocol]Formatter)

// Register registers a formatter for a protocol.
// Called automatically via init() in formatter implementations.s
func Register(f Formatter) {

	formatters[f.Protocol()] = f

}

func out(res domain.Result) error {

	renderer, ok := renderers[res.RenderType]

	if !ok {
		return fmt.Errorf("undifined RenderType %s", res.RenderType)
	}

	f, ok := formatters[res.Protocol]

	if !ok {
		return fmt.Errorf("undifined formatter protocol %s", res.Protocol)

	}
	output := renderer(f, res)

	////TODO: add file output
	fmt.Println(output)

	return nil

}
