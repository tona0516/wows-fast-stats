package wargaming

import (
	"net/http"
	"testing"
	"wfs/backend2/testutil"

	"github.com/imroc/req/v3"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
)

const provideName = "WargamingAPIClient"
const provideAppIDName = "WargamingAppID"

func NewStubInjector(t *testing.T, statusCode int, body map[string]interface{}) *do.Injector {
	t.Helper()

	injector := do.New()
	do.ProvideNamed(injector, provideName, func(i *do.Injector) (*req.Client, error) {
		server := testutil.NewStubServer(t, statusCode, body)
		client := req.C()
		client.SetBaseURL(server.URL)

		return client, nil
	})
	do.ProvideNamedValue(injector, provideAppIDName, "test_app_id")

	return injector
}

func TestAPI_AccountInfo(t *testing.T) {
	t.Parallel()

	injector := NewStubInjector(t, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"meta": map[string]interface{}{
			"count":  1,
			"hidden": nil,
		},
		"data": map[string]interface{}{
			"2010342809": map[string]interface{}{
				"statistics": map[string]interface{}{
					"pvp": map[string]interface{}{
						"battles": 20620,
					},
				},
			},
		},
	})

	api, _ := NewAPI(injector)

	accountID := 2010342809
	result, err := api.AccountInfo([]int{accountID})

	assert.NoError(t, err)
	assert.Equal(t, uint(20620), result[accountID].Statistics.Pvp.Battles)
}

func TestAPI_AccountList(t *testing.T) {
	t.Parallel()

	injector := NewStubInjector(t, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"meta": map[string]interface{}{
			"count": 1,
		},
		"data": []interface{}{
			map[string]interface{}{
				"nickname":   "tonango",
				"account_id": 2010342809,
			},
		},
	})

	api, _ := NewAPI(injector)
	nickName := "tonango"
	result, err := api.AccountList([]string{nickName})

	assert.NoError(t, err)
	assert.Equal(t, nickName, result[0].NickName)
	assert.Equal(t, 2010342809, result[0].AccountID)
}

func TestAPI_AccountListForSearch(t *testing.T) {
	t.Parallel()

	injector := NewStubInjector(t, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"meta": map[string]interface{}{
			"count": 1,
		},
		"data": []interface{}{
			map[string]interface{}{
				"nickname":   "tonango",
				"account_id": 2010342809,
			},
		},
	})

	api, _ := NewAPI(injector)
	nickName := "tonango"
	result, err := api.AccountListForSearch(nickName)

	assert.NoError(t, err)
	assert.Equal(t, nickName, result[0].NickName)
	assert.Equal(t, 2010342809, result[0].AccountID)
}

func TestAPI_ClansAccountInfo(t *testing.T) {
	t.Parallel()

	injector := NewStubInjector(t, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"meta": map[string]interface{}{
			"count": 1,
		},
		"data": map[string]interface{}{
			"2010342809": map[string]interface{}{
				"clan_id": 2000036632,
			},
		},
	})

	api, _ := NewAPI(injector)
	accountID := 2010342809
	result, err := api.ClansAccountInfo([]int{accountID})

	assert.NoError(t, err)
	assert.Equal(t, 2000036632, result[accountID].ClanID)
}

func TestAPI_ClansInfo(t *testing.T) {
	t.Parallel()

	injector := NewStubInjector(t, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"meta": map[string]interface{}{
			"count": 1,
		},
		"data": map[string]interface{}{
			"2000036632": map[string]interface{}{
				"tag": "-K2-",
			},
		},
	})

	api, _ := NewAPI(injector)
	clanID := 2000036632
	result, err := api.ClansInfo([]int{clanID})

	assert.NoError(t, err)
	assert.Equal(t, "-K2-", result[clanID].Tag)
}

func TestAPI_ShipsStats(t *testing.T) {
	t.Parallel()

	injector := NewStubInjector(t, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"meta": map[string]interface{}{
			"count":  1,
			"hidden": nil,
		},
		"data": map[string]interface{}{
			"2010342809": []interface{}{
				map[string]interface{}{
					"ship_id": 3769513264,
					"pvp": map[string]interface{}{
						"battles": 568,
					},
				},
			},
		},
	})

	api, _ := NewAPI(injector)
	accountID := 2010342809
	result, err := api.ShipsStats(accountID)

	assert.NoError(t, err)
	assert.Equal(t, uint(568), result[accountID][0].Pvp.Battles)
	assert.Equal(t, 3769513264, result[accountID][0].ShipID)
}

func TestAPI_EncycShips(t *testing.T) {
	t.Parallel()

	injector := NewStubInjector(t, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"meta": map[string]interface{}{
			"count":      1,
			"page":       2,
			"page_total": 3,
		},
		"data": map[string]interface{}{
			"3769513264": map[string]interface{}{
				"tier":       7,
				"is_premium": true,
				"type":       "Destroyer",
				"name":       "Blyskawica",
				"nation":     "europe",
			},
		},
	})

	api, _ := NewAPI(injector)
	result, err := api.EncycShips(2)

	assert.NoError(t, err)
	assert.Equal(t, 2, result.Meta.Page)
	assert.Equal(t, 3, result.Meta.PageTotal)

	shipID := 3769513264
	assert.Equal(t, "Blyskawica", result.Data[shipID].Name)
	assert.Equal(t, "Destroyer", result.Data[shipID].Type)
	assert.Equal(t, "europe", result.Data[shipID].Nation)
	assert.Equal(t, uint(7), result.Data[shipID].Tier)
	assert.True(t, result.Data[shipID].IsPremium)
}

func TestAPI_BattleArenas(t *testing.T) {
	t.Parallel()

	injector := NewStubInjector(t, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"meta": map[string]interface{}{
			"count": 3,
		},
		"data": map[string]interface{}{
			"0": map[string]interface{}{
				"name": "大海原",
			},
			"1": map[string]interface{}{
				"name": "ソロモン諸島",
			},
			"2": map[string]interface{}{
				"name": "列島",
			},
		},
	})

	api, _ := NewAPI(injector)
	result, err := api.BattleArenas()

	assert.NoError(t, err)
	assert.Equal(t, "大海原", result[0].Name)
	assert.Equal(t, "ソロモン諸島", result[1].Name)
	assert.Equal(t, "列島", result[2].Name)
}

func TestAPI_BattleTypes(t *testing.T) {
	t.Parallel()

	injector := NewStubInjector(t, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"meta": map[string]interface{}{
			"count": 3,
		},
		"data": map[string]interface{}{
			"CLAN": map[string]interface{}{
				"name": "クラン戦",
			},
			"PVP": map[string]interface{}{
				"name": "ランダム戦",
			},
			"RANKED": map[string]interface{}{
				"name": "ランク戦",
			},
		},
	})

	api, _ := NewAPI(injector)
	result, err := api.BattleTypes()

	assert.NoError(t, err)
	assert.Equal(t, "クラン戦", result["CLAN"].Name)
	assert.Equal(t, "ランダム戦", result["PVP"].Name)
	assert.Equal(t, "ランク戦", result["RANKED"].Name)
}

func TestAPI_GameVersion(t *testing.T) {
	t.Parallel()

	injector := NewStubInjector(t, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"meta": map[string]interface{}{
			"count": 1,
		},
		"data": map[string]interface{}{
			"game_version": "14.3.0",
		},
	})

	api, _ := NewAPI(injector)
	result, err := api.GameVersion()

	assert.NoError(t, err)
	assert.Equal(t, "14.3.0", result)
}
