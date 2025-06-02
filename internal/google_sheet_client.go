package internal

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/sheets/v4"
)

type GoogleSheetClient interface {
	FetchAllSheetData() ([]Record, error)
	ProcessRecord(record Record) error
	GetType() string
	MarkAllAsExpired() error
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
	header, err := c.srv.Spreadsheets.Values.Get(c.sheetId, c.sheetName+"!A1:R").Do()
	if err != nil {
		return nil, err
	}

	headerIndexMap := make(map[string]int)
	for i, h := range header.Values[0] {
		if hStr, ok := h.(string); ok {
			headerIndexMap[hStr] = i
		} else {
			logrus.Infof("sheet %s header %v is not a string, skipping", c.sheetName, h)
		}
	}

	resp, err := c.srv.Spreadsheets.Values.Get(c.sheetId, c.sheetName+"!A2:R").Do()
	if err != nil {
		return nil, err
	}

	var records []Record
	for _, row := range resp.Values {
		r := Record{}

		r.LevelType = fmt.Sprint(c.getValue(row, headerIndexMap, "关卡"))
		r.LevelNumber = fmt.Sprint(c.getValue(row, headerIndexMap, "关数"))
		r.Attack = fmt.Sprint(c.getValue(row, headerIndexMap, "攻击"))
		r.HP = fmt.Sprint(c.getValue(row, headerIndexMap, "生命"))
		r.Defense = fmt.Sprint(c.getValue(row, headerIndexMap, "防御"))
		r.Matching = fmt.Sprint(c.getValue(row, headerIndexMap, "对谱"))
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
		r.Buffer = fmt.Sprint(c.getValue(row, headerIndexMap, "加成"))
		r.Time = fmt.Sprint(c.getValue(row, headerIndexMap, "时间"))

		records = append(records, r)
	}

	return records, nil
}

func (c *GoogleSheetClientImpl) getValue(row []interface{}, headerIndexMap map[string]int, key string) interface{} {
	if index, ok := headerIndexMap[key]; ok && index < len(row) {
		return row[index]
	}
	logrus.Warnf("sheet %s header %s not found in row %v", c.sheetName, key, row)
	return nil
}

func (c *GoogleSheetClientImpl) ProcessRecord(record Record) error {
	keys := []string{"关卡", "关数", "攻击", "防御", "生命", "对谱", "暴击", "暴伤", "加速回能", "虚弱增伤", "誓约增伤", "誓约回能", "搭档身份", "日卡", "阶数", "武器", "加成", "时间"}
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
		case "加成":
			row = append(row, record.Buffer)
		case "时间":
			row = append(row, record.Time)
		default:
		}
	}

	_, err := c.srv.Spreadsheets.Values.Append(c.sheetId, c.sheetName+"!A1", &sheets.ValueRange{
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

func (c *GoogleSheetClientImpl) MarkAllAsExpired() error {
	resp, err := c.srv.Spreadsheets.Values.Get(c.sheetId, c.sheetName+"!A1:Z").Do()
	if err != nil {
		logrus.Errorf("failed to read sheet %s with error: %v", c.sheetName, err)
		return err
	}

	if len(resp.Values) == 0 {
		logrus.Infof("%s sheet is empty", c.sheetName)
		return nil
	}

	headers := resp.Values[0]
	expiredCol := -1
	for i, h := range headers {
		if h == "已过期" {
			expiredCol = i
			break
		}
	}
	if expiredCol == -1 {
		logrus.Error("expired column not found")
		return fmt.Errorf("未找到 '已过期' 字段")
	}

	// 计算哪些需要更新
	var updates []*sheets.Request
	sheetId := int64(1105225329)
	for rowIndex, row := range resp.Values[1:] {
		// 记录从第2行开始（index 1），所以 +2 表示行号
		rowNum := rowIndex + 2

		// 如果已过期字段存在且已经是 true，跳过
		if expiredCol < len(row) && strings.TrimSpace(fmt.Sprint(row[expiredCol])) == "true" {
			continue
		}

		// 构造更新请求
		updates = append(updates, &sheets.Request{
			UpdateCells: &sheets.UpdateCellsRequest{
				Range: &sheets.GridRange{
					SheetId:          sheetId, // 你需要先获取实际的 sheetId
					StartRowIndex:    int64(rowNum - 1),
					EndRowIndex:      int64(rowNum),
					StartColumnIndex: int64(expiredCol),
					EndColumnIndex:   int64(expiredCol + 1),
				},
				Rows: []*sheets.RowData{{
					Values: []*sheets.CellData{{
						UserEnteredValue: &sheets.ExtendedValue{BoolValue: googleapi.Bool(true)},
					}},
				}},
				Fields: "userEnteredValue",
			},
		})
	}

	if len(updates) == 0 {
		logrus.Info("All records are already marked as expired")
		return nil
	}

	_, err = c.srv.Spreadsheets.BatchUpdate(c.sheetId, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: updates,
	}).Do()
	if err != nil {
		return fmt.Errorf("sheet %s failed to update for expiration: %w", c.sheetName, err)
	}

	logrus.Infof("updated %d records to expired", len(updates))
	return nil
}
