<script lang="ts">
  import { SHIP_TYPES } from "@libs/constants";
  import type { data } from "@wails/go/models";

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

  function formatValue(value: number): string {
    if (!value || value === 0) return "0";
    return value.toFixed(0);
  }

  function getDiffColor(diff: number): string {
    if (diff > 0) return "text-green-600";
    if (diff < 0) return "text-red-600";
    return "text-gray-600";
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
</script>

<div class="flex gap-8 mt-4 px-4 justify-center flex-wrap">
  <!-- Team Average Column -->
  <div class="bg-base-200 rounded-xl border border-base-300 p-4 min-w-fit">
    <h3 class="font-semibold text-lg mb-4 text-center">チーム平均</h3>
    <div class="space-y-3">
      <div class="flex items-center gap-3 text-sm">
        <span class="w-12">艦PR</span>
        <div class="flex gap-2">
          <div class="text-center w-12">
            <div class="text-xs text-gray-500">味方</div>
            <div>{formatValue(friendTeamAvg.ship_avg_pr)}</div>
          </div>
          <div class="text-center w-12">
            <div class="text-xs text-gray-500">敵</div>
            <div>{formatValue(enemyTeamAvg.ship_avg_pr)}</div>
          </div>
          <div class="text-center w-12">
            <div class="text-xs text-gray-500">差分</div>
            <div
              class={getDiffColor(
                friendTeamAvg.ship_avg_pr - enemyTeamAvg.ship_avg_pr,
              )}
            >
              {formatValue(
                friendTeamAvg.ship_avg_pr - enemyTeamAvg.ship_avg_pr,
              )}
            </div>
          </div>
        </div>
      </div>
      <div class="flex items-center gap-3 text-sm">
        <span class="w-12">艦ダメ</span>
        <div class="flex gap-2">
          <div class="text-center w-12">
            <div class="text-xs text-gray-500">味方</div>
            <div>{formatValue(friendTeamAvg.ship_avg_damage)}</div>
          </div>
          <div class="text-center w-12">
            <div class="text-xs text-gray-500">敵</div>
            <div>{formatValue(enemyTeamAvg.ship_avg_damage)}</div>
          </div>
          <div class="text-center w-12">
            <div class="text-xs text-gray-500">差分</div>
            <div
              class={getDiffColor(
                friendTeamAvg.ship_avg_damage - enemyTeamAvg.ship_avg_damage,
              )}
            >
              {formatValue(
                friendTeamAvg.ship_avg_damage - enemyTeamAvg.ship_avg_damage,
              )}
            </div>
          </div>
        </div>
      </div>
      <div class="flex items-center gap-3 text-sm">
        <span class="w-12">艦勝率</span>
        <div class="flex gap-2">
          <div class="text-center w-12">
            <div class="text-xs text-gray-500">味方</div>
            <div>{friendTeamAvg.ship_win_rate.toFixed(1)}%</div>
          </div>
          <div class="text-center w-12">
            <div class="text-xs text-gray-500">敵</div>
            <div>{enemyTeamAvg.ship_win_rate.toFixed(1)}%</div>
          </div>
          <div class="text-center w-12">
            <div class="text-xs text-gray-500">差分</div>
            <div
              class={getDiffColor(
                friendTeamAvg.ship_win_rate - enemyTeamAvg.ship_win_rate,
              )}
            >
              {(
                friendTeamAvg.ship_win_rate - enemyTeamAvg.ship_win_rate
              ).toFixed(1)}%
            </div>
          </div>
        </div>
      </div>
      <div class="flex items-center gap-3 text-sm">
        <span class="w-12">総PR</span>
        <div class="flex gap-2">
          <div class="text-center w-12">
            <div class="text-xs text-gray-500">味方</div>
            <div>{formatValue(friendTeamAvg.overall_avg_pr)}</div>
          </div>
          <div class="text-center w-12">
            <div class="text-xs text-gray-500">敵</div>
            <div>{formatValue(enemyTeamAvg.overall_avg_pr)}</div>
          </div>
          <div class="text-center w-12">
            <div class="text-xs text-gray-500">差分</div>
            <div
              class={getDiffColor(
                friendTeamAvg.overall_avg_pr - enemyTeamAvg.overall_avg_pr,
              )}
            >
              {formatValue(
                friendTeamAvg.overall_avg_pr - enemyTeamAvg.overall_avg_pr,
              )}
            </div>
          </div>
        </div>
      </div>
      <div class="flex items-center gap-3 text-sm">
        <span class="w-12">総ダメ</span>
        <div class="flex gap-2">
          <div class="text-center w-12">
            <div class="text-xs text-gray-500">味方</div>
            <div>{formatValue(friendTeamAvg.overall_avg_damage)}</div>
          </div>
          <div class="text-center w-12">
            <div class="text-xs text-gray-500">敵</div>
            <div>{formatValue(enemyTeamAvg.overall_avg_damage)}</div>
          </div>
          <div class="text-center w-12">
            <div class="text-xs text-gray-500">差分</div>
            <div
              class={getDiffColor(
                friendTeamAvg.overall_avg_damage -
                  enemyTeamAvg.overall_avg_damage,
              )}
            >
              {formatValue(
                friendTeamAvg.overall_avg_damage -
                  enemyTeamAvg.overall_avg_damage,
              )}
            </div>
          </div>
        </div>
      </div>
      <div class="flex items-center gap-3 text-sm">
        <span class="w-12">総勝率</span>
        <div class="flex gap-2">
          <div class="text-center w-12">
            <div class="text-xs text-gray-500">味方</div>
            <div>{friendTeamAvg.overall_win_rate.toFixed(1)}%</div>
          </div>
          <div class="text-center w-12">
            <div class="text-xs text-gray-500">敵</div>
            <div>{enemyTeamAvg.overall_win_rate.toFixed(1)}%</div>
          </div>
          <div class="text-center w-12">
            <div class="text-xs text-gray-500">差分</div>
            <div
              class={getDiffColor(
                friendTeamAvg.overall_win_rate - enemyTeamAvg.overall_win_rate,
              )}
            >
              {(
                friendTeamAvg.overall_win_rate - enemyTeamAvg.overall_win_rate
              ).toFixed(1)}%
            </div>
          </div>
        </div>
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
          <div class="flex items-center gap-3 text-sm">
            <span class="w-12">艦PR</span>
            <div class="flex gap-2">
              <div class="text-center w-12">
                <div class="text-xs text-gray-500">味方</div>
                <div>{formatValue(friendStats.ship_avg_pr)}</div>
              </div>
              <div class="text-center w-12">
                <div class="text-xs text-gray-500">敵</div>
                <div>{formatValue(enemyStats.ship_avg_pr)}</div>
              </div>
              <div class="text-center w-12">
                <div class="text-xs text-gray-500">差分</div>
                <div
                  class={getDiffColor(
                    friendStats.ship_avg_pr - enemyStats.ship_avg_pr,
                  )}
                >
                  {formatValue(
                    friendStats.ship_avg_pr - enemyStats.ship_avg_pr,
                  )}
                </div>
              </div>
            </div>
          </div>
          <div class="flex items-center gap-3 text-sm">
            <span class="w-12">艦ダメ</span>
            <div class="flex gap-2">
              <div class="text-center w-12">
                <div class="text-xs text-gray-500">味方</div>
                <div>{formatValue(friendStats.ship_avg_damage)}</div>
              </div>
              <div class="text-center w-12">
                <div class="text-xs text-gray-500">敵</div>
                <div>{formatValue(enemyStats.ship_avg_damage)}</div>
              </div>
              <div class="text-center w-12">
                <div class="text-xs text-gray-500">差分</div>
                <div
                  class={getDiffColor(
                    friendStats.ship_avg_damage - enemyStats.ship_avg_damage,
                  )}
                >
                  {formatValue(
                    friendStats.ship_avg_damage - enemyStats.ship_avg_damage,
                  )}
                </div>
              </div>
            </div>
          </div>
          <div class="flex items-center gap-3 text-sm">
            <span class="w-12">艦勝率</span>
            <div class="flex gap-2">
              <div class="text-center w-12">
                <div class="text-xs text-gray-500">味方</div>
                <div>{friendStats.ship_win_rate.toFixed(1)}%</div>
              </div>
              <div class="text-center w-12">
                <div class="text-xs text-gray-500">敵</div>
                <div>{enemyStats.ship_win_rate.toFixed(1)}%</div>
              </div>
              <div class="text-center w-12">
                <div class="text-xs text-gray-500">差分</div>
                <div
                  class={getDiffColor(
                    friendStats.ship_win_rate - enemyStats.ship_win_rate,
                  )}
                >
                  {(
                    friendStats.ship_win_rate - enemyStats.ship_win_rate
                  ).toFixed(1)}%
                </div>
              </div>
            </div>
          </div>
          <div class="flex items-center gap-3 text-sm">
            <span class="w-12">総PR</span>
            <div class="flex gap-2">
              <div class="text-center w-12">
                <div class="text-xs text-gray-500">味方</div>
                <div>{formatValue(friendStats.overall_avg_pr)}</div>
              </div>
              <div class="text-center w-12">
                <div class="text-xs text-gray-500">敵</div>
                <div>{formatValue(enemyStats.overall_avg_pr)}</div>
              </div>
              <div class="text-center w-12">
                <div class="text-xs text-gray-500">差分</div>
                <div
                  class={getDiffColor(
                    friendStats.overall_avg_pr - enemyStats.overall_avg_pr,
                  )}
                >
                  {formatValue(
                    friendStats.overall_avg_pr - enemyStats.overall_avg_pr,
                  )}
                </div>
              </div>
            </div>
          </div>
          <div class="flex items-center gap-3 text-sm">
            <span class="w-12">総ダメ</span>
            <div class="flex gap-2">
              <div class="text-center w-12">
                <div class="text-xs text-gray-500">味方</div>
                <div>{formatValue(friendStats.overall_avg_damage)}</div>
              </div>
              <div class="text-center w-12">
                <div class="text-xs text-gray-500">敵</div>
                <div>{formatValue(enemyStats.overall_avg_damage)}</div>
              </div>
              <div class="text-center w-12">
                <div class="text-xs text-gray-500">差分</div>
                <div
                  class={getDiffColor(
                    friendStats.overall_avg_damage -
                      enemyStats.overall_avg_damage,
                  )}
                >
                  {formatValue(
                    friendStats.overall_avg_damage -
                      enemyStats.overall_avg_damage,
                  )}
                </div>
              </div>
            </div>
          </div>
          <div class="flex items-center gap-3 text-sm">
            <span class="w-12">総勝率</span>
            <div class="flex gap-2">
              <div class="text-center w-12">
                <div class="text-xs text-gray-500">味方</div>
                <div>{friendStats.overall_win_rate.toFixed(1)}%</div>
              </div>
              <div class="text-center w-12">
                <div class="text-xs text-gray-500">敵</div>
                <div>{enemyStats.overall_win_rate.toFixed(1)}%</div>
              </div>
              <div class="text-center w-12">
                <div class="text-xs text-gray-500">差分</div>
                <div
                  class={getDiffColor(
                    friendStats.overall_win_rate - enemyStats.overall_win_rate,
                  )}
                >
                  {(
                    friendStats.overall_win_rate - enemyStats.overall_win_rate
                  ).toFixed(1)}%
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    {/if}
  {/each}
</div>

<style>
  :global(.text-green-600) {
    color: #16a34a;
  }

  :global(.text-red-600) {
    color: #dc2626;
  }

  :global(.text-gray-600) {
    color: #4b5563;
  }
</style>
