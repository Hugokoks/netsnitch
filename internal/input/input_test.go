package input_test

import (
	"netsnitch/internal/input"
	_ "netsnitch/internal/scans/tcp"
	"testing"
)

func TestInput(t *testing.T) {

	args := []string{
		//"tcp", "-p", "66,222,122", "-open", "192.168.1.1",
		"tcp", "-h",
		//"-h",
	}

	query, err := input.Parse(args)
	if err != nil {

		t.Fatalf("Parse failed %v", err)

	}

	t.Logf("configs %v", query.Configs)

}
