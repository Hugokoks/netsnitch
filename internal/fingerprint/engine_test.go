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
		name        string
		port        int
		banner      string
		wantServ    string
		wantProduct string
		wantVersion string
	}{
		{
			"Valid SSH",
			22,
			"SSH-2.0-OpenSSH_8.2p1",
			"ssh",
			"OpenSSH",
			"8.2p1",
		},
		{
			"HTTP on SSH port uses generic match",
			22,
			"HTTP/1.1 200 OK",
			"http",
			"",
			"1.1",
		},
		{
			"Empty banner",
			80,
			"",
			"",
			"",
			"",
		},
		{
			"Valid HTTP",
			80,
			"HTTP/1.1 200 OK",
			"http",
			"",
			"",
		},
		{
			"HTTP Apache",
			80,
			"HTTP/1.1 200 OK\r\nServer: Apache/2.4.58\r\n\r\n",
			"http",
			"Apache",
			"2.4.58",
		},
		{
			"HTTP Generic no port",
			6000,
			"HTTP/1.1 200 OK",
			"http",
			"",
			"1.1",
		},
		{
			"FTP Banner",
			21,
			"220 (vsFTPd 3.0.3)",
			"ftp",
			"vsFTPd",
			"3.0.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fp.Detect(tt.port, tt.banner)

			if tt.wantServ == "" {
				if got != nil {
					t.Errorf("%s: expected nil but got %+v", tt.name, got)
				}
				return
			}

			if got == nil {
				t.Fatalf("%s: Detect() returned nil, expected %s", tt.name, tt.wantServ)
			}

			if got.Service != tt.wantServ {
				t.Errorf("%s: wrong service\n got: %s\n want: %s",
					tt.name, got.Service, tt.wantServ)
			}

			if got.Product != tt.wantProduct {
				t.Errorf("%s: wrong product\n got: %s\n want: %s",
					tt.name, got.Product, tt.wantProduct)
			}

			if got.Version != tt.wantVersion {
				t.Errorf("%s: wrong version\n got: %s\n want: %s",
					tt.name, got.Version, tt.wantVersion)
			}

			t.Logf("OK %s -> Service=%s Product=%s Version=%s",
				tt.name, got.Service, got.Product, got.Version)
		})
	}
}
