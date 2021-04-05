package models

// Endpoint is an endpoint configuration for the oracle
type Endpoint struct {
	Endpoint      string
	JSONPaths     map[string]string
	NormalizeFunc func(map[string]interface{}) (map[string]float64, error)
}
