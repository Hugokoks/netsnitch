package fingerprint_test

import (
	"netsnitch/internal/fingerprint"
	"testing"
)

func TestEngineDetect(t *testing.T) {

	fp := fingerprint.NewEngine()

	if err := fp.LoadRules("../../data/rules.json"); err != nil {

		t.Fatalf("load rules %s", err)
	}
	if err := fp.LoadProbes("../../data/probes.json"); err != nil {

		t.Fatalf("load probes %s", err)
	}

	tests := []struct {
		name     string
		port     int
		banner   string
		wantServ string
	}{
		{"Valid SSH", 22, "SSH-2.0-OpenSSH_8.2p1", "ssh"},
		{"Invalid banner", 22, "HTTP/1.1 200 OK", ""},
		{"Empty banner", 80, "", ""},
		{"Valid HTTP", 80, "HTTP/1.1 200 OK", "http"},
		{"Valid HTTP without Port", 6000, "HTTP/1.1 200 OK", "http"},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			got := fp.Detect(tt.port, tt.banner)

			// Kontrola, jestli jsme něco dostali, když jsme to čekali
			if tt.wantServ != "" && got == nil {
				t.Errorf("%s: Detect() return nil, expected %s", tt.name, tt.wantServ)
				return
			}

			// Pokud chceme vidět detaily, když se služba neshoduje
			if got != nil && got.Service != tt.wantServ {
				// %+v vypíše strukturu i s názvy polí (Service, Product, Version...)
				t.Errorf("%s: wrong detect!\n got: %+v\n expected: %s",
					tt.name, got, tt.wantServ)
			}

			if got != nil {

				t.Logf("Right result %s: Service=%s, Product=%s, Version=%s",
					tt.name, got.Service, got.Product, got.Version)
			}
		})
	}
}
