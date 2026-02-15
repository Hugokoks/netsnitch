package output_test

import (
	"context"
	"fmt"
	"net"
	"netsnitch/internal/domain"
	"netsnitch/internal/output"
	_ "netsnitch/internal/scans/arp"
	_ "netsnitch/internal/scans/tcp"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"
)

var Wg sync.WaitGroup

func TestOutput(t *testing.T) {

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()
	resChan := make(chan domain.Result)

	resData := createResults()

	consumer := output.NewConsumer(ctx, resChan)

	go consumer.Start()

	Wg.Add(1)

	go func() {

		sendResults(ctx, resChan, resData)
		close(resChan)

	}()

	Wg.Wait()
	consumer.Wait()

}
func sendResults(ctx context.Context, resChan chan<- domain.Result, resData []domain.Result) {

	defer Wg.Done()

	for _, res := range resData {
		time.Sleep(200 * time.Microsecond)
		select {
		case <-ctx.Done():
			return
		case resChan <- res:

		}

	}

}

func createResults() []domain.Result {
	return []domain.Result{
		// 1. TCP nález - Otevřený port s bannerem
		{
			IP:         net.ParseIP("192.168.1.15"),
			Port:       80,
			Protocol:   domain.TCP,
			RenderType: domain.ROWS_OUT,
			Open:       true,
			Alive:      true,
			RTT:        12 * time.Millisecond,
			Service:    "http",
			Banner:     "Apache/2.4.41 (Ubuntu)",
		},

		// 2. TCP nález - JSON výstup (pro integraci s jinými nástroji)
		{
			IP:         net.ParseIP("10.0.0.1"),
			Port:       22,
			Protocol:   domain.TCP,
			RenderType: domain.JSON_OUT,
			Open:       true,
			Alive:      true,
			RTT:        45 * time.Millisecond,
			Service:    "ssh",
			Banner:     "OpenSSH_8.2p1",
		},

		// 3. ARP nález - Lokální zařízení (MAC adresa)
		{
			IP:         net.ParseIP("192.168.1.1"),
			MAC:        net.HardwareAddr{0x00, 0x50, 0x56, 0xC0, 0x00, 0x08},
			Protocol:   domain.ARP,
			RenderType: domain.ROWS_OUT,
			Alive:      true,
			RTT:        500 * time.Microsecond,
		},
		{
			IP:         net.ParseIP("192.168.1.1"),
			MAC:        net.HardwareAddr{0x00, 0x50, 0x56, 0xC0, 0x00, 0x08},
			Protocol:   domain.ARP,
			RenderType: domain.JSON_OUT,
			Alive:      true,
			RTT:        500 * time.Microsecond,
		},
		// 4. ICMP - Host je naživu (Ping)
		{
			IP:         net.ParseIP("8.8.8.8"),
			Protocol:   domain.ICMP,
			RenderType: domain.ROWS_OUT,
			Alive:      true,
			RTT:        22 * time.Millisecond,
		},

		// 5. Selhání - Port je zavřený nebo nastala chyba
		{
			IP:         net.ParseIP("192.168.1.50"),
			Port:       443,
			Protocol:   domain.TCP,
			RenderType: domain.ROWS_OUT,
			Open:       true,
			Alive:      true,
			Error:      fmt.Errorf("connection refused"),
		},
	}
}
