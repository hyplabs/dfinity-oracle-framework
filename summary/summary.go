package summary

import (
	"math"
)

// Mean: Returns the man of the dataset
func Mean(dataset []map[string]float64) map[string]float64 {
	result := make(map[string]float64)

	for key, values := range groupByKey(dataset) {
		result[key] = meanOfArray(values)
	}

	return result
}

// Median: Returns the median of the dataset
func Median(dataset []map[string]float64) map[string]float64 {
	result := make(map[string]float64)

	for key, values := range groupByKey(dataset) {
		result[key] = medianOfArray(values)
	}

	return result
}

// Mode: Returns the mode of the dataset
func Mode(dataset []map[string]float64) map[string]float64 {
	result := make(map[string]float64)

	for key, values := range groupByKey(dataset) {
		result[key] = modeOfArray(values)
	}

	return result
}

// MeanWithoutOutliers: Returns the mean of the dataset after removing values outside 2 standard deviations from the median
func MeanWithoutOutliers(dataset []map[string]float64) map[string]float64 {
	result := make(map[string]float64)

	for key, values := range groupByKey(dataset) {
		result[key] = meanOfArray(RemoveOutlier(values))
	}

	return result
}

// MedianWithoutOutliers: Returns the median of the dataset after removing values outside 2 std deviations from the median
func MedianWithoutOutliers(dataset []map[string]float64) map[string]float64 {
	result := make(map[string]float64)

	for key, values := range groupByKey(dataset) {
		result[key] = medianOfArray(RemoveOutlier(values))
	}

	return result
}

// groupBykey: returns the dataset grouped by keys from each entry in the dataset
func groupByKey(dataset []map[string]float64) map[string][]float64 {
	result := make(map[string][]float64)
	for _, entry := range dataset {
		for k, v := range entry {
			result[k] = append(result[k], v)
		}
	}
	return result
}

func RemoveOutlier(dataset []float64) []float64 {
	if len(dataset) <= 2 {
		return dataset
	}

	var sum float64 = 0.0
	var squaredSum float64 = 0.0
	for _, x := range dataset {
		sum += x
		squaredSum += x * x
	}
	var mean float64 = sum / float64(len(dataset))
	var standardDeviation float64 = math.Sqrt(squaredSum/float64(len(dataset)) - (mean * mean))

	slicedData := make([]float64, 0)

	median := medianOfArray(dataset)
	for _, x := range dataset {
		if median-(2*standardDeviation) <= x && x <= median+(2*standardDeviation) {
			slicedData = append(slicedData, x)
		}
	}
	return slicedData
}

func meanOfArray(dataset []float64) float64 {
	if len(dataset) == 1 {
		return dataset[0]
	}

	var sum float64 = 0
	for i := 0; i < len(dataset); i++ {
		sum += dataset[i]
	}

	return sum / float64(len(dataset))
}

func medianOfArray(dataset []float64) float64 {
	if len(dataset)%2 == 0 {
		return (dataset[len(dataset)/2] + dataset[(len(dataset)/2)-1]) / 2
	} else {
		return dataset[(len(dataset)-1)/2]
	}
}

func modeOfArray(dataset []float64) float64 {
	counter := make(map[float64]int)
	for _, x := range dataset {
		counter[x]++
	}
	modeX, modeCount := 0.0, 0
	for x, count := range counter {
		if count > modeCount {
			modeCount = count
			modeX = x
		}
	}
	return modeX
}
