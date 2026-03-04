package fingerprint

import (
	"encoding/xml"
	"os"
	"regexp"
)

func (e *Engine) LoadRecogFile(path string) error {

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var db RecogXML

	err = xml.Unmarshal(data, &db)
	if err != nil {
		return err
	}

	for _, f := range db.Fingerprints {

		re, err := regexp.Compile("(?i)" + f.Pattern)
		if err != nil {
			continue
		}

		p := Pattern{
			Regex:   re,
			Service: db.Protocol,
			Params:  f.Params,
		}

		if db.Protocol != "" {
			e.byProtocol[db.Protocol] =
				append(e.byProtocol[db.Protocol], p)

		} else {

			e.generic =
				append(e.generic, p)
		}
	}

	return nil
}
