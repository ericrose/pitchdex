package main

import (
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

func Gather(r *Reviews, indexName string) StatisticalData {
	first, min, total, max, count := true, 0, 0, 0, 0
	for _, review := range *r {
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
	sqdev := 0.0
	for _, review := range *r {
		// http://en.wikipedia.org/wiki/Standard_deviation
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

func GatherAll(r *Reviews) AllStatisticalData {
	allStats := AllStatisticalData{}
	for indexName, _ := range IndexDefinitions {
		allStats[indexName] = Gather(r, indexName)
	}
	return allStats
}

func DeviationsFromMinimum(score int, stats StatisticalData) int {
	for i := 1; i <= 10; i++ {
		if score <= stats.Minimum+(i*stats.StandardDeviation) {
			return i
		}
	}
	return 10
}

func calculateBullshit(review Review, allStats AllStatisticalData) int {
	return 10*DeviationsFromMinimum(
		review.Scores["Pitchformulaity"],
		allStats["Pitchformulaity"],
	) +
		5*DeviationsFromMinimum(
			review.Scores["Sentence length"],
			allStats["Sentence length"],
		) +
		2*DeviationsFromMinimum(
			review.Scores["Word count"],
			allStats["Word count"],
		) +
		review.Scores["Words invented"]
}
