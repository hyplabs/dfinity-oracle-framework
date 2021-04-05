package summary

import (
	"reflect"
	"testing"
)

func TestRemoveOutlierLarge(t *testing.T) {
	dataset := []float64{16.6, 23.4, 13.5, 1.1, 52.2, 5.5, 17.1, 50.2, 35.5, 100000000000000}
	expectedDataset := []float64{1.1, 5.5, 13.5, 16.6, 17.1, 23.4, 35.5, 50.2, 52.2}

	result := RemoveOutlier(dataset)

	if !reflect.DeepEqual(result, expectedDataset) {
		t.Errorf("Incorrect dataset from remove outlier, expected %v, got %v", expectedDataset, result)
	}
}

func TestRemoveOutlieSmall(t *testing.T) {
	dataset := []float64{-16.0, -100.0, -18.0, -6.0, -2000, 16.6, 23.4, 13.5, 1.1, 52.2, 5.5, 17.1, 35.5}
	expectedDataset := []float64{-100, -18, -16, -6, 1.1, 5.5, 13.5, 16.6, 17.1, 23.4, 35.5, 52.2}

	result := RemoveOutlier(dataset)

	if !reflect.DeepEqual(result, expectedDataset) {
		t.Errorf("Incorrect dataset from remove outlier, expected %v, got %v", expectedDataset, result)
	}
}
