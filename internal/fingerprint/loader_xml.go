package fingerprint

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// 1. TYPY DAT (Krabičky, do kterých to budeme sypat)

// XMLParam odpovídá tagu <param pos="1" name="service.version" />
type XMLParam struct {
	Pos   int    `xml:"pos,attr"`
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// XMLFingerprint odpovídá celému bloku <fingerprint>...</fingerprint>
type XMLFingerprint struct {
	Pattern     string     `xml:"pattern,attr"`
	Description string     `xml:"description"`
	Params      []XMLParam `xml:"param"` // Tímhle Go najde všechny vnořené parametry
}

// RecogXML je obal pro celý soubor (např. mysql_banners.xml)
type RecogXML struct {
	Fingerprints []XMLFingerprint `xml:"fingerprint"`
}

// Pattern je to, co budeme mít v paměti pro skenování
type Pattern struct {
	Regex   *regexp.Regexp
	Service string
	Params  []XMLParam // Tady si necháme instrukce, jak z regexu vytáhnout verzi
}

// Tady bude náš "mozek" - seznam všech pravidel
var AllFingerprints []Pattern

// 2. FUNKCE PRO NAČTENÍ (Vysavač)

func LoadAllBanners(bannersDir string) error {
	// Najdeme všechny .xml soubory ve složce (ftp.xml, mysql.xml, atd.)
	files, err := filepath.Glob(filepath.Join(bannersDir, "*.xml"))
	if err != nil {
		return err
	}

	for _, file := range files {
		// A. Přečteme soubor z disku do paměti (jako bajty)
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		// B. Použijeme ten Unmarshal, o kterém jsme mluvili
		var recog RecogXML
		if err := xml.Unmarshal(data, &recog); err != nil {
			fmt.Printf("Chyba v XML %s: %v\n", file, err)
			continue
		}

		// C. Teď ty surové XML data převedeme na "živé" Patterny
		for _, f := range recog.Fingerprints {
			// Zkompilujeme regex (přidáme (?i) pro case-insensitive)
			re, err := regexp.Compile("(?i)" + f.Pattern)
			if err != nil {
				// Pokud je v XML blbě napsaný regex, přeskočíme ho
				continue
			}

			// Vytvoříme náš finální objekt a hodíme ho do pole
			p := Pattern{
				Regex:   re,
				Service: f.Description,
				Params:  f.Params, // Tady se zkopírují ty pos="1" instrukce
			}
			AllFingerprints = append(AllFingerprints, p)
		}
	}

	fmt.Printf("[*] Úspěšně načteno %d fingerprintů\n", len(AllFingerprints))
	return nil
}
