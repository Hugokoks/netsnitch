package fingerprint

func (e *Engine) Detect(raw string) *ServiceInfo {

	if raw == "" {
		return nil
	}

	if cached := e.getCached(raw); cached != nil {
		return cached
	}

	proto := guessProtocol(raw)

	var candidates []Pattern

	if proto != "" {
		candidates = e.byProtocol[proto]
	} else {
		candidates = e.generic
	}

	for _, p := range candidates {

		matches := p.Regex.FindStringSubmatch(raw)

		if matches == nil {
			continue
		}

		info := &ServiceInfo{
			Name: p.Service,
			Raw:  raw,
		}

		for _, param := range p.Params {

			if param.Pos >= len(matches) {
				continue
			}

			switch param.Name {

			case "service.version":

				info.Version = matches[param.Pos]

			case "service.product":

				info.Product = matches[param.Pos]

			case "service.vendor":

				info.Vendor = matches[param.Pos]
			}
		}

		e.storeCache(raw, info)

		return info
	}

	return nil
}
