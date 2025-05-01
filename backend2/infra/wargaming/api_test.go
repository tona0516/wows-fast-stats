package wargaming

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imroc/req/v3"
	"github.com/stretchr/testify/assert"
)

func newMockServer(responseBody string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responseBody))
	}))
}

func newMockClient(server *httptest.Server) *req.Client {
	client := req.C()
	client.SetBaseURL(server.URL)

	return client
}

func TestAPI_AccountInfo(t *testing.T) {
	t.Parallel()

	server := newMockServer(`{
  "status":"ok",
  "meta":{
    "count":1,
    "hidden":null
  },
  "data":{
    "2010342809":{
      "last_battle_time":1744454827,
      "account_id":2010342809,
      "leveling_tier":17,
      "created_at":1447501603,
      "leveling_points":25859,
      "updated_at":1744455875,
      "private":null,
      "hidden_profile":false,
      "logout_at":1744455864,
      "karma":null,
      "statistics":{
        "distance":1207086,
        "battles":25266,
        "pvp":{
          "max_xp":5808,
          "damage_to_buildings":136340,
          "main_battery":{
            "max_frags_battle":10,
            "frags":16621,
            "hits":1523979,
            "max_frags_ship_id":4188977104,
            "shots":3857966
          },
          "max_ships_spotted_ship_id":3542071248,
          "max_damage_scouting":288762,
          "art_agro":4000000001,
          "max_xp_ship_id":4178491216,
          "ships_spotted":20445,
          "second_battery":{
            "max_frags_battle":2,
            "frags":209,
            "hits":61324,
            "max_frags_ship_id":4284429648,
            "shots":336160
          },
          "max_frags_ship_id":4188977104,
          "xp":39108833,
          "survived_battles":11129,
          "dropped_capture_points":196,
          "max_damage_dealt_to_buildings":81900,
          "torpedo_agro":1534193941,
          "draws":20,
          "control_captured_points":384847,
          "battles_since_510":18682,
          "max_total_agro_ship_id":4179572176,
          "planes_killed":80477,
          "battles":20581,
          "max_ships_spotted":12,
          "max_suppressions_ship_id":null,
          "survived_wins":9139,
          "frags":22040,
          "damage_scouting":565697645,
          "max_total_agro":5553237,
          "max_frags_battle":10,
          "capture_points":0,
          "ramming":{
            "max_frags_battle":1,
            "frags":98,
            "max_frags_ship_id":4284429648
          },
          "suppressions_count":0,
          "max_suppressions_count":0,
          "torpedoes":{
            "max_frags_battle":5,
            "frags":1972,
            "hits":12837,
            "max_frags_ship_id":4181636304,
            "shots":179967
          },
          "max_planes_killed_ship_id":4179506480,
          "aircraft":{
            "max_frags_battle":5,
            "frags":729,
            "max_frags_ship_id":4183799248
          },
          "team_capture_points":2243726,
          "control_dropped_points":184283,
          "max_damage_dealt":294322,
          "max_damage_dealt_to_buildings_ship_id":4281284304,
          "max_damage_dealt_ship_id":3760141776,
          "wins":12640,
          "losses":7921,
          "damage_dealt":1414957613,
          "max_planes_killed":62,
          "max_scouting_damage_ship_id":4179506480,
          "team_dropped_capture_points":1180957,
          "battles_since_512":18184
        }
      },
      "nickname":"tonango",
      "stats_updated_at":1744455875
    }
  }
}`)
	defer server.Close()

	mockClient := newMockClient(server)
	api := NewAPI(mockClient, "test_app_id")

	accountID := 2010342809
	result, err := api.AccountInfo([]int{accountID})

	assert.NoError(t, err)
	assert.NotEmpty(t, result[accountID])
}

func TestAPI_AccountList(t *testing.T) {
	t.Parallel()

	server := newMockServer(`{
  "status":"ok",
  "meta":{
    "count":1
  },
  "data":[
    {
      "nickname":"tonango",
      "account_id":2010342809
    }
  ]
}`)
	defer server.Close()

	mockClient := newMockClient(server)
	api := NewAPI(mockClient, "test_app_id")

	nickName := "tonango"
	result, err := api.AccountList([]string{nickName})

	assert.NoError(t, err)
	assert.Equal(t, nickName, result[0].NickName)
	assert.Equal(t, 2010342809, result[0].AccountID)
}

func TestAPI_AccountListForSearch(t *testing.T) {
	t.Parallel()

	server := newMockServer(`{
  "status":"ok",
  "meta":{
    "count":1
  },
  "data":[
    {
      "nickname":"tonango",
      "account_id":2010342809
    }
  ]
}`)
	defer server.Close()

	mockClient := newMockClient(server)
	api := NewAPI(mockClient, "test_app_id")

	nickName := "tonango"
	result, err := api.AccountListForSearch(nickName)

	assert.NoError(t, err)
	assert.Equal(t, nickName, result[0].NickName)
	assert.Equal(t, 2010342809, result[0].AccountID)
}

func TestAPI_ClansAccountInfo(t *testing.T) {
	t.Parallel()

	server := newMockServer(`{
  "status":"ok",
  "meta":{
    "count":1
  },
  "data":{
    "2010342809":{
      "role":"executive_officer",
      "clan_id":2000036632,
      "joined_at":1701953479,
      "account_id":2010342809,
      "account_name":"tonango"
    }
  }
}`)
	defer server.Close()

	mockClient := newMockClient(server)
	api := NewAPI(mockClient, "test_app_id")

	accountID := 2010342809
	result, err := api.ClansAccountInfo([]int{accountID})

	assert.NoError(t, err)
	assert.Equal(t, 2000036632, result[accountID].ClanID)
}

func TestAPI_ClansInfo(t *testing.T) {
	t.Parallel()

	//nolint:lll
	server := newMockServer(`{
  "status":"ok",
  "meta":{
    "count":1
  },
  "data":{
    "2000036632":{
      "members_count":11,
      "name":"神風-s",
      "creator_name":"myouko02",
      "clan_id":2000036632,
      "created_at":1701159041,
      "updated_at":1741329671,
      "leader_name":"myouko02",
      "members_ids":[
        2008949743,
        2010342809,
        2011651752,
        2020139332,
        2026082124,
        2038889037,
        2041245324,
        2041825589,
        2042655198,
        2049607306,
        2049770926
      ],
      "creator_id":2038889037,
      "tag":"-K2-",
      "old_name":null,
      "is_clan_disbanded":false,
      "renamed_at":null,
      "old_tag":null,
      "leader_id":2038889037,
      "description":"『-K2-』神風‐sではTyphoonリーグを目指すクランとなります。\nクラン加入については下記のDiscordの招待URLを通して面接申請をお願いします。\nもしDiscordをお持ちでない方は、ゲーム内チャットにてmyouko02もしくは、gaku0083ご連絡してください。"
    }
  }
}`)
	defer server.Close()

	mockClient := newMockClient(server)
	api := NewAPI(mockClient, "test_app_id")

	clanID := 2000036632
	result, err := api.ClansInfo([]int{clanID})

	assert.NoError(t, err)
	assert.Equal(t, "-K2-", result[clanID].Tag)
	assert.NotEmpty(t, result[clanID].Description)
}

func TestAPI_ShipsStats(t *testing.T) {
	t.Parallel()

	server := newMockServer(`{
  "status":"ok",
  "meta":{
    "count":1,
    "hidden":null
  },
  "data":{
    "2010342809":[
      {
        "pvp":{
          "max_xp":3474,
          "damage_to_buildings":0,
          "main_battery":{
            "max_frags_battle":6,
            "frags":458,
            "hits":59809,
            "shots":132967
          },
          "suppressions_count":0,
          "max_damage_scouting":128174,
          "art_agro":283622650,
          "ships_spotted":1294,
          "second_battery":{
            "max_frags_battle":0,
            "frags":0,
            "hits":0,
            "shots":0
          },
          "xp":933415,
          "survived_battles":345,
          "dropped_capture_points":8290,
          "max_damage_dealt_to_buildings":0,
          "torpedo_agro":47804110,
          "draws":1,
          "battles_since_510":568,
          "planes_killed":218,
          "battles":568,
          "max_ships_spotted":9,
          "team_capture_points":64110,
          "frags":650,
          "damage_scouting":15582844,
          "max_total_agro":1844261,
          "max_frags_battle":6,
          "capture_points":21473,
          "ramming":{
            "max_frags_battle":0,
            "frags":0
          },
          "torpedoes":{
            "max_frags_battle":2,
            "frags":80,
            "hits":346,
            "shots":8595
          },
          "aircraft":{
            "max_frags_battle":0,
            "frags":0
          },
          "survived_wins":276,
          "max_damage_dealt":137327,
          "wins":366,
          "losses":201,
          "damage_dealt":24062152,
          "max_planes_killed":13,
          "max_suppressions_count":0,
          "team_dropped_capture_points":33513,
          "battles_since_512":568
        },
        "last_battle_time":1729775719,
        "account_id":2010342809,
        "distance":41927,
        "updated_at":1729776457,
        "battles":693,
        "ship_id":3769513264,
        "private":null
      }
    ]
  }
}`)
	defer server.Close()

	mockClient := newMockClient(server)
	api := NewAPI(mockClient, "test_app_id")

	accountID := 2010342809
	result, err := api.ShipsStats(accountID)

	assert.NoError(t, err)
	assert.NotEmpty(t, result[accountID][0].Pvp)
	assert.Equal(t, 3769513264, result[accountID][0].ShipID)
}

func TestAPI_EncycShips(t *testing.T) {
	t.Parallel()

	server := newMockServer(`{
  "status":"ok",
  "meta":{
    "count":1,
    "page_total":1,
    "total":1,
    "limit":100,
    "page":1
  },
  "data":{
    "3769513264":{
      "tier":7,
      "is_premium":true,
      "type":"Destroyer",
      "name":"Blyskawica",
      "nation":"europe"
    }
  }
}`)
	defer server.Close()

	mockClient := newMockClient(server)
	api := NewAPI(mockClient, "test_app_id")

	result, totalPages, err := api.EncycShips(1)

	assert.NoError(t, err)
	assert.Equal(t, 1, totalPages)
	assert.Equal(t, "Blyskawica", result.Data[3769513264].Name)
	assert.Equal(t, "Destroyer", result.Data[3769513264].Type)
	assert.Equal(t, "europe", result.Data[3769513264].Nation)
	assert.Equal(t, uint(7), result.Data[3769513264].Tier)
	assert.True(t, result.Data[3769513264].IsPremium)
}

func TestAPI_BattleArenas(t *testing.T) {
	t.Parallel()

	server := newMockServer(`{
  "status":"ok",
  "meta":{
    "count":55
  },
  "data":{
    "0":{
      "name":"大海原"
    },
    "1":{
      "name":"ソロモン諸島"
    },
    "2":{
      "name":"列島"
    },
    "3":{
      "name":"リング"
    },
    "4":{
      "name":"海峡"
    },
    "5":{
      "name":"ビッグレース"
    },
    "6":{
      "name":"新たなる夜明け"
    },
    "7":{
      "name":"大西洋"
    },
    "8":{
      "name":"北方"
    },
    "9":{
      "name":"ホットスポット"
    },
    "10":{
      "name":"断層線"
    },
    "11":{
      "name":"氷の群島"
    },
    "12":{
      "name":"罠"
    },
    "13":{
      "name":"二人の兄弟"
    },
    "14":{
      "name":"火の地"
    },
    "15":{
      "name":"破片"
    },
    "16":{
      "name":"幸運の海"
    },
    "17":{
      "name":"砂漠の涙"
    },
    "18":{
      "name":"極地"
    },
    "19":{
      "name":"群島"
    },
    "20":{
      "name":"北極光"
    },
    "21":{
      "name":"山岳地帯"
    },
    "22":{
      "name":"粉砕"
    },
    "23":{
      "name":"沖縄"
    },
    "24":{
      "name":"トライデント"
    },
    "25":{
      "name":"隣接勢力"
    },
    "26":{
      "name":"戦士の道"
    },
    "27":{
      "name":"ループ"
    },
    "28":{
      "name":"河口"
    },
    "29":{
      "name":"眠れる巨人"
    },
    "30":{
      "name":"安息の地"
    },
    "31":{
      "name":"ギリシャ"
    },
    "32":{
      "name":"クラッシュゾーン α"
    },
    "33":{
      "name":"北方海域"
    },
    "34":{
      "name":"フェロー諸島"
    },
    "35":{
      "name":"セーシェル"
    },
    "36":{
      "name":"反撃"
    },
    "37":{
      "name":"キラー・ホエール"
    },
    "38":{
      "name":"ニューポート海軍基地"
    },
    "39":{
      "name":"迷宮"
    },
    "41":{
      "name":"ルーアン環礁"
    },
    "42":{
      "name":"スンダ列島"
    },
    "43":{
      "name":"エルメス"
    },
    "44":{
      "name":"エンプレス・オーガスタ湾"
    },
    "45":{
      "name":"ノルマンディー海岸"
    },
    "46":{
      "name":"ノルマンディー海岸"
    },
    "47":{
      "name":"水没都市"
    },
    "48":{
      "name":"Transylvania を救え"
    },
    "49":{
      "name":"ドラゴン湾"
    },
    "53":{
      "name":"ポリゴン"
    },
    "54":{
      "name":"お風呂"
    },
    "55":{
      "name":"カラフル・アイランズ"
    },
    "56":{
      "name":"ヴァルカン星"
    },
    "57":{
      "name":"ヴァルカン星"
    },
    "58":{
      "name":"ヴァルカン星"
    }
  }
}`)
	defer server.Close()

	mockClient := newMockClient(server)
	api := NewAPI(mockClient, "test_app_id")

	result, err := api.BattleArenas()

	assert.NoError(t, err)
	assert.Len(t, result, 55)
	assert.Equal(t, "大海原", result[0].Name)
}

func TestAPI_BattleTypes(t *testing.T) {
	t.Parallel()

	server := newMockServer(`{
  "status":"ok",
  "meta":{
    "count":8
  },
  "data":{
    "CLAN":{
      "name":"クラン戦"
    },
    "PVP":{
      "name":"ランダム戦"
    },
    "BRAWL":{
      "name":"闘争"
    },
    "RANKED":{
      "name":"ランク戦"
    },
    "PVE":{
      "name":"オペレーション"
    },
    "PVE_PREMADE":{
      "name":"オペレーション"
    },
    "COOPERATIVE":{
      "name":"Co-op 戦"
    },
    "EVENT":{
      "name":"期間限定戦闘タイプ"
    }
  }
}`)
	defer server.Close()

	mockClient := newMockClient(server)
	api := NewAPI(mockClient, "test_app_id")

	result, err := api.BattleTypes()

	assert.NoError(t, err)
	assert.Len(t, result, 8)
	assert.Equal(t, "クラン戦", result["CLAN"].Name)
}

func TestAPI_GameVersion(t *testing.T) {
	t.Parallel()

	server := newMockServer(`{
  "status":"ok",
  "meta":{
    "count":1
  },
  "data":{
    "game_version":"14.3.0"
  }
}`)
	defer server.Close()

	mockClient := newMockClient(server)
	api := NewAPI(mockClient, "test_app_id")

	result, err := api.GameVersion()

	assert.NoError(t, err)
	assert.Equal(t, "14.3.0", result)
}
