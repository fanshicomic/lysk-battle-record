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
	header, err := c.srv.Spreadsheets.Values.Get(c.sheetId, c.sheetName+"!A1:Q").Do()
	if err != nil {
		return nil, err
	}

	headerIndexMap := make(map[string]int)
	for i, h := range header.Values[0] {
		if hStr, ok := h.(string); ok {
			headerIndexMap[hStr] = i
		} else {
			logrus.Info("header %v is not a string, skipping", h)
		}
	}

	resp, err := c.srv.Spreadsheets.Values.Get(c.sheetId, c.sheetName+"!A2:Q").Do()
	if err != nil {
		return nil, err
	}

	var records []Record
	for _, row := range resp.Values {
		r := Record{}

		r.LevelType = fmt.Sprint(row[headerIndexMap["关卡"]])
		r.LevelNumber = fmt.Sprint(row[headerIndexMap["关数"]])
		r.Attack = fmt.Sprint(row[headerIndexMap["攻击"]])
		r.HP = fmt.Sprint(row[headerIndexMap["生命"]])
		r.Defense = fmt.Sprint(row[headerIndexMap["防御"]])
		r.Matching = fmt.Sprint(row[headerIndexMap["对谱"]])
		r.CritRate = fmt.Sprint(row[headerIndexMap["暴击"]])
		r.CritDmg = fmt.Sprint(row[headerIndexMap["暴伤"]])
		r.EnergyRegen = fmt.Sprint(row[headerIndexMap["加速回能"]])
		r.WeakenBoost = fmt.Sprint(row[headerIndexMap["虚弱增伤"]])
		r.OathBoost = fmt.Sprint(row[headerIndexMap["誓约增伤"]])
		r.OathRegen = fmt.Sprint(row[headerIndexMap["誓约回能"]])
		r.Partner = fmt.Sprint(row[headerIndexMap["搭档身份"]])
		r.SetCard = fmt.Sprint(row[headerIndexMap["日卡"]])
		r.Stage = fmt.Sprint(row[headerIndexMap["阶数"]])
		r.Weapon = fmt.Sprint(row[headerIndexMap["武器"]])
		r.Time = fmt.Sprint(row[headerIndexMap["时间"]])

		records = append(records, r)
	}

	return records, nil
}

func (c *GoogleSheetClientImpl) ProcessRecord(record Record) error {
	keys := []string{"关卡", "关数", "攻击", "防御", "生命", "对谱", "暴击", "暴伤", "加速回能", "虚弱增伤", "誓约增伤", "誓约回能", "搭档身份", "日卡", "阶数", "武器", "时间"}
	row := make([]interface{}, 0)
	for _, key := range keys {
		switch key {
		case "关卡":
			row = append(row, record.LevelType)
		case "关数":
			row = append(row, record.LevelNumber)
		case "攻击":
			row = append(row, record.Attack)
		case "生命":
			row = append(row, record.HP)
		case "防御":
			row = append(row, record.Defense)
		case "对谱":
			row = append(row, record.Matching)
		case "暴击":
			row = append(row, record.CritRate)
		case "暴伤":
			row = append(row, record.CritDmg)
		case "加速回能":
			row = append(row, record.EnergyRegen)
		case "虚弱增伤":
			row = append(row, record.WeakenBoost)
		case "誓约增伤":
			row = append(row, record.OathBoost)
		case "誓约回能":
			row = append(row, record.OathRegen)
		case "搭档身份":
			row = append(row, record.Partner)
		case "日卡":
			row = append(row, record.SetCard)
		case "阶数":
			row = append(row, record.Stage)
		case "武器":
			row = append(row, record.Weapon)
		case "时间":
			row = append(row, record.Time)
		default:
		}
	}

	_, err := c.srv.Spreadsheets.Values.Append(c.sheetId, c.sheetName+"!A1", &sheets.ValueRange{
		Values: [][]interface{}{row},
	}).ValueInputOption("RAW").Do()
	if err != nil {
		logrus.Errorf("failed to append record to Google Sheets: %v", err)
		return err
	}

	return nil
}

func (c *GoogleSheetClientImpl) GetType() string {
	return c.sheetName
}
