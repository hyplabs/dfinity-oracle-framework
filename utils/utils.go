package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/hyplabs/dfinity-oracle-framework/models"
	"github.com/oliveagle/jsonpath"
)

func getEndpoint(endpoint string) ([]byte, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// GetAPIInfo takes a given endpoint and parses the endpoint data
// Currently assumes all output is in map of floats format
func GetAPIInfo(e models.Endpoint) (map[string]float64, error) {
	responseBody, err := getEndpoint(e.Endpoint)
	if err != nil {
		return map[string]float64{}, err
	}

	var jsonData interface{}
	json.Unmarshal(responseBody, &jsonData)

	result := make(map[string]interface{})
	for fieldName, jsonPath := range e.JSONPaths {
		resp, err := jsonpath.JsonPathLookup(jsonData, jsonPath)
		if err != nil {
			return map[string]float64{}, err
		}
		result[fieldName] = resp
	}

	if e.NormalizeFunc != nil {
		normalizedResult, err := e.NormalizeFunc(result)
		if err != nil {
			return map[string]float64{}, err
		}
		return normalizedResult, nil
	} else {
		normalizedResult := make(map[string]float64)
		for k, v := range result {
			normalizedResult[k] = v.(float64)
		}
		return normalizedResult, nil
	}
}
