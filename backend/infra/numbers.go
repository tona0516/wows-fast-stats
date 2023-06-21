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

var (
	errNoTimeKey = errors.New("no time key")
	errNoDataKey = errors.New("no data key")
)

type Numbers struct {
	URL string
}

func NewNumbers(url string) *Numbers {
	return &Numbers{
		URL: url,
	}
}

func (n *Numbers) ExpectedStats() (vo.NSExpectedStats, error) {
	var result vo.NSExpectedStats

	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)
	body, err := backoff.RetryWithData(n.fetch, b)
	if err != nil {
		return result, err
	}

	result, err = parse(body)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (n *Numbers) fetch() ([]byte, error) {
	res, err := http.Get(n.URL)
	if err != nil {
		return []byte{}, apperr.New(apperr.HTTPRequest, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, apperr.New(apperr.HTTPRequest, err)
	}

	return body, nil
}

func parse(body []byte) (vo.NSExpectedStats, error) {
	var result vo.NSExpectedStats

	depth1 := make(map[string]interface{})
	if err := json.Unmarshal(body, &depth1); err != nil {
		return result, apperr.New(apperr.NumbersAPIParse, err)
	}

	time, ok := depth1["time"].(float64)
	if !ok {
		return result, apperr.New(apperr.NumbersAPIParse, errNoTimeKey)
	}
	depth2, ok := depth1["data"].(map[string]interface{})
	if !ok {
		return result, apperr.New(apperr.NumbersAPIParse, errNoDataKey)
	}

	data := make(map[int]vo.NSExpectedStatsData)
	for key, value := range depth2 {
		shipID, err := strconv.Atoi(key)
		if err != nil {
			continue
		}

		valueMap, ok := value.(map[string]interface{})
		if !ok {
			continue
		}

		damage, ok := valueMap["average_damage_dealt"].(float64)
		if !ok {
			continue
		}

		frags, ok := valueMap["average_frags"].(float64)
		if !ok {
			continue
		}

		wr, ok := valueMap["win_rate"].(float64)
		if !ok {
			continue
		}

		data[shipID] = vo.NSExpectedStatsData{
			AverageDamageDealt: damage,
			AverageFrags:       frags,
			WinRate:            wr,
		}
	}

	result = vo.NSExpectedStats{
		Time: int(time),
		Data: data,
	}

	return result, nil
}
