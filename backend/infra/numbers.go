package infra

import (
	"changeme/backend/apperr"
	"changeme/backend/vo"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
)

type Numbers struct{}

func (n *Numbers) ExpectedStats() (vo.NSExpectedStats, error) {
	var result vo.NSExpectedStats

	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)
	body, err := backoff.RetryWithData(fetch, b)
	if err != nil {
		return result, err
	}

	result, err = parse(body)
	if err != nil {
		return result, err
	}

	return result, nil
}

func fetch() ([]byte, error) {
	errDetail := apperr.Ns.Req

	res, err := http.Get("https://api.wows-numbers.com/personal/rating/expected/json/")
	if err != nil {
		return []byte{}, errors.WithStack(errDetail.WithRaw(err))
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, errors.WithStack(errDetail.WithRaw(err))
	}

	return body, nil
}

func parse(body []byte) (vo.NSExpectedStats, error) {
	errDetail := apperr.Ns.Parse

	var result vo.NSExpectedStats

	depth1 := make(map[string]interface{})
	if err := json.Unmarshal(body, &depth1); err != nil {
		return result, errors.WithStack(errDetail)
	}

	time, ok := depth1["time"].(float64)
	if !ok {
		return result, errors.WithStack(errDetail.WithRaw(apperr.ErrNoTimeKey))
	}
	depth2, ok := depth1["data"].(map[string]interface{})
	if !ok {
		return result, errors.WithStack(errDetail.WithRaw(apperr.ErrNoDataKey))
	}

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

		avgDmg, ok := valueMap["average_damage_dealt"].(float64)
		if !ok {
			continue
		}

		avgFlgs, ok := valueMap["average_frags"].(float64)
		if !ok {
			continue
		}

		winRate, ok := valueMap["win_rate"].(float64)
		if !ok {
			continue
		}

		data[keyInt] = vo.NSExpectedStatsData{
			AverageDamageDealt: avgDmg,
			AverageFrags:       avgFlgs,
			WinRate:            winRate,
		}
	}

	result = vo.NSExpectedStats{
		Time: int(time),
		Data: data,
	}

	return result, nil
}
