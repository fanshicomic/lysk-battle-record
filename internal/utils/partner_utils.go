package utils

func GetPartnerSetCardMap() map[string]map[string]bool {
	return map[string]map[string]bool{
		"沈星回": {
			"夜誓": true, "末夜": true, "逐光": true, "鎏光": true, "睱日": true, "弦光": true, "心晴": true, "匿光": true, "无套装": true,
		},
		"黎深": {
			"神谕": true, "拥雪": true, "永恒": true, "终序": true, "夜色": true, "静谧": true, "心晴": true, "深林": true, "无套装": true,
		},
		"祁煜": {
			"雾海": true, "神殿": true, "深海": true, "坠浪": true, "点染": true, "斑斓": true, "心晴": true, "碧海": true, "无套装": true,
		},
		"秦彻": {
			"猩红": true, "深渊": true, "掠心": true, "纯白": true, "锋尖": true, "戮夜": true, "无套装": true,
		},
		"夏以昼": {
			"寂路": true, "远空": true, "长昼": true, "离途": true, "无套装": true,
		},
	}
}

func GetPartnerCompanionMap() map[string][]string {
	return map[string][]string{
		"沈星回": {"暗蚀国王", "光猎", "逐光骑士", "遥远少年", "Evol特警", "深空猎人"},
		"黎深":   {"终末之神", "九黎司命", "永恒先知", "极地军医", "黎明抹杀者", "临空医生"},
		"祁煜":   {"利莫里亚海神", "潮汐之神", "深海潜行者", "画坛新锐", "海妖魅影", "艺术家"},
		"秦彻":   {"银翼恶魔", "深渊主宰", "无尽掠夺者", "异界来客"},
		"夏以昼": {"终极兵器X-02", "远空执舰官", "深空飞行员"},
	}
}
