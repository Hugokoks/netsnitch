package fingerprint

type Engine struct {
	byProtocol map[string][]Pattern
	generic    []Pattern
	cache      map[string]*ServiceInfo
}

func NewEngine() *Engine {

	return &Engine{
		byProtocol: make(map[string][]Pattern),
		cache:      make(map[string]*ServiceInfo),
	}
}
