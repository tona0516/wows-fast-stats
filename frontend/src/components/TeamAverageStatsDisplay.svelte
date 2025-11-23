<script lang="ts">
  import { SHIP_TYPES } from "@libs/constants";
  import type { data } from "@wails/go/models";
  import StatCell from "./StatCell.svelte";
  import StatSection from "./StatSection.svelte";

  export let friendTeam: data.Team;
  export let enemyTeam: data.Team;

  let friendTeamAvg: data.TeamAverageStats = {
    ship_avg_pr: 0,
    ship_avg_damage: 0,
    ship_win_rate: 0,
    overall_avg_pr: 0,
    overall_avg_damage: 0,
    overall_win_rate: 0,
  };
  let enemyTeamAvg: data.TeamAverageStats = {
    ship_avg_pr: 0,
    ship_avg_damage: 0,
    ship_win_rate: 0,
    overall_avg_pr: 0,
    overall_avg_damage: 0,
    overall_win_rate: 0,
  };

  $: if (friendTeam) {
    friendTeamAvg = calculateTeamAverage(friendTeam);
  }
  $: if (enemyTeam) {
    enemyTeamAvg = calculateTeamAverage(enemyTeam);
  }

  function calculateTeamAverage(team: data.Team): data.TeamAverageStats {
    const shipTypes = ["cv", "bb", "cl", "dd", "ss"] as const;
    let totalShipPR = 0;
    let totalShipDamage = 0;
    let totalShipWinRate = 0;
    let totalOverallPR = 0;
    let totalOverallDamage = 0;
    let totalOverallWinRate = 0;

    shipTypes.forEach((type) => {
      const stats = team.ship_type_stats?.[type];
      if (stats) {
        totalShipPR += stats.ship_avg_pr;
        totalShipDamage += stats.ship_avg_damage;
        totalShipWinRate += stats.ship_win_rate;
        totalOverallPR += stats.overall_avg_pr;
        totalOverallDamage += stats.overall_avg_damage;
        totalOverallWinRate += stats.overall_win_rate;
      }
    });

    const count = shipTypes.length;
    return {
      ship_avg_pr: totalShipPR / count,
      ship_avg_damage: totalShipDamage / count,
      ship_win_rate: totalShipWinRate / count,
      overall_avg_pr: totalOverallPR / count,
      overall_avg_damage: totalOverallDamage / count,
      overall_win_rate: totalOverallWinRate / count,
    };
  }

  function formatValue(value: number): string {
    if (!value || value === 0) return "0";
    return value.toFixed(0);
  }

  function formatPercent(value: number): string {
    return value.toFixed(1);
  }
</script>

<div class="flex gap-8 mt-4 px-4 justify-center flex-wrap">
  <!-- Team Average Column -->
  <div class="bg-base-200 rounded-xl border border-base-300 p-4 min-w-fit">
    <h3 class="font-semibold text-lg mb-4 text-center">チーム平均</h3>
    <div class="space-y-3">
      <div class="pb-2">
        <StatSection title="艦成績">
          <StatCell
            label="PR"
            friendValue={formatValue(friendTeamAvg.ship_avg_pr)}
            enemyValue={formatValue(enemyTeamAvg.ship_avg_pr)}
            diffValue={formatValue(
              friendTeamAvg.ship_avg_pr - enemyTeamAvg.ship_avg_pr,
            )}
            diffValueNum={friendTeamAvg.ship_avg_pr - enemyTeamAvg.ship_avg_pr}
          />
          <StatCell
            label="Dmg"
            friendValue={formatValue(friendTeamAvg.ship_avg_damage)}
            enemyValue={formatValue(enemyTeamAvg.ship_avg_damage)}
            diffValue={formatValue(
              friendTeamAvg.ship_avg_damage - enemyTeamAvg.ship_avg_damage,
            )}
            diffValueNum={friendTeamAvg.ship_avg_damage -
              enemyTeamAvg.ship_avg_damage}
          />
          <StatCell
            label="勝率"
            friendValue={formatPercent(friendTeamAvg.ship_win_rate)}
            enemyValue={formatPercent(enemyTeamAvg.ship_win_rate)}
            diffValue={formatPercent(
              friendTeamAvg.ship_win_rate - enemyTeamAvg.ship_win_rate,
            )}
            diffValueNum={friendTeamAvg.ship_win_rate -
              enemyTeamAvg.ship_win_rate}
            isPercentage={true}
          />
        </StatSection>
      </div>

      <div class="border-t border-base-300 pt-2">
        <StatSection title="総合成績">
          <StatCell
            label="PR"
            friendValue={formatValue(friendTeamAvg.overall_avg_pr)}
            enemyValue={formatValue(enemyTeamAvg.overall_avg_pr)}
            diffValue={formatValue(
              friendTeamAvg.overall_avg_pr - enemyTeamAvg.overall_avg_pr,
            )}
            diffValueNum={friendTeamAvg.overall_avg_pr -
              enemyTeamAvg.overall_avg_pr}
          />
          <StatCell
            label="Dmg"
            friendValue={formatValue(friendTeamAvg.overall_avg_damage)}
            enemyValue={formatValue(enemyTeamAvg.overall_avg_damage)}
            diffValue={formatValue(
              friendTeamAvg.overall_avg_damage -
                enemyTeamAvg.overall_avg_damage,
            )}
            diffValueNum={friendTeamAvg.overall_avg_damage -
              enemyTeamAvg.overall_avg_damage}
          />
          <StatCell
            label="勝率"
            friendValue={formatPercent(friendTeamAvg.overall_win_rate)}
            enemyValue={formatPercent(enemyTeamAvg.overall_win_rate)}
            diffValue={formatPercent(
              friendTeamAvg.overall_win_rate - enemyTeamAvg.overall_win_rate,
            )}
            diffValueNum={friendTeamAvg.overall_win_rate -
              enemyTeamAvg.overall_win_rate}
            isPercentage={true}
          />
        </StatSection>
      </div>

      <div class="border-t border-base-300 pt-2">
        <StatSection title="戦力評価(by 178usagi)">
          <StatCell
            label="脅威度"
            friendValue={friendTeam.pvp_all.team_threat_level.average.toFixed(
              1,
            )}
            enemyValue={enemyTeam.pvp_all.team_threat_level.average.toFixed(1)}
            diffValue={(
              friendTeam.pvp_all.team_threat_level.average -
              enemyTeam.pvp_all.team_threat_level.average
            ).toFixed(1)}
            diffValueNum={friendTeam.pvp_all.team_threat_level.average -
              enemyTeam.pvp_all.team_threat_level.average}
          />
          <StatCell
            label="確度"
            friendValue={friendTeam.pvp_all.team_threat_level.accuracy.toFixed(
              1,
            )}
            enemyValue={enemyTeam.pvp_all.team_threat_level.accuracy.toFixed(1)}
            diffValue={(
              friendTeam.pvp_all.team_threat_level.accuracy -
              enemyTeam.pvp_all.team_threat_level.accuracy
            ).toFixed(1)}
            diffValueNum={friendTeam.pvp_all.team_threat_level.accuracy -
              enemyTeam.pvp_all.team_threat_level.accuracy}
          />
          <StatCell
            label="介護指数"
            friendValue={friendTeam.pvp_all.team_threat_level.dissociation_degree.toFixed(
              1,
            )}
            enemyValue={enemyTeam.pvp_all.team_threat_level.dissociation_degree.toFixed(
              1,
            )}
            diffValue={(
              friendTeam.pvp_all.team_threat_level.dissociation_degree -
              enemyTeam.pvp_all.team_threat_level.dissociation_degree
            ).toFixed(1)}
            diffValueNum={friendTeam.pvp_all.team_threat_level
              .dissociation_degree -
              enemyTeam.pvp_all.team_threat_level.dissociation_degree}
          />
        </StatSection>
      </div>
    </div>
  </div>

  <!-- Ship Type Columns -->
  {#each Array.from(SHIP_TYPES.keys()) as shipType (shipType)}
    {@const shipTypeName = SHIP_TYPES.get(shipType) ?? shipType}
    {@const friendStats = friendTeam.ship_type_stats?.[shipType]}
    {@const enemyStats = enemyTeam.ship_type_stats?.[shipType]}

    {#if friendStats && enemyStats}
      <div class="bg-base-200 rounded-xl border border-base-300 p-4 min-w-fit">
        <h3 class="font-semibold text-lg mb-4 text-center">{shipTypeName}</h3>
        <div class="space-y-3">
          <div class="pb-2">
            <StatSection title="艦成績">
              <StatCell
                label="PR"
                friendValue={formatValue(friendStats.ship_avg_pr)}
                enemyValue={formatValue(enemyStats.ship_avg_pr)}
                diffValue={formatValue(
                  friendStats.ship_avg_pr - enemyStats.ship_avg_pr,
                )}
                diffValueNum={friendStats.ship_avg_pr - enemyStats.ship_avg_pr}
              />
              <StatCell
                label="Dmg"
                friendValue={formatValue(friendStats.ship_avg_damage)}
                enemyValue={formatValue(enemyStats.ship_avg_damage)}
                diffValue={formatValue(
                  friendStats.ship_avg_damage - enemyStats.ship_avg_damage,
                )}
                diffValueNum={friendStats.ship_avg_damage -
                  enemyStats.ship_avg_damage}
              />
              <StatCell
                label="勝率"
                friendValue={formatPercent(friendStats.ship_win_rate)}
                enemyValue={formatPercent(enemyStats.ship_win_rate)}
                diffValue={formatPercent(
                  friendStats.ship_win_rate - enemyStats.ship_win_rate,
                )}
                diffValueNum={friendStats.ship_win_rate -
                  enemyStats.ship_win_rate}
                isPercentage={true}
              />
            </StatSection>
          </div>

          <div class="border-t border-base-300 pt-2">
            <StatSection title="総合成績">
              <StatCell
                label="PR"
                friendValue={formatValue(friendStats.overall_avg_pr)}
                enemyValue={formatValue(enemyStats.overall_avg_pr)}
                diffValue={formatValue(
                  friendStats.overall_avg_pr - enemyStats.overall_avg_pr,
                )}
                diffValueNum={friendStats.overall_avg_pr -
                  enemyStats.overall_avg_pr}
              />
              <StatCell
                label="Dmg"
                friendValue={formatValue(friendStats.overall_avg_damage)}
                enemyValue={formatValue(enemyStats.overall_avg_damage)}
                diffValue={formatValue(
                  friendStats.overall_avg_damage -
                    enemyStats.overall_avg_damage,
                )}
                diffValueNum={friendStats.overall_avg_damage -
                  enemyStats.overall_avg_damage}
              />
              <StatCell
                label="勝率"
                friendValue={formatPercent(friendStats.overall_win_rate)}
                enemyValue={formatPercent(enemyStats.overall_win_rate)}
                diffValue={formatPercent(
                  friendStats.overall_win_rate - enemyStats.overall_win_rate,
                )}
                diffValueNum={friendStats.overall_win_rate -
                  enemyStats.overall_win_rate}
                isPercentage={true}
              />
            </StatSection>
          </div>
        </div>
      </div>
    {/if}
  {/each}
</div>
