package input

import (
	"fmt"
	"netsnitch/internal/domain"
	"sort"
	"strings"
)

type Query struct {
	Configs []domain.Config
}

var ErrHelpRequested = fmt.Errorf("help requested")

type FlagSpec struct {
	HasValue    bool
	Default     string
	Usage       string
	Description string
}

type Flags map[string]string

var FlagRegistry = map[string]FlagSpec{
	"p": {
		HasValue:    true,
		Usage:       "-p (1-100 | 80,443)",
		Description: "Target ports: supports ranges, lists, or single ports.",
	},
	"mode": {
		HasValue:    true,
		Default:     "full",
		Usage:       "-mode (f | s)",
		Description: "Scan strategy: 'f' for full handshake, 's' for stealth SYN.",
	},
	"render": {
		HasValue:    true,
		Default:     "rows",
		Usage:       "-render (rows | json)",
		Description: "Output format: table rows or machine-readable JSON.",
	},
	"t": {
		HasValue:    true,
		Default:     "1s",
		Usage:       "-t (1s | 500ms)",
		Description: "Network timeout: higher value means better accuracy on slow links.",
	},
	"open": {
		HasValue:    false,
		Usage:       "-open",
		Description: "Show only active/open ports; hide everything else.",
	},
	"h": {
		HasValue:    false,
		Usage:       "-h",
		Description: "Show help and exit.",
	},
}

type ScanSpec struct {
	Usage       string
	Description string
	Example     string
	UsableFlags map[string]FlagSpec
}

var ScanRegistry = map[string]ScanSpec{
	"tcp": {
		Usage:       "tcp [flags] <target>",
		Description: "Performs deep inspection of TCP ports to identify active services and potential entry points.",
		Example:     "tcp -p 1-1024 -mode s -open 192.168.1.1",
		UsableFlags: map[string]FlagSpec{
			"p":      FlagRegistry["p"],
			"mode":   FlagRegistry["mode"],
			"t":      FlagRegistry["t"],
			"open":   FlagRegistry["open"],
			"render": FlagRegistry["render"],
		},
	},
	"arp": {
		Usage:       "arp [flags] <target>",
		Description: "Maps the local network by resolving IP addresses to MAC addresses using ARP requests. Ideal for fast host discovery.",
		Example:     "arp -t 500ms 192.168.1.0/24",
		UsableFlags: map[string]FlagSpec{
			"t":      FlagRegistry["t"],
			"render": FlagRegistry["render"],
		},
	},
	"udp": {
		Usage:       "udp [flags] <target>",
		Description: "Scans for open UDP ports. Note: UDP is connectionless and results may be less reliable than TCP.",
		Example:     "udp -p 53,161 192.168.1.1",
		UsableFlags: map[string]FlagSpec{
			"p":      FlagRegistry["p"],
			"t":      FlagRegistry["t"],
			"render": FlagRegistry["render"],
		},
	},
}

func GetUsage(proto domain.Protocol) string {
	spec, ok := ScanRegistry[string(proto)]
	if !ok {
		return fmt.Sprintf("Usage: %s [flags] <target>", proto)
	}

	// Horní lišta pro konzistenci
	out := fmt.Sprintf("\n%s HELP\n", strings.ToUpper(string(proto)))
	out += strings.Repeat("-", 40) + "\n"

	out += fmt.Sprintf("Usage:   %s\n", spec.Usage)
	out += fmt.Sprintf("Example: %s\n", spec.Example)

	out += "\nAVAILABLE FLAGS FOR THIS SCAN:"

	// Seřadíme flagy, ať to nelítá
	keys := make([]string, 0, len(spec.UsableFlags))
	for k := range spec.UsableFlags {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		f := spec.UsableFlags[k]
		out += fmt.Sprintf("\n[%s]", k)
		out += fmt.Sprintf("\n  usage:       %s", f.Usage)

		if f.Description != "" {
			out += fmt.Sprintf("\n  description: %s", f.Description)
		}

		if f.Default != "" {
			out += fmt.Sprintf("\n  default:     %s", f.Default)
		}
		out += "\n"
	}

	out += strings.Repeat("-", 40)
	return out
}

func PrintHelp() {
	fmt.Println("NetSnitch - Intelligent Network Scanner")
	fmt.Println("=======================================")
	fmt.Println("\nUsage: netsnitch <protocol> [flags] <target>")

	// 1. Sekce SCANNERY
	fmt.Println("\nAVAILABLE SCANS:")

	// Seřadíme protokoly, aby nápověda neskákala
	protos := make([]string, 0, len(ScanRegistry))
	for p := range ScanRegistry {
		protos = append(protos, p)
	}
	sort.Strings(protos)

	for _, p := range protos {
		spec := ScanRegistry[p]
		fmt.Printf("\n[%s]\n", strings.ToUpper(p))
		fmt.Printf("  Usage:   %s\n", spec.Usage)
		fmt.Printf("  Example: %s\n", spec.Example)

		// Vypíšeme jen názvy flagů, které tento scan podporuje
		usable := make([]string, 0)
		for f := range spec.UsableFlags {
			usable = append(usable, "-"+f)
		}
		sort.Strings(usable)
		fmt.Printf("  Flags:   %s\n", strings.Join(usable, ", "))
	}

	// 2. Sekce GLOBÁLNÍ DETAILY FLAGŮ
	fmt.Println("\n" + strings.Repeat("-", 40))
	fmt.Println("FLAG DETAILS:")

	// Seřadíme všechny flagy z registru
	fKeys := make([]string, 0, len(FlagRegistry))
	for k := range FlagRegistry {
		fKeys = append(fKeys, k)
	}
	sort.Strings(fKeys)

	for _, k := range fKeys {
		f := FlagRegistry[k]
		fmt.Printf("\n[%s]\n", k)
		fmt.Printf("  usage:       %s\n", f.Usage)
		fmt.Printf("  description: %s\n", f.Description)
		if f.Default != "" {
			fmt.Printf("  default:     %s\n", f.Default)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 40))
	fmt.Println("Tip: Use 'netsnitch <protocol> -h' for protocol-specific help.")
}
