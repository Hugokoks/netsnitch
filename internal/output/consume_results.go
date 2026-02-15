package output

import (
	"context"
	"netsnitch/internal/domain"
)

func ConsumeResults(ctx context.Context, results <-chan domain.Result) {
	for {
		select {
		case <-ctx.Done():
			return

		case res, ok := <-results:
			if !ok {
				return
			}

			if f, ok := formatters[res.Protocol]; ok {

				f.Format(res)
			}

		}
	}
}
