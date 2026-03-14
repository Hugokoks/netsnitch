package fingerprint

import (
	"context"
	"fmt"
	"net"
	"time"
)

func grabWithContext(
	ctx context.Context,
	conn net.Conn,
	timeout time.Duration,
) string {
	buf := make([]byte, 4096)

	_ = conn.SetReadDeadline(time.Now().Add(timeout))

	type result struct {
		n   int
		err error
	}

	ch := make(chan result, 1)

	go func() {
		n, err := conn.Read(buf)
		ch <- result{n, err}
	}()

	select {
	case <-ctx.Done():
		_ = conn.Close()
		return ""

	case res := <-ch:
		if res.err != nil || res.n == 0 {
			fmt.Printf("[grab] read err=%v n=%d\n", res.err, res.n)
			return ""
		}

		return string(buf[:res.n])
	}
}
