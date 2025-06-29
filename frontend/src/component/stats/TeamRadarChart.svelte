<script lang="ts">
  import { formatWithSuffix, toPlayerStats } from "src/lib/util";
  import { storedConfig } from "src/stores";
  import type { data } from "wailsjs/go/models";

  export let battle: data.Battle;

  function getPlayerStats(team: data.Team): data.PlayerStats[] {
    return team.players
      .filter((p) => p.player_info.id !== 0 && !p.player_info.is_hidden)
      .map((p) => toPlayerStats(p, $storedConfig.stats_pattern));
  }

  function average(values: number[]): number {
    if (values.length === 0) {
      return 0;
    }
    return values.reduce((a, b) => a + b) / values.length;
  }

  function averagePR(team: data.Team): number {
    const values = getPlayerStats(team).map((p) => p.overall.pr);
    return average(values);
  }

  function averageThreatLevel(team: data.Team): number {
    const values = getPlayerStats(team).map(
      (p) => p.overall.threat_level.modified,
    );
    return average(values);
  }

  function averageDamage(team: data.Team): number {
    const values = getPlayerStats(team).map((p) => p.overall.damage);
    return average(values);
  }

  function averageWinRate(team: data.Team): number {
    const values = getPlayerStats(team).map((p) => p.overall.win_rate);
    return average(values);
  }

  function averageBattle(team: data.Team): number {
    const values = getPlayerStats(team).map((p) => p.overall.battles);
    return average(values);
  }

  const chartData = [
    {
      label: "PR",
      friend: averagePR(battle.teams[0]),
      enemy: averagePR(battle.teams[1]),
    },
    {
      label: "戦力評価",
      friend: averageThreatLevel(battle.teams[0]),
      enemy: averageThreatLevel(battle.teams[1]),
    },
    {
      label: "ダメージ",
      friend: averageDamage(battle.teams[0]),
      enemy: averageDamage(battle.teams[1]),
    },
    {
      label: "勝率",
      friend: averageWinRate(battle.teams[0]),
      enemy: averageWinRate(battle.teams[1]),
    },
    {
      label: "戦闘数",
      friend: averageBattle(battle.teams[0]),
      enemy: averageBattle(battle.teams[1]),
    },
  ];
</script>

<div class="w-2xl">
  <table
    class="charts-css column multiple show-labels data-spacing-10 datasets-spacing-4"
  >
    <tbody>
      {#each chartData as item}
        {@const max = Math.max(item.friend, item.enemy)}
        <tr>
          <th scope="row">{item.label}</th>
          <td style="--size: calc({item.friend}/{max}; --color: #00A95F;"
            ><span class="text-nowrap text-neutral-100"
              >{formatWithSuffix(item.friend)}</span
            ></td
          >
          <td style="--size: calc({item.enemy}/{max}); --color: #EA5532;"
            ><span class="text-nowrap text-neutral-100"
              >{formatWithSuffix(item.enemy)}</span
            ></td
          >
        </tr>
      {/each}
    </tbody>
  </table>
</div>
