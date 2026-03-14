package output

import (
	"fmt"
)

func IsMostlyPrintable(s string) bool {
	if s == "" {
		return true
	}

	runes := []rune(s)
	if len(runes) == 0 {
		return true
	}

	printable := 0
	for _, r := range runes {
		if r == '\n' || r == '\r' || r == '\t' || (r >= 32 && r <= 126) {
			printable++
		}
	}

	return printable*100/len(runes) >= 80
}

func TextPreview(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + "..."
}

func HexPreview(s string, maxBytes int) string {
	b := []byte(s)
	if len(b) > maxBytes {
		b = b[:maxBytes]
	}
	return fmt.Sprintf("%x", b)
}

func FormatTextOrHex(s string, maxText int, maxHex int) (label string, value string) {
	if IsMostlyPrintable(s) {
		return "banner", TextPreview(s, maxText)
	}
	return "banner_hex", HexPreview(s, maxHex)
}
