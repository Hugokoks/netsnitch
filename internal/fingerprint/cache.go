package fingerprint

func (e *Engine) getCached(banner string) *ServiceInfo {

	if v, ok := e.cache[banner]; ok {
		return v
	}

	return nil
}

func (e *Engine) storeCache(banner string, info *ServiceInfo) {
	e.cache[banner] = info
}
