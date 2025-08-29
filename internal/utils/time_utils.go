package utils

import "time"

func GetCurrentChampionshipsRound() (time.Time, time.Time) {
	return GetChampionshipsRoundByTime(time.Now())
}

func GetChampionshipsRoundByTime(targetTime time.Time) (time.Time, time.Time) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	firstRoundStartDate := time.Date(2025, time.June, 2, 0, 0, 0, 0, loc)
	roundDuration := 14 * 24 * time.Hour

	elapsed := targetTime.In(loc).Sub(firstRoundStartDate)
	roundsPassed := int(elapsed / roundDuration)

	roundStartDate := firstRoundStartDate.Add(time.Duration(roundsPassed) * roundDuration)
	roundEndDate := roundStartDate.Add(roundDuration)

	return roundStartDate, roundEndDate
}
