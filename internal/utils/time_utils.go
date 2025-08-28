package utils

import "time"

func GetCurrentChampionshipsRound() (time.Time, time.Time) {
	firstRoundStartDate := time.Date(2025, time.June, 2, 0, 0, 0, 0, time.UTC)
	roundDuration := 14 * 24 * time.Hour

	elapsed := time.Now().UTC().Sub(firstRoundStartDate)
	roundsPassed := int(elapsed / roundDuration)

	currentRoundStartDate := firstRoundStartDate.Add(time.Duration(roundsPassed) * roundDuration)
	currentRoundEndDate := currentRoundStartDate.Add(roundDuration)

	return currentRoundStartDate, currentRoundEndDate
}
