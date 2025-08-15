package sheet_clients

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"

	"lysk-battle-record/internal/models"
)

type RecordSheetClient interface {
	FetchAllSheetData() ([]models.Record, error)
	ProcessRecord(record models.Record) (*models.Record, error)
	UpdateRecord(record models.Record) error
	DeleteRecord(record models.Record) error
	GetType() string
}

type RecordSheetClientImpl struct {
	sheetId   string
	sheetName string
	srv       *sheets.Service
}

func NewRecordSheetClient(sheetId, sheetName string) *RecordSheetClientImpl {
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
		return &RecordSheetClientImpl{
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
	return &RecordSheetClientImpl{
		srv:       srv,
		sheetId:   sheetId,
		sheetName: sheetName,
	}
}

func (c *RecordSheetClientImpl) FetchAllSheetData() ([]models.Record, error) {
	header, err := c.srv.Spreadsheets.Values.Get(c.sheetId, c.sheetName+"!A1:Z").Do()
	if err != nil {
		return nil, err
	}

	headerIndexMap := make(map[string]int)
	for i, h := range header.Values[0] {
		if hStr, ok := h.(string); ok {
			headerIndexMap[hStr] = i
		}
	}

	resp, err := c.srv.Spreadsheets.Values.Get(c.sheetId, c.sheetName+"!A2:Z").Do()
	if err != nil {
		return nil, err
	}

	var records []models.Record
	for i, row := range resp.Values {
		r := models.Record{}

		r.RowNumber = i + 2
		r.LevelType = c.getValue(row, headerIndexMap, "关卡")
		r.LevelNumber = c.getValue(row, headerIndexMap, "关数")
		r.LevelMode = c.getValue(row, headerIndexMap, "模式")
		r.Attack = c.getValue(row, headerIndexMap, "攻击")
		r.HP = c.getValue(row, headerIndexMap, "生命")
		r.Defense = c.getValue(row, headerIndexMap, "防御")
		r.Matching = c.getValue(row, headerIndexMap, "对谱")
		r.MatchingBuff = c.getValue(row, headerIndexMap, "对谱加成")
		r.CritRate = c.getValue(row, headerIndexMap, "暴击")
		r.CritDmg = c.getValue(row, headerIndexMap, "暴伤")
		r.EnergyRegen = c.getValue(row, headerIndexMap, "加速回能")
		r.WeakenBoost = c.getValue(row, headerIndexMap, "虚弱增伤")
		r.OathBoost = c.getValue(row, headerIndexMap, "誓约增伤")
		r.OathRegen = c.getValue(row, headerIndexMap, "誓约回能")
		r.Partner = c.getValue(row, headerIndexMap, "搭档身份")
		r.SetCard = c.getValue(row, headerIndexMap, "日卡")
		r.Stage = c.getValue(row, headerIndexMap, "阶数")
		r.Weapon = c.getValue(row, headerIndexMap, "武器")
		r.StarRank = c.getValue(row, headerIndexMap, "星级")
		r.Buff = c.getValue(row, headerIndexMap, "加成")
		r.TotalLevel = c.getValue(row, headerIndexMap, "卡总等级")
		r.Note = c.getValue(row, headerIndexMap, "备注")
		r.Time = c.getValue(row, headerIndexMap, "时间")
		r.UserID = c.getValue(row, headerIndexMap, "用户ID")
		r.Id = c.getValue(row, headerIndexMap, "id")

		deletedStr := c.getValue(row, headerIndexMap, "deleted")
		deleted, _ := strconv.ParseBool(deletedStr)
		r.Deleted = deleted

		records = append(records, r)
	}

	return records, nil
}

func (c *RecordSheetClientImpl) getValue(row []interface{}, headerIndexMap map[string]int, key string) string {
	if index, ok := headerIndexMap[key]; ok && index < len(row) {
		return fmt.Sprint(row[index])
	}
	return ""
}

func (c *RecordSheetClientImpl) ProcessRecord(record models.Record) (*models.Record, error) {
	record.Id = uuid.New().String()
	header, err := c.srv.Spreadsheets.Values.Get(c.sheetId, c.sheetName+"!A1:Z").Do()
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
		case "星级":
			row[index] = record.StarRank
		case "加成":
			row[index] = record.Buff
		case "卡总等级":
			row[index] = record.TotalLevel
		case "备注":
			row[index] = record.Note
		case "时间":
			row[index] = record.Time
		case "用户ID":
			row[index] = record.UserID
		case "id":
			row[index] = record.Id
		default:
		}
	}

	resp, err := c.srv.Spreadsheets.Values.Append(c.sheetId, c.sheetName+"!A1", &sheets.ValueRange{
		Values: [][]interface{}{row},
	}).ValueInputOption("RAW").Do()
	if err != nil {
		// : assign record.RowNumber after appending
		logrus.Errorf("sheet %s failed to append record to Google Sheets: %v", c.sheetName, err)
		return nil, err
	}

	rowNum, err := extractRowNumber(resp.Updates.UpdatedRange)
	if err != nil {
		return nil, err
	}
	record.RowNumber = rowNum

	return &record, nil
}

func extractRowNumber(updatedRange string) (int, error) {
	re := regexp.MustCompile(`A(\d+)`)
	matches := re.FindStringSubmatch(updatedRange)
	if len(matches) < 2 {
		return 0, fmt.Errorf("could not extract row number from range: %s", updatedRange)
	}
	return strconv.Atoi(matches[1])
}

func (c *RecordSheetClientImpl) GetType() string {
	return c.sheetName
}

func (c *RecordSheetClientImpl) UpdateRecord(record models.Record) error {
	header, err := c.srv.Spreadsheets.Values.Get(c.sheetId, c.sheetName+"!A1:Z").Do()
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
		case "星级":
			row[index] = record.StarRank
		case "加成":
			row[index] = record.Buff
		case "时间":
			row[index] = record.Time
		case "卡总等级":
			row[index] = record.TotalLevel
		case "备注":
			row[index] = record.Note
		case "用户ID":
			row[index] = record.UserID
		case "id":
			row[index] = record.Id
		case "deleted":
			row[index] = record.Deleted
		default:
		}
	}

	updateRange := fmt.Sprintf("%s!A%d", c.sheetName, record.RowNumber)
	_, err = c.srv.Spreadsheets.Values.Update(c.sheetId, updateRange, &sheets.ValueRange{
		Values: [][]interface{}{row},
	}).ValueInputOption("RAW").Do()
	if err != nil {
		logrus.Errorf("sheet %s failed to update record to Google Sheets: %v", c.sheetName, err)
		return err
	}

	return nil
}

func (c *RecordSheetClientImpl) DeleteRecord(record models.Record) error {
	header, err := c.srv.Spreadsheets.Values.Get(c.sheetId, c.sheetName+"!A1:Z").Do()
	if err != nil {
		return err
	}

	headerIndexMap := make(map[string]int)
	for i, h := range header.Values[0] {
		if hStr, ok := h.(string); ok {
			headerIndexMap[hStr] = i
		}
	}

	deleteColumnIndex, ok := headerIndexMap["deleted"]
	if !ok {
		return fmt.Errorf("deleted column not found in sheet %s", c.sheetName)
	}

	updateRange := fmt.Sprintf("%s!%s%d", c.sheetName, toCharStr(deleteColumnIndex+1), record.RowNumber)
	_, err = c.srv.Spreadsheets.Values.Update(c.sheetId, updateRange, &sheets.ValueRange{
		Values: [][]interface{}{{true}},
	}).ValueInputOption("RAW").Do()
	if err != nil {
		logrus.Errorf("sheet %s failed to delete record from Google Sheets: %v", c.sheetName, err)
		return err
	}

	return nil
}

func toCharStr(i int) string {
	s := ""
	for i > 0 {
		i--
		s = string('A'+i%26) + s
		i /= 26
	}
	return s
}
