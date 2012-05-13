package main

import (
	"log"
	"math"
)

type StatisticalData struct {
	IndexName         string
	Instances         int
	Minimum           int
	Mean              int
	Maximum           int
	StandardDeviation int
}

func Gather(reviews Reviews, indexName string) StatisticalData {
	first, min, total, max, count := true, 0, 0, 0, 0
	for _, review := range reviews {
		if score, ok := review.Scores[indexName]; ok {
			if first || score < min {
				min = score
				first = false
			}
			if first || score > max {
				max = score
				first = false
			}
			total += score
			count++
		}
	}
	mean := int(float64(total) / float64(count))

	// http://en.wikipedia.org/wiki/Standard_deviation
	sqdev := 0.0
	for _, review := range reviews {
		x := review.Scores[indexName] - mean
		sqdev += math.Pow(float64(x), 2)
	}
	stdev := int(math.Sqrt(sqdev / float64(count)))

	return StatisticalData{
		IndexName:         indexName,
		Instances:         count,
		Minimum:           min,
		Mean:              mean,
		Maximum:           max,
		StandardDeviation: stdev,
	}
}

type AllStatisticalData map[string]StatisticalData

func GatherAll(reviews Reviews) AllStatisticalData {
	allStats := AllStatisticalData{}
	for indexName, _ := range IndexDefinitions {
		allStats[indexName] = Gather(reviews, indexName)
	}
	return allStats
}

func DeviationsFromMinimum(score int, stats StatisticalData) int {
	for i := 1; i <= 10; i++ {
		if score <= stats.Minimum+(i*stats.StandardDeviation) {
			return i
		}
	}
	log.Printf(
		"%d >10 standard deviations (%d) from minimum %d",
		score,
		stats.StandardDeviation,
		stats.Minimum,
	)
	return 10
}
