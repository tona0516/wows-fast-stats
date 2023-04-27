package vo

type Teams []Team

func (t Teams) Comparition() Comparision {
    if len(t) < 2 {
        return Comparision{}
    }

    friendTeamAvg := t[0].Players.TeamAverage()
    enemyTeamAvg := t[1].Players.TeamAverage()

    friendShipStats := friendTeamAvg.ShipStats
    enemyShipStas := enemyTeamAvg.ShipStats
    friendOverallStats := friendTeamAvg.PlayerStats
    enemyOverallStats := enemyTeamAvg.PlayerStats

    return Comparision{
        Ship: ShipComp{
            PR: Between{
                Friend: friendShipStats.PR,
                Enemy: enemyShipStas.PR,
                Diff: friendShipStats.PR - enemyShipStas.PR,
            },
            Damage: Between{
                Friend: friendShipStats.Damage,
                Enemy: enemyShipStas.Damage,
                Diff: friendShipStats.Damage - enemyShipStas.Damage,
            },
            WinRate: Between{
                Friend: friendShipStats.WinRate,
                Enemy: enemyShipStas.WinRate,
                Diff: friendShipStats.WinRate -enemyShipStas.WinRate,
            },
            KdRate: Between{
                Friend: friendShipStats.KdRate,
                Enemy: enemyShipStas.KdRate,
                Diff: friendShipStats.KdRate - enemyShipStas.KdRate,
            },
        },
        Overall: OverallComp{
            Damage: Between{
                Friend: friendOverallStats.Damage,
                Enemy:enemyOverallStats.Damage,
                Diff: friendOverallStats.Damage - enemyOverallStats.Damage,
            },
            WinRate: Between{
                Friend: friendOverallStats.WinRate,
                Enemy: enemyOverallStats.WinRate,
                Diff: friendOverallStats.WinRate - enemyOverallStats.WinRate,
            },
            KdRate: Between{
                Friend: friendOverallStats.KdRate,
                Enemy: enemyOverallStats.KdRate,
                Diff: friendOverallStats.KdRate - enemyOverallStats.KdRate,
            },
        },
    }
}
