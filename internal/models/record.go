package models

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"lysk-battle-record/internal/pkg"
	"lysk-battle-record/internal/utils"
)

type Record struct {
	Id           string      `json:"id"`
	RowNumber    int         `json:"row_number"`
	UserID       string      `json:"userID"`
	Nickname     string      `json:"nickname,omitempty"`
	LevelType    string      `json:"关卡"`
	LevelNumber  string      `json:"关数"`
	LevelMode    string      `json:"模式"`
	Attack       string      `json:"攻击"`
	HP           string      `json:"生命"`
	Defense      string      `json:"防御"`
	Matching     string      `json:"对谱"`
	MatchingBuff string      `json:"对谱加成"`
	CritRate     string      `json:"暴击"`
	CritDmg      string      `json:"暴伤"`
	EnergyRegen  string      `json:"加速回能"`
	WeakenBoost  string      `json:"虚弱增伤"`
	OathBoost    string      `json:"誓约增伤"`
	OathRegen    string      `json:"誓约回能"`
	TotalLevel   string      `json:"卡总等级"`
	Note         string      `json:"备注"`
	Companion    string      `json:"搭档身份"`
	SetCard      string      `json:"日卡"`
	Stage        string      `json:"阶数"`
	Weapon       string      `json:"武器"`
	Buff         string      `json:"加成"`
	Time         string      `json:"时间"`
	StarRank     string      `json:"星级"`
	CombatPower  CombatPower `json:"战力值"`
	Deleted      bool        `json:"deleted"`
}

type Records []Record

func (r Record) validateCommon() (bool, error) {
	if !r.validateLevelType() {
		return false, fmt.Errorf("无效的关卡类型: %s", r.LevelType)
	}

	if !r.validateAttack() {
		return false, fmt.Errorf("攻击数值错误: %s", r.Attack)
	}

	if _, err := r.validateDefence(); err != nil {
		return false, err
	}

	if _, err := r.validateHP(); err != nil {
		return false, err
	}

	if !r.validateMatching() {
		return false, fmt.Errorf("对谱类型错误: %s", r.Matching)
	}

	if !r.validateMatchingBuff() {
		return false, fmt.Errorf("对谱加成错误: %s", r.MatchingBuff)
	}

	if !r.validateCritRate() {
		return false, fmt.Errorf("暴击率错误: %s", r.CritRate)
	}

	if !r.validateCritDmg() {
		return false, fmt.Errorf("暴击伤害错误: %s", r.CritDmg)
	}

	if !r.validateWeakenBoost() {
		return false, fmt.Errorf("虚弱增伤错误: %s", r.WeakenBoost)
	}

	if !r.validateOathBoost() {
		return false, fmt.Errorf("誓约增伤错误: %s", r.OathBoost)
	}

	if !r.validateOathRegen() {
		return false, fmt.Errorf("誓约回能错误: %s", r.OathRegen)
	}

	if !r.validateEnergyRegen() {
		return false, fmt.Errorf("加速回能错误: %s", r.EnergyRegen)
	}

	if !r.validateRegen() {
		return false, fmt.Errorf("回能总和错误: %s + %s，面板总回能不能大于48", r.EnergyRegen, r.OathRegen)
	}

	if !r.validateStage() {
		return false, fmt.Errorf("阶数错误: %s", r.Stage)
	}

	if !r.validateWeapon() {
		return false, fmt.Errorf("武器错误: %s", r.Weapon)
	}

	if !r.validateCompanionSetCard() {
		return false, fmt.Errorf("搭档身份与日卡不匹配: %s - %s", r.Companion, r.SetCard)
	}

	if r.SetCard == "无套装" && r.Stage != "无套装" {
		return false, fmt.Errorf("无套装时阶数必须为无套装: %s", r.Stage)
	}

	if r.SetCard != "无套装" && r.Stage == "无套装" {
		return false, fmt.Errorf("有套装时阶数不能为无套装: %s", r.Stage)
	}

	if !r.validateTotalLevel() {
		return false, fmt.Errorf("卡面总等级错误, 请填写卡面等级总和。如不确定请填留空: %s", r.TotalLevel)
	}

	if pass, err := r.ValidateNote(); !pass {
		return false, err
	}

	return true, nil
}

func (r Record) validateLevelType() bool {
	validTypes := map[string]bool{
		"光":   true,
		"火":   true,
		"冰":   true,
		"能量": true,
		"引力": true,
		"开放": true,
		"A4":   true,
		"B4":   true,
		"C4":   true,
	}
	return validTypes[r.LevelType]
}

func (r Record) validateLevelNumber() bool {
	maxEasyLevelNumber := map[string]int{
		"光":   210,
		"火":   210,
		"冰":   210,
		"能量": 180,
		"引力": 150,
		"开放": 300,
	}

	maxHardLevelNumber := map[string]int{
		"光":   0,
		"火":   0,
		"冰":   0,
		"能量": 0,
		"引力": 0,
		"开放": 60,
	}
	levelInfo := strings.Split(r.LevelNumber, "_")
	var levelPart string

	if len(levelInfo) == 2 {
		levelPart = levelInfo[1]
	}

	levelNumber, _ := strconv.Atoi(levelInfo[0])
	validEasyLevelNumber := r.LevelMode == "稳定" && levelNumber <= maxEasyLevelNumber[r.LevelType]
	validHardLevelNumber := r.LevelMode == "波动" && levelNumber <= maxHardLevelNumber[r.LevelType]
	isValidNumber := levelNumber > 0 && (validEasyLevelNumber || validHardLevelNumber)
	isValidPart := (levelNumber%10 != 0 && levelPart == "") ||
		(levelNumber%10 == 0 && (levelPart == "上" || levelPart == "下")) ||
		(r.LevelMode == "波动" && levelNumber%5 == 0 && (levelPart == "上" || levelPart == "下"))
	if !isValidNumber || !isValidPart {
		return false
	}

	return true
}

func (r Record) validateLevelMode() bool {
	validModes := map[string]bool{
		"稳定": true,
		"波动": true,
	}

	if !validModes[r.LevelMode] {
		return false
	}

	if r.LevelType != "开放" && r.LevelMode == "波动" {
		return false
	}

	return true
}

func (r Record) validateAttack() bool {
	maxAttack := 1229 * 1.9 * 6
	n, err := strconv.ParseFloat(r.Attack, 64)
	if err != nil || n <= 0 || n > maxAttack {
		return false
	}

	return true
}

func (r Record) validateDefence() (bool, error) {
	maxDefence := 614 * 1.9 * 6
	if r.Defense == "" {
		r.Defense = "0"
	}
	n, err := strconv.ParseFloat(r.Defense, 64)
	if err != nil || n < 0 || n > maxDefence {
		return false, fmt.Errorf("防御值错误: %s", r.Defense)
	}

	defenceCompanions := map[string]bool{
		"光猎":         true,
		"永恒先知":     true,
		"远空执舰官":   true,
		"利莫里亚海神": true,
		"银翼恶魔":     true,
	}

	if _, ok := defenceCompanions[r.Companion]; ok && n == 0 {
		return false, fmt.Errorf("搭档 %s 的防御值不能为 0", r.Companion)
	}

	return true, nil
}

func (r Record) validateHP() (bool, error) {
	maxHP := 24594 * 1.9 * 6
	if r.HP == "" {
		r.HP = "0"
	}
	n, err := strconv.ParseFloat(r.HP, 64)
	if err != nil || n < 0 || n > maxHP {
		return false, fmt.Errorf("生命值错误: %s", r.HP)
	}

	hpCompanions := map[string]bool{
		"潮汐之神": true,
		"深渊主宰": true,
		"暗蚀国王": true,
		"终末之神": true,
	}

	if _, ok := hpCompanions[r.Companion]; ok && n == 0 {
		return false, fmt.Errorf("搭档 %s 的生命值不能为 0", r.Companion)
	}

	return true, nil
}

func (r Record) validateMatching() bool {
	validMatching := map[string]bool{
		"顺":     true,
		"逆":     true,
		"不确定": true,
	}
	return validMatching[r.Matching]
}

func (r Record) validateMatchingBuff() bool {
	validMatchingBuff := map[string]bool{
		"30":     true,
		"25":     true,
		"20":     true,
		"15":     true,
		"10":     true,
		"5":      true,
		"0":      true,
		"不确定": true,
	}
	return validMatchingBuff[r.MatchingBuff]
}

func (r Record) validateWeapon() bool {
	validWeapons := map[string]bool{
		"专武":   true,
		"重剑":   true,
		"手枪":   true,
		"法杖":   true,
		"单手剑": true,
	}
	return validWeapons[r.Weapon]
}

func (r Record) validateBuff() bool {
	validBuffs := map[string]bool{
		"0":  true,
		"10": true,
		"20": true,
		"30": true,
		"40": true,
	}
	return validBuffs[r.Buff]
}

func (r Record) validateCompanionSetCard() bool {
	partnerMap := utils.GetPartnerCompanionMap()
	setCardMap := utils.GetPartnerSetCardMap()

	for partner, companions := range partnerMap {
		for _, companion := range companions {
			if companion == r.Companion {
				if setCards, exists := setCardMap[partner]; exists {
					return setCards[r.SetCard]
				}
				return false
			}
		}
	}
	return false
}

func (r Record) validatePartnerAndLevelType() bool {
	levelPartnerMap := map[string]string{
		"光":   "沈星回",
		"火":   "祁煜",
		"冰":   "黎深",
		"能量": "秦彻",
		"引力": "夏以昼",
	}

	// Check if level type has specific partner requirement first
	requiredMain, hasRequirement := levelPartnerMap[r.LevelType]
	if !hasRequirement {
		return true
	}

	partnerCompanionsMap := utils.GetPartnerCompanionMap()

	// Find which main character this partner belongs to
	var partner string
	found := false
	for pn, companions := range partnerCompanionsMap {
		for _, companion := range companions {
			if companion == r.Companion {
				partner = pn
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		return false
	}

	// Check if partner matches required main character for this level type
	if partner != requiredMain {
		return false
	}

	return true
}

func (r Record) validateCritRate() bool {
	n, err := strconv.ParseFloat(r.CritRate, 64)
	if err != nil || n < 0 || n > 100 {
		return false
	}

	return true
}

func (r Record) validateCritDmg() bool {
	maxCritDmg := float64(150 + 20*2) // 20 max from each sun card itself
	maxCritDmg += 22.4 * 4            // 22.4 max from each moon card core
	maxCritDmg += 14.4 * 2 * 6        // 14.4 max from each core attribute
	n, err := strconv.ParseFloat(r.CritDmg, 64)
	if err != nil || n < 0 || n > maxCritDmg {
		return false
	}

	return true
}

func (r Record) validateEnergyRegen() bool {
	if r.EnergyRegen == "" {
		return true
	}

	n, err := strconv.ParseFloat(r.EnergyRegen, 64)
	if err != nil || n < 0 || n > 48 {
		return false
	}

	return true
}

func (r Record) validateOathRegen() bool {
	if r.OathRegen == "" {
		return true
	}

	n, err := strconv.ParseFloat(r.OathRegen, 64)
	if err != nil || n < 0 || n > 40 {
		return false
	}

	return true
}

func (r Record) validateWeakenBoost() bool {
	maxWeakenBoost := 18.2 * 4   // 18.2 max from each moon card core
	maxWeakenBoost += 11 * 2 * 6 // 11 max from each core attribute
	n, err := strconv.ParseFloat(r.WeakenBoost, 64)
	if err != nil || n < 0 || n > maxWeakenBoost {
		return false
	}

	return true
}

func (r Record) validateOathBoost() bool {
	if r.OathBoost == "" {
		return true
	}

	maxOathBoost := float64(14 * 2) // 62.4 max from each moon card core
	maxOathBoost += 5.6 * 2 * 6     // 8.4 max from each core attribute
	n, err := strconv.ParseFloat(r.OathBoost, 64)
	if err != nil || n < 0 || n > maxOathBoost {
		return false
	}

	return true
}

func (r Record) validateRegen() bool {
	energy, _ := strconv.ParseFloat(r.EnergyRegen, 64)
	oath, _ := strconv.ParseFloat(r.OathRegen, 64)

	if energy+oath > 48 {
		return false
	}

	return true
}

func (r Record) validateStage() bool {
	validStages := map[string]bool{
		"I":      true,
		"II":     true,
		"III":    true,
		"IV":     true,
		"无套装": true,
	}

	return validStages[r.Stage]
}

func (r Record) validateTotalLevel() bool {
	if r.TotalLevel == "" {
		return true
	}

	totalLevel, err := strconv.Atoi(r.TotalLevel)
	if err != nil || totalLevel <= 0 || totalLevel > 480 {
		return false
	}

	return true
}

func (r Record) validateStarRank() bool {
	if r.LevelMode != "波动" {
		return r.StarRank == ""
	}

	validStarRanks := map[string]bool{
		"零星": true,
		"一星": true,
		"二星": true,
		"三星": true,
	}

	return validStarRanks[r.StarRank]
}

func (r Record) ValidateNote() (bool, error) {
	if utf8.RuneCountInString(r.Note) > 40 {
		return false, fmt.Errorf("备注最长30个字: %s", r.Note)
	}

	detector, err := pkg.NewDetector()
	if err != nil {
		return false, err
	}

	if detector.ContainsSensitiveWords(r.Note) {
		return false, fmt.Errorf("备注中包含敏感词")
	}

	return true, nil
}

func (r Record) ValidateOrbit() (bool, error) {
	if !r.validateLevelNumber() {
		return false, fmt.Errorf("关数错误: %s - %s - %s", r.LevelType, r.LevelMode, r.LevelNumber)
	}

	if !r.validateLevelMode() {
		return false, fmt.Errorf("关卡模式错误: %s", r.LevelMode)
	}

	if !r.validateStarRank() {
		return false, fmt.Errorf("波动关卡通关星级错误: %s", r.StarRank)
	}

	if !r.validatePartnerAndLevelType() {
		return false, fmt.Errorf("搭档身份与关卡类型不匹配: %s - %s", r.Companion, r.LevelType)
	}

	return r.validateCommon()
}

func (r Record) ValidateChampionships() (bool, error) {
	if !r.validateBuff() {
		return false, fmt.Errorf("锦标赛加成错误: %s", r.Buff)
	}

	return r.validateCommon()
}

func (r Record) GetHash() string {
	data := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s",
		r.LevelType, r.LevelNumber, r.LevelMode, r.Attack, r.HP, r.Defense, r.Matching, r.MatchingBuff,
		r.CritRate, r.CritDmg, r.EnergyRegen, r.WeakenBoost, r.OathBoost,
		r.OathRegen, r.Companion, r.SetCard, r.Stage, r.Weapon, r.Buff, r.TotalLevel, r.StarRank,
	)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (r Records) Filter(filterFunc func(Record) bool) []Record {
	result := Records{}
	for _, r := range r {
		if filterFunc(r) {
			result = append(result, r)
		}
	}
	return result
}

func (r Records) SortByTimeDesc() []Record {
	sort.Slice(r, func(i, j int) bool {
		ti, _ := time.Parse("2006-01-02T15:04:05Z", r[i].Time)
		tj, _ := time.Parse("2006-01-02T15:04:05Z", r[j].Time)
		return tj.Before(ti)
	})
	return r
}

func (r Records) Pagination(offset, size int) []Record {
	total := len(r)
	if offset >= total {
		return []Record{}
	}
	end := offset + size
	if end > total {
		end = total
	}
	return r[offset:end]
}

func (r Record) ToStats() Stats {
	attack, _ := strconv.Atoi(r.Attack)
	hp, _ := strconv.Atoi(r.HP)
	defense, _ := strconv.Atoi(r.Defense)
	totalLevel, _ := strconv.Atoi(r.TotalLevel)
	if totalLevel == 0 {
		totalLevel = 480
	}

	matchingBuff, _ := strconv.ParseFloat(r.MatchingBuff, 64)
	if r.MatchingBuff == "不确定" {
		matchingBuff = 0
	}

	critRate, _ := strconv.ParseFloat(r.CritRate, 64)
	critDmg, _ := strconv.ParseFloat(r.CritDmg, 64)
	energyRegen, _ := strconv.ParseFloat(r.EnergyRegen, 64)
	weakenBoost, _ := strconv.ParseFloat(r.WeakenBoost, 64)
	oathBoost, _ := strconv.ParseFloat(r.OathBoost, 64)
	oathRegen, _ := strconv.ParseFloat(r.OathRegen, 64)
	buff, _ := strconv.ParseFloat(r.Buff, 64)

	return Stats{
		Attack:       attack,
		HP:           hp,
		Defense:      defense,
		Matching:     r.Matching,
		MatchingBuff: matchingBuff,
		CritRate:     critRate,
		CritDmg:      critDmg,
		EnergyRegen:  energyRegen,
		WeakenBoost:  weakenBoost,
		OathBoost:    oathBoost,
		OathRegen:    oathRegen,
		TotalLevel:   totalLevel,
		Companion:    r.Companion,
		SetCard:      r.SetCard,
		Stage:        r.Stage,
		Weapon:       r.Weapon,
		Buff:         buff,
	}
}

func (r Record) GenerateLevelKey() string {
	// For championships: round-leveltype
	if r.LevelType == "A4" || r.LevelType == "B4" || r.LevelType == "C4" {
		recordTime, err := time.Parse(time.RFC3339, r.Time)
		if err != nil {
			recordTime = time.Now()
		}
		start, _ := utils.GetChampionshipsRoundByTime(recordTime)
		roundKey := start.Format("2006-01-02") // Use start date as round identifier
		return roundKey + "-" + r.LevelType
	}
	// For orbit: leveltype-levelnumber-levelmode
	return r.LevelType + "-" + r.LevelNumber + "-" + r.LevelMode
}
