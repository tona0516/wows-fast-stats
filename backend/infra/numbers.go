package infra

import (
	"changeme/backend/vo"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type Numbers struct {
}

func (n *Numbers) ExpectedStats() (*vo.NSExpectedStats, error) {
	res, err := http.Get("https://api.wows-numbers.com/personal/rating/expected/json/")
	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	depth1 := make(map[string]interface{})
	err = json.Unmarshal(body, &depth1)
	if err != nil {
		return nil, err
	}

	time := depth1["time"].(float64)
	depth2 := depth1["data"].(map[string]interface{})
	data := make(map[int]vo.NSExpectedStatsData)
	for key, value := range depth2 {
		keyInt, err := strconv.Atoi(key)
		if err != nil {
			continue
		}

		valueMap, ok := value.(map[string]interface{})
		if !ok {
			continue
		}

		data[keyInt] = vo.NSExpectedStatsData{
			AverageDamageDealt: valueMap["average_damage_dealt"].(float64),
			AverageFrags:       valueMap["average_frags"].(float64),
			WinRate:            valueMap["win_rate"].(float64),
		}
	}

	response := vo.NSExpectedStats{
		Time: int(time),
		Data: data,
	}

	return &response, nil
}
