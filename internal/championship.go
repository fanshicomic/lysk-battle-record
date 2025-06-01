package internal

import (
	"time"
)

type TournamentPhase int

const (
	BeforeFirstTournament TournamentPhase = iota
	Cooldown
	Active
)

type Tournament interface {
	GetPhase() TournamentPhase
}

type TournamentInfo struct {
	Phase         TournamentPhase
	CurrentStart  time.Time // 当前这轮锦标赛的开始时间（周一早上5点）
	CurrentEnd    time.Time // 当前这轮锦标赛的结束时间（下下周周一凌晨0点）
	CooldownStart time.Time // 冷却期开始时间（凌晨0点）
	CooldownEnd   time.Time // 冷却期结束时间（凌晨5点）
	ShouldRefresh bool      // 是否需要刷新
}

func InitTournament() Tournament {
	t := &TournamentInfo{}
	t.updateTournamentInfo(time.Now())
	return t
}

func (t TournamentInfo) GetPhase() TournamentPhase {
	t.updateTournamentInfo(time.Now())
	return t.Phase
}

func (t TournamentInfo) updateTournamentInfo(now time.Time) {
	// 第一轮锦标赛开始时间：2025年5月19日 周一 5点
	firstStart := time.Date(2025, 5, 19, 5, 0, 0, 0, time.Local)

	if now.Before(firstStart) {
		t.Phase = BeforeFirstTournament
		t.CurrentStart = firstStart
		return
	}

	// 每轮14天
	elapsed := now.Sub(firstStart)
	days := int(elapsed.Hours() / 24)
	cycle := days / 14

	// 当前轮开始时间（本轮周一 5:00）
	currentStart := firstStart.AddDate(0, 0, cycle*14)
	currentEnd := currentStart.AddDate(0, 0, 14).Add(-5 * time.Hour) // 当前轮结束（冷却期开始）
	cooldownStart := currentEnd
	cooldownEnd := currentEnd.Add(5 * time.Hour) // 冷却期结束（下轮开始）
	shouldRefresh := false

	var phase TournamentPhase
	switch {
	case now.Before(cooldownStart):
		phase = Active
	case now.Before(cooldownEnd):
		phase = Cooldown
		shouldRefresh = true
	default:
		// 进入下一轮
		currentStart = currentStart.AddDate(0, 0, 14)
		currentEnd = currentStart.AddDate(0, 0, 14).Add(-5 * time.Hour)
		cooldownStart = currentEnd
		cooldownEnd = currentEnd.Add(5 * time.Hour)
		phase = Active
		shouldRefresh = true
	}

	t.Phase = phase
	t.CurrentStart = currentStart
	t.CurrentEnd = currentEnd
	t.CooldownStart = cooldownStart
	t.CooldownEnd = cooldownEnd
	t.ShouldRefresh = shouldRefresh
}
