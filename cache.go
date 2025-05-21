package main

import "time"

var (
	cachedABData      [][]interface{}
	cachedABTimestamp time.Time
	cacheDuration     = 5 * time.Minute // 缓存有效期
)

func getCachedABData() ([][]interface{}, error) {
	if time.Since(cachedABTimestamp) < cacheDuration && cachedABData != nil {
		return cachedABData, nil
	}

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, sheetName+"!A2:B").Do()
	if err != nil {
		return nil, err
	}

	cachedABData = resp.Values
	cachedABTimestamp = time.Now()
	return cachedABData, nil
}
