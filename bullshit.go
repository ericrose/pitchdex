package main

func calculateBullshit(review Review, allStats AllStatisticalData) int {
	pf := DeviationsFromMinimum(
		review.Scores["Pitchformulaity"],
		allStats["Pitchformulaity"],
	)
	sl := DeviationsFromMinimum(
		review.Scores["Sentence length"],
		allStats["Sentence length"],
	)
	wc := DeviationsFromMinimum(
		review.Scores["Word count"],
		allStats["Word count"],
	)
	wi := DeviationsFromMinimum(
		review.Scores["Words invented"],
		allStats["Words invented"],
	)
	return 10*pf + 5*sl + 2*wc + 1*wi
}
