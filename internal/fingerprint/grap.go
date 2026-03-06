package fingerprint

import (
	"context"
	"net"
	"time"
)

func grabWithContext(
	ctx context.Context,
	conn net.Conn,
	timeout time.Duration,
) string {

	buf := make([]byte, 4096)

	conn.SetReadDeadline(time.Now().Add(timeout))

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
		conn.Close()
		return ""

	case res := <-ch:

		if res.err != nil || res.n == 0 {
			return ""
		}

		return string(buf[:res.n])
	}
}
