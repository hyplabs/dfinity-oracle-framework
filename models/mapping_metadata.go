package models

// MappingMetadata is the data required for the smart contract to store arbitrary key-values
type MappingMetadata struct {
	Key         string
	SummaryFunc func([]map[string]float64) map[string]float64
	Endpoints   []Endpoint
}
