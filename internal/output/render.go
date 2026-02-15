package output

import "netsnitch/internal/domain"

type RenderFunc func(t Formatter, res domain.Result) string

var renderers = map[domain.RenderType]RenderFunc{
	domain.JSON_OUT: renderJson,
	domain.ROWS_OUT: renderRows,
}

func renderJson(f Formatter, res domain.Result) string {

	return f.FormatJson(res)

}
func renderRows(f Formatter, res domain.Result) string {

	return f.FormatRows(res)

}
