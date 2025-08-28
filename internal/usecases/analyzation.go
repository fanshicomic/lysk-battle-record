package usecases

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"lysk-battle-record/internal/estimator"
	"lysk-battle-record/internal/models"
	"lysk-battle-record/internal/pkg"
)

type AnalyzeResponse struct {
	CombatPower models.CombatPower `json:"combat_power"`
}

func (s *LyskServer) AnalyzeCombatPower(c *gin.Context) {
	var input map[string]interface{}
	if err := c.BindJSON(&input); err != nil {
		logrus.Errorf("[Analysis] Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误", "detail": err.Error()})
		return
	}

	record := models.Record{}
	record.LevelType = pkg.GetValue(input, "关卡")
	record.LevelNumber = pkg.GetValue(input, "关数")
	record.LevelMode = pkg.GetValue(input, "模式")
	record.Attack = pkg.GetValue(input, "攻击")
	record.HP = pkg.GetValue(input, "生命")
	record.Defense = pkg.GetValue(input, "防御")
	record.Matching = pkg.GetValue(input, "对谱")
	record.MatchingBuff = pkg.GetValue(input, "对谱加成")
	record.CritRate = pkg.GetValue(input, "暴击")
	record.CritDmg = pkg.GetValue(input, "暴伤")
	record.EnergyRegen = pkg.GetValue(input, "加速回能")
	record.WeakenBoost = pkg.GetValue(input, "虚弱增伤")
	record.OathBoost = pkg.GetValue(input, "誓约增伤")
	record.OathRegen = pkg.GetValue(input, "誓约回能")
	record.Companion = pkg.GetValue(input, "搭档身份")
	record.SetCard = pkg.GetValue(input, "日卡")
	record.Stage = pkg.GetValue(input, "阶数")
	record.Weapon = pkg.GetValue(input, "武器")
	record.Buff = pkg.GetValue(input, "加成")
	record.TotalLevel = pkg.GetValue(input, "卡总等级")
	record.StarRank = pkg.GetValue(input, "星级")

	cpEstimator := estimator.NewCombatPowerEstimator()
	combatPower := cpEstimator.EstimateCombatPower(record)

	if record.LevelType != "" && record.LevelNumber != "" {
		record.CombatPower = combatPower
		evaluation := s.orbitRecordStore.EvaluateRecord(record)
		combatPower.Evaluation = evaluation
	}

	response := AnalyzeResponse{
		CombatPower: combatPower,
	}

	c.JSON(http.StatusOK, response)
}
