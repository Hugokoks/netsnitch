package fingerprint

import (
	"context"
	"fmt"
	"net"
	"time"
)

func Identify(
	ctx context.Context,
	engine *Engine,
	firstConn net.Conn,
	ip net.IP,
	port int,
	timeout time.Duration,
) *ServiceInfo {

	raw := grabWithContext(ctx, firstConn, timeout)

	firstConn.Close()

	if raw != "" {

		return engine.Detect(raw)
	}

	probes := GetProbesForPort(port)

	address := fmt.Sprintf("%s:%d", ip, port)

	for _, p := range probes {

		select {
		case <-ctx.Done():
			return nil
		default:
		}

		d := net.Dialer{Timeout: timeout}

		conn, err := d.DialContext(ctx, "tcp", address)

		if err != nil {
			continue
		}

		conn.SetWriteDeadline(time.Now().Add(timeout))

		_, err = conn.Write(p.RawData)

		if err != nil {

			conn.Close()
			continue
		}

		response := grabWithContext(ctx, conn, timeout)

		conn.Close()

		if response == "" {
			continue
		}

		info := engine.Detect(response)

		if info != nil {
			return info
		}
	}

	return nil
}

func grabWithContext(
	ctx context.Context,
	conn net.Conn,
	timeout time.Duration,
) string {

	buf := make([]byte, 4096)

	conn.SetReadDeadline(time.Now().Add(timeout))

	type readResult struct {
		n   int
		err error
	}

	resChan := make(chan readResult, 1)

	go func() {
		n, err := conn.Read(buf)
		resChan <- readResult{n: n, err: err}
	}()

	select {

	case <-ctx.Done():
		// kill blocking read
		conn.Close()

		return ""

	case res := <-resChan:

		if res.err != nil || res.n == 0 {
			return ""
		}

		return string(buf[:res.n])
	}
}
