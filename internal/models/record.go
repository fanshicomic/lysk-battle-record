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
)

type Record struct {
	Id           string `json:"id"`
	RowNumber    int    `json:"row_number"`
	UserID       string `json:"userID"`
	LevelType    string `json:"关卡"`
	LevelNumber  string `json:"关数"`
	LevelMode    string `json:"模式"`
	Attack       string `json:"攻击"`
	HP           string `json:"生命"`
	Defense      string `json:"防御"`
	Matching     string `json:"对谱"`
	MatchingBuff string `json:"对谱加成"`
	CritRate     string `json:"暴击"`
	CritDmg      string `json:"暴伤"`
	EnergyRegen  string `json:"加速回能"`
	WeakenBoost  string `json:"虚弱增伤"`
	OathBoost    string `json:"誓约增伤"`
	OathRegen    string `json:"誓约回能"`
	TotalLevel   string `json:"卡总等级"`
	Note         string `json:"备注"`
	Partner      string `json:"搭档身份"`
	SetCard      string `json:"日卡"`
	Stage        string `json:"阶数"`
	Weapon       string `json:"武器"`
	Buff         string `json:"加成"`
	Time         string `json:"时间"`
	Deleted      bool   `json:"deleted"`
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

	if !r.validatePartner() {
		return false, fmt.Errorf("搭档身份错误: %s", r.Partner)
	}

	if !r.validateSetCard() {
		return false, fmt.Errorf("日卡错误: %s", r.SetCard)
	}

	if !r.validateStage() {
		return false, fmt.Errorf("阶数错误: %s", r.Stage)
	}

	if !r.validateWeapon() {
		return false, fmt.Errorf("武器错误: %s", r.Weapon)
	}

	if !r.validatePartnerSetCard() {
		return false, fmt.Errorf("搭档身份与日卡不匹配: %s - %s", r.Partner, r.SetCard)
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
		"冰":   180,
		"能量": 150,
		"引力": 120,
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
	isValidPart := (levelNumber%10 != 0 && levelPart == "") || (levelNumber%10 == 0 && levelPart == "上" || levelPart == "下")
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

	defencePartner := map[string]bool{
		"光猎":         true,
		"永恒先知":     true,
		"远空执舰官":   true,
		"利莫里亚海神": true,
	}

	if _, ok := defencePartner[r.Partner]; ok && n == 0 {
		return false, fmt.Errorf("搭档 %s 的防御值不能为 0", r.Partner)
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

	hpPartner := map[string]bool{
		"潮汐之神": true,
		"深渊主宰": true,
		"暗蚀国王": true,
	}

	if _, ok := hpPartner[r.Partner]; ok && n == 0 {
		return false, fmt.Errorf("搭档 %s 的生命值不能为 0", r.Partner)
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

func (r Record) validatePartnerSetCard() bool {
	partnerSetCardMap := map[string]map[string]bool{
		"沈星回": {
			"夜誓": true, "末夜": true, "逐光": true, "鎏光": true, "睱日": true, "弦光": true, "心晴": true, "匿光": true, "无套装": true,
		},
		"黎深": {
			"拥雪": true, "永恒": true, "夜色": true, "静谧": true, "心晴": true, "深林": true, "无套装": true,
		},
		"祁煜": {
			"雾海": true, "神殿": true, "深海": true, "坠浪": true, "点染": true, "斑斓": true, "心晴": true, "碧海": true, "无套装": true,
		},
		"秦彻": {
			"深渊": true, "掠心": true, "锋尖": true, "戮夜": true, "无套装": true,
		},
		"夏以昼": {
			"寂路": true, "远空": true, "长昼": true, "离途": true, "无套装": true,
		},
	}

	// 搭档身份到主角名的映射
	partnerToMain := map[string]string{
		"暗蚀国王": "沈星回", "光猎": "沈星回", "逐光骑士": "沈星回", "遥远少年": "沈星回", "Evol特警": "沈星回", "深空猎人": "沈星回",
		"九黎司命": "黎深", "永恒先知": "黎深", "极地军医": "黎深", "黎明抹杀者": "黎深", "临空医生": "黎深",
		"利莫里亚海神": "祁煜", "潮汐之神": "祁煜", "深海潜行者": "祁煜", "画坛新锐": "祁煜", "海妖魅影": "祁煜", "艺术家": "祁煜",
		"深渊主宰": "秦彻", "无尽掠夺者": "秦彻", "异界来客": "秦彻",
		"终极兵器X-02": "夏以昼", "远空执舰官": "夏以昼", "深空飞行员": "夏以昼",
	}

	main, ok := partnerToMain[r.Partner]
	if !ok {
		return false
	}
	setMap, ok := partnerSetCardMap[main]
	if !ok {
		return false
	}
	return setMap[r.SetCard]
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

func (r Record) validatePartner() bool {
	validPartner := map[string]bool{
		"暗蚀国王":     true,
		"光猎":         true,
		"逐光骑士":     true,
		"遥远少年":     true,
		"Evol特警":     true,
		"深空猎人":     true,
		"九黎司命":     true,
		"永恒先知":     true,
		"极地军医":     true,
		"黎明抹杀者":   true,
		"临空医生":     true,
		"利莫里亚海神": true,
		"潮汐之神":     true,
		"深海潜行者":   true,
		"画坛新锐":     true,
		"海妖魅影":     true,
		"艺术家":       true,
		"深渊主宰":     true,
		"无尽掠夺者":   true,
		"异界来客":     true,
		"终极兵器X-02": true,
		"远空执舰官":   true,
		"深空飞行员":   true,
	}

	return validPartner[r.Partner]
}

func (r Record) validateSetCard() bool {
	validSetCard := map[string]bool{
		"夜誓":   true,
		"鎏光":   true,
		"末夜":   true,
		"逐光":   true,
		"睱日":   true,
		"弦光":   true,
		"心晴":   true,
		"匿光":   true,
		"拥雪":   true,
		"永恒":   true,
		"夜色":   true,
		"静谧":   true,
		"深林":   true,
		"雾海":   true,
		"神殿":   true,
		"深海":   true,
		"坠浪":   true,
		"点染":   true,
		"斑斓":   true,
		"碧海":   true,
		"深渊":   true,
		"掠心":   true,
		"锋尖":   true,
		"戮夜":   true,
		"寂路":   true,
		"远空":   true,
		"长昼":   true,
		"离途":   true,
		"无套装": true,
	}

	return validSetCard[r.SetCard]
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

func (r Record) ValidateNote() (bool, error) {
	if utf8.RuneCountInString(r.Note) > 20 {
		return false, fmt.Errorf("备注最长20个字: %s", r.Note)
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

	return r.validateCommon()
}

func (r Record) ValidateChampionships() (bool, error) {
	if !r.validateBuff() {
		return false, fmt.Errorf("锦标赛加成错误: %s", r.Buff)
	}

	return r.validateCommon()
}

func (r Record) GetHash() string {
	data := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s",
		r.LevelType, r.LevelNumber, r.LevelMode, r.Attack, r.HP, r.Defense, r.Matching, r.MatchingBuff,
		r.CritRate, r.CritDmg, r.EnergyRegen, r.WeakenBoost, r.OathBoost,
		r.OathRegen, r.Partner, r.SetCard, r.Stage, r.Weapon, r.Buff, r.TotalLevel,
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
