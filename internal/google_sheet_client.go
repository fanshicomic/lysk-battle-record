package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

type GoogleSheetClient interface {
	FetchAllSheetData() ([]Record, error)
	ProcessRecord(record Record) error
	GetType() string
}

type GoogleSheetClientImpl struct {
	sheetId   string
	sheetName string
	srv       *sheets.Service
}

func NewGoogleSheetClient(sheetId, sheetName string) *GoogleSheetClientImpl {
	ctx := context.Background()
	var srv *sheets.Service

	// 优先使用本地 credentials.json，如果不存在就走默认（Cloud Run）
	if b, err := os.ReadFile("credentials.json"); err == nil {
		config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsScope)
		if err != nil {
			logrus.Fatalf("%s sheet failed to load credentials.json: %v", sheetName, err)
		}

		client := config.Client(ctx)
		srv, err = sheets.New(client)
		if err != nil {
			logrus.Fatalf("%s failed to init Sheets client with credentials.json: %v", sheetName, err)
		}

		logrus.Infof("%s use credentials.json to init Sheets client", sheetName)
		return &GoogleSheetClientImpl{
			srv:       srv,
			sheetId:   sheetId,
			sheetName: sheetName,
		}
	}

	// fallback 到 Cloud Run 默认凭证
	client, err := google.DefaultClient(ctx, sheets.SpreadsheetsScope)
	if err != nil {
		logrus.Fatalf("%s failed to fetch service account credential: %v", sheetName, err)
	}

	srv, err = sheets.New(client)
	if err != nil {
		logrus.Fatalf("%s failed to init Sheets client with service account credential: %v", sheetName, err)
	}

	logrus.Infof("%s using default client (Cloud Run) to init Sheets client", sheetName)
	return &GoogleSheetClientImpl{
		srv:       srv,
		sheetId:   sheetId,
		sheetName: sheetName,
	}
}

func (c *GoogleSheetClientImpl) FetchAllSheetData() ([]Record, error) {
	header, err := c.srv.Spreadsheets.Values.Get(c.sheetId, c.sheetName+"!A1:W").Do()
	if err != nil {
		return nil, err
	}

	headerIndexMap := make(map[string]int)
	for i, h := range header.Values[0] {
		if hStr, ok := h.(string); ok {
			headerIndexMap[hStr] = i
		}
	}

	resp, err := c.srv.Spreadsheets.Values.Get(c.sheetId, c.sheetName+"!A2:W").Do()
	if err != nil {
		return nil, err
	}

	var records []Record
	for _, row := range resp.Values {
		r := Record{}

		r.LevelType = fmt.Sprint(c.getValue(row, headerIndexMap, "关卡"))
		r.LevelNumber = fmt.Sprint(c.getValue(row, headerIndexMap, "关数"))
		r.LevelMode = fmt.Sprint(c.getValue(row, headerIndexMap, "模式"))
		r.Attack = fmt.Sprint(c.getValue(row, headerIndexMap, "攻击"))
		r.HP = fmt.Sprint(c.getValue(row, headerIndexMap, "生命"))
		r.Defense = fmt.Sprint(c.getValue(row, headerIndexMap, "防御"))
		r.Matching = fmt.Sprint(c.getValue(row, headerIndexMap, "对谱"))
		r.MatchingBuff = fmt.Sprint(c.getValue(row, headerIndexMap, "对谱加成"))
		r.CritRate = fmt.Sprint(c.getValue(row, headerIndexMap, "暴击"))
		r.CritDmg = fmt.Sprint(c.getValue(row, headerIndexMap, "暴伤"))
		r.EnergyRegen = fmt.Sprint(c.getValue(row, headerIndexMap, "加速回能"))
		r.WeakenBoost = fmt.Sprint(c.getValue(row, headerIndexMap, "虚弱增伤"))
		r.OathBoost = fmt.Sprint(c.getValue(row, headerIndexMap, "誓约增伤"))
		r.OathRegen = fmt.Sprint(c.getValue(row, headerIndexMap, "誓约回能"))
		r.Partner = fmt.Sprint(c.getValue(row, headerIndexMap, "搭档身份"))
		r.SetCard = fmt.Sprint(c.getValue(row, headerIndexMap, "日卡"))
		r.Stage = fmt.Sprint(c.getValue(row, headerIndexMap, "阶数"))
		r.Weapon = fmt.Sprint(c.getValue(row, headerIndexMap, "武器"))
		r.Buff = fmt.Sprint(c.getValue(row, headerIndexMap, "加成"))
		r.Time = fmt.Sprint(c.getValue(row, headerIndexMap, "时间"))
		r.UserID = fmt.Sprint(c.getValue(row, headerIndexMap, "用户ID"))

		records = append(records, r)
	}

	return records, nil
}

func (c *GoogleSheetClientImpl) getValue(row []interface{}, headerIndexMap map[string]int, key string) interface{} {
	if index, ok := headerIndexMap[key]; ok && index < len(row) {
		return row[index]
	}
	return nil
}

func (c *GoogleSheetClientImpl) ProcessRecord(record Record) error {
	header, err := c.srv.Spreadsheets.Values.Get(c.sheetId, c.sheetName+"!A1:W").Do()
	if err != nil {
		return err
	}

	headerIndexMap := make(map[string]int)
	for i, h := range header.Values[0] {
		if hStr, ok := h.(string); ok {
			headerIndexMap[hStr] = i
		} else {
			logrus.Infof("sheet %s header %v is not a string, skipping", c.sheetName, h)
		}
	}

	row := make([]interface{}, len(headerIndexMap))
	for key, index := range headerIndexMap {
		switch key {
		case "关卡":
			row[index] = record.LevelType
		case "关数":
			row[index] = record.LevelNumber
		case "模式":
			row[index] = record.LevelMode
		case "攻击":
			row[index] = record.Attack
		case "防御":
			row[index] = record.Defense
		case "生命":
			row[index] = record.HP
		case "对谱":
			row[index] = record.Matching
		case "对谱加成":
			row[index] = record.MatchingBuff
		case "暴击":
			row[index] = record.CritRate
		case "暴伤":
			row[index] = record.CritDmg
		case "加速回能":
			row[index] = record.EnergyRegen
		case "虚弱增伤":
			row[index] = record.WeakenBoost
		case "誓约增伤":
			row[index] = record.OathBoost
		case "誓约回能":
			row[index] = record.OathRegen
		case "搭档身份":
			row[index] = record.Partner
		case "日卡":
			row[index] = record.SetCard
		case "阶数":
			row[index] = record.Stage
		case "武器":
			row[index] = record.Weapon
		case "加成":
			row[index] = record.Buff
		case "时间":
			row[index] = record.Time
		case "用户ID":
			row[index] = record.UserID
		default:
		}
	}

	_, err = c.srv.Spreadsheets.Values.Append(c.sheetId, c.sheetName+"!A1", &sheets.ValueRange{
		Values: [][]interface{}{row},
	}).ValueInputOption("RAW").Do()
	if err != nil {
		logrus.Errorf("sheet %s failed to append record to Google Sheets: %v", c.sheetName, err)
		return err
	}

	return nil
}

func (c *GoogleSheetClientImpl) GetType() string {
	return c.sheetName
}
