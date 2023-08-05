package infra

import (
	"encoding/json"
	"strconv"
	"wfs/backend/apperr"
	"wfs/backend/domain"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
)

const (
	NumbersAvgDamage = "average_damage_dealt"
	NumbersAvgFrags  = "average_frags"
	NumbersWinrate   = "win_rate"
)

var (
	errNoTimeKey = errors.New("no time key")
	errNoDataKey = errors.New("no data key")
)

type Numbers struct {
	config RequestConfig
}

func NewNumbers(config RequestConfig) *Numbers {
	return &Numbers{config: config}
}

func (n *Numbers) ExpectedStats() (domain.NSExpectedStats, error) {
	var result domain.NSExpectedStats

	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), n.config.Retry)
	operation := func() (APIResponse[any], error) {
		return getRequest[any](n.config.URL, map[string]string{}, n.config.Retry)
	}

	response, err := backoff.RetryWithData(operation, b)
	if err != nil {
		return result, err
	}

	result, err = parse([]byte(response.BodyString))
	if err != nil {
		return result, err
	}

	return result, nil
}

func parse(body []byte) (domain.NSExpectedStats, error) {
	var result domain.NSExpectedStats

	root := make(map[string]interface{})
	if err := json.Unmarshal(body, &root); err != nil {
		return result, apperr.New(apperr.ErrNumbersAPI, err)
	}

	time, ok := root["time"].(float64)
	if !ok {
		return result, apperr.New(apperr.ErrNumbersAPI, errNoTimeKey)
	}
	data, ok := root["data"].(map[string]interface{})
	if !ok {
		return result, apperr.New(apperr.ErrNumbersAPI, errNoDataKey)
	}

	esd := make(map[int]domain.NSExpectedStatsData)
	for key, value := range data {
		shipID, err := strconv.Atoi(key)
		if err != nil {
			continue
		}

		values, ok := value.(map[string]interface{})
		if !ok {
			continue
		}

		damage, ok := values[NumbersAvgDamage].(float64)
		if !ok {
			continue
		}

		frags, ok := values[NumbersAvgFrags].(float64)
		if !ok {
			continue
		}

		wr, ok := values[NumbersWinrate].(float64)
		if !ok {
			continue
		}

		esd[shipID] = domain.NSExpectedStatsData{
			AverageDamageDealt: damage,
			AverageFrags:       frags,
			WinRate:            wr,
		}
	}

	result = domain.NSExpectedStats{
		Time: int(time),
		Data: esd,
	}

	return result, nil
}
