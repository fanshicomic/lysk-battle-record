package sheet_clients

import (
	"context"
	"fmt"
	"os"

	"lysk-battle-record/internal/models"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

type UserSheetClient interface {
	FetchAllSheetData() ([]models.User, error)
	ProcessUser(user models.User) (*models.User, error)
	UpdateUser(user models.User) error
	GetType() string
}

type UserSheetClientImpl struct {
	sheetId   string
	sheetName string
	srv       *sheets.Service
}

func NewUserSheetClient(sheetId, sheetName string) *UserSheetClientImpl {
	ctx := context.Background()
	var srv *sheets.Service

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
		return &UserSheetClientImpl{
			srv:       srv,
			sheetId:   sheetId,
			sheetName: sheetName,
		}
	}

	client, err := google.DefaultClient(ctx, sheets.SpreadsheetsScope)
	if err != nil {
		logrus.Fatalf("%s failed to fetch service account credential: %v", sheetName, err)
	}

	srv, err = sheets.New(client)
	if err != nil {
		logrus.Fatalf("%s failed to init Sheets client with service account credential: %v", sheetName, err)
	}

	logrus.Infof("%s using default client (Cloud Run) to init Sheets client", sheetName)
	return &UserSheetClientImpl{
		srv:       srv,
		sheetId:   sheetId,
		sheetName: sheetName,
	}
}

func (c *UserSheetClientImpl) FetchAllSheetData() ([]models.User, error) {
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

	var users []models.User
	for i, row := range resp.Values {
		u := models.User{}

		u.RowNumber = i + 2
		u.ID = c.getValue(row, headerIndexMap, "id")
		u.Nickname = c.getValue(row, headerIndexMap, "nickname")

		users = append(users, u)
	}

	return users, nil
}

func (c *UserSheetClientImpl) getValue(row []interface{}, headerIndexMap map[string]int, key string) string {
	if index, ok := headerIndexMap[key]; ok && index < len(row) {
		return fmt.Sprint(row[index])
	}
	return ""
}

func (c *UserSheetClientImpl) ProcessUser(user models.User) (*models.User, error) {
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
		case "id":
			row[index] = user.ID
		case "nickname":
			row[index] = user.Nickname
		default:
		}
	}

	resp, err := c.srv.Spreadsheets.Values.Append(c.sheetId, c.sheetName+"!A1", &sheets.ValueRange{
		Values: [][]interface{}{row},
	}).ValueInputOption("RAW").Do()
	if err != nil {
		logrus.Errorf("sheet %s failed to append user to Google Sheets: %v", c.sheetName, err)
		return nil, err
	}

	rowNum, err := extractRowNumber(resp.Updates.UpdatedRange)
	if err != nil {
		return nil, err
	}
	user.RowNumber = rowNum

	return &user, nil
}

func (c *UserSheetClientImpl) GetType() string {
	return c.sheetName
}

func (c *UserSheetClientImpl) UpdateUser(user models.User) error {
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
		case "id":
			row[index] = user.ID
		case "nickname":
			row[index] = user.Nickname
		default:
		}
	}

	updateRange := fmt.Sprintf("%s!A%d", c.sheetName, user.RowNumber)
	_, err = c.srv.Spreadsheets.Values.Update(c.sheetId, updateRange, &sheets.ValueRange{
		Values: [][]interface{}{row},
	}).ValueInputOption("RAW").Do()
	if err != nil {
		logrus.Errorf("sheet %s failed to update user to Google Sheets: %v", c.sheetName, err)
		return err
	}

	return nil
}
