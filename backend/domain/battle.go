package domain

import "fmt"

const SampleTeamLen = 8

type Battle struct {
	Meta  Meta   `json:"meta"`
	Teams []Team `json:"teams"`
}

type Meta struct {
	Unixtime int64  `json:"unixtime"`
	Arena    string `json:"arena"`
	Type     string `json:"type"`
	OwnShip  string `json:"own_ship"`
}

type Team struct {
	Players Players `json:"players"`
	Name    string  `json:"name"`
}

type samplePlayer struct {
	tier        uint
	shipType    ShipType
	pr          float64
	damageRatio float64
	winRate     float64
}

func newSamplePlayer(
	tier uint,
	shipType ShipType,
	pr float64,
	damageRatio float64,
	winRate float64,
) samplePlayer {
	return samplePlayer{
		tier:        tier,
		shipType:    shipType,
		pr:          pr,
		damageRatio: damageRatio,
		winRate:     winRate,
	}
}

func SampleTeams() []Team {
	samplePlayers := []samplePlayer{
		newSamplePlayer(11, CV, 0, 0, 0),
		newSamplePlayer(10, BB, 750, 0.6, 47),
		newSamplePlayer(9, BB, 1100, 0.8, 50),
		newSamplePlayer(8, CL, 1350, 1.0, 52),
		newSamplePlayer(7, CL, 1550, 1.2, 54),
		newSamplePlayer(6, DD, 1750, 1.4, 56),
		newSamplePlayer(5, DD, 2100, 1.5, 60),
		newSamplePlayer(4, SS, 2450, 1.6, 65),
	}
	players := make([]Player, len(samplePlayers))

	for i, p := range samplePlayers {
		playerInfo := PlayerInfo{
			ID:   1,
			Name: fmt.Sprintf("player_name%d", i+1),
			Clan: Clan{
				Tag: "TEST",
			},
		}
		shipInfo := ShipInfo{
			Name:      "Test Ship",
			Nation:    "japan",
			Tier:      p.tier,
			Type:      p.shipType,
			AvgDamage: 10000,
		}
		shipStats := ShipStats{
			Battles:            10,
			Damage:             10000 * p.damageRatio,
			WinRate:            p.winRate,
			WinSurvivedRate:    50,
			LoseSurvivedRate:   50,
			KdRate:             1,
			Kill:               1,
			Exp:                1000,
			MainBatteryHitRate: 50,
			TorpedoesHitRate:   5,
			PlanesKilled:       5,
			PR:                 p.pr,
		}
		overallStats := OverallStats{
			Battles:          10,
			Damage:           10000 * p.damageRatio,
			WinRate:          p.winRate,
			WinSurvivedRate:  50,
			LoseSurvivedRate: 50,
			KdRate:           1,
			Kill:             1,
			Exp:              1000,
			AvgTier:          5,
			UsingShipTypeRate: ShipTypeGroup{
				SS: 20,
				DD: 20,
				CL: 20,
				BB: 20,
				CV: 20,
			},
			UsingTierRate: TierGroup{
				Low:    33.3,
				Middle: 33.3,
				High:   33.4,
			},
		}
		players[i] = Player{
			PlayerInfo: playerInfo,
			ShipInfo:   shipInfo,
			PvPSolo: PlayerStats{
				ShipStats:    shipStats,
				OverallStats: overallStats,
			},
			PvPAll: PlayerStats{
				ShipStats:    shipStats,
				OverallStats: overallStats,
			},
		}
	}

	return []Team{
		{
			Players: players,
		},
	}
}
