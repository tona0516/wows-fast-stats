<script lang="ts">
  import type { data } from "@wails/go/models";
  import BarChart from "./BarChart.svelte";
  import { storedUserConfig } from "@libs/stores";
  import { get } from "svelte/store";
  import type { StatsExtra } from "@libs/types";

  export let friendTeam: data.Team;
  export let enemyTeam: data.Team;

  type ChartData = {
    label: string;
    friendValue: number;
    enemyValue: number;
  };

  let friendTeamStats: data.TeamStats;
  let enemyTeamStats: data.TeamStats;

  let shipChartData: ChartData[] = [];
  let overallChartData: ChartData[] = [];
  let threatChartData: ChartData[] = [];

  $: if (friendTeam) {
    friendTeamStats = getTeamStats(friendTeam);
  }
  $: if (enemyTeam) {
    enemyTeamStats = getTeamStats(enemyTeam);
  }

  $: if (friendTeam && enemyTeam) {
    shipChartData = [
      {
        label: "PR",
        friendValue: friendTeamStats.team_average_stats.ship_pr,
        enemyValue: enemyTeamStats.team_average_stats.ship_pr,
      },
      {
        label: "Dmg",
        friendValue: friendTeamStats.team_average_stats.ship_damage,
        enemyValue: enemyTeamStats.team_average_stats.ship_damage,
      },
      {
        label: "勝率",
        friendValue: friendTeamStats.team_average_stats.ship_win_rate,
        enemyValue: enemyTeamStats.team_average_stats.ship_win_rate,
      },
      {
        label: "戦闘数",
        friendValue: friendTeamStats.team_average_stats.ship_battles,
        enemyValue: enemyTeamStats.team_average_stats.ship_battles,
      },
    ];

    overallChartData = [
      {
        label: "PR",
        friendValue: friendTeamStats.team_average_stats.overall_pr,
        enemyValue: enemyTeamStats.team_average_stats.overall_pr,
      },
      {
        label: "Dmg",
        friendValue: friendTeamStats.team_average_stats.overall_damage,
        enemyValue: enemyTeamStats.team_average_stats.overall_damage,
      },
      {
        label: "勝率",
        friendValue: friendTeamStats.team_average_stats.overall_win_rate,
        enemyValue: enemyTeamStats.team_average_stats.overall_win_rate,
      },
      {
        label: "戦闘数",
        friendValue: friendTeamStats.team_average_stats.overall_battles,
        enemyValue: enemyTeamStats.team_average_stats.overall_battles,
      },
    ];

    threatChartData = [
      {
        label: "脅威度",
        friendValue: friendTeam.pvp_all.team_threat_level.average,
        enemyValue: enemyTeam.pvp_all.team_threat_level.average,
      },
      {
        label: "確度",
        friendValue: friendTeam.pvp_all.team_threat_level.accuracy,
        enemyValue: enemyTeam.pvp_all.team_threat_level.accuracy,
      },
      {
        label: "介護指数",
        friendValue: friendTeam.pvp_all.team_threat_level.dissociation_degree,
        enemyValue: enemyTeam.pvp_all.team_threat_level.dissociation_degree,
      },
    ];
  }

  function getTeamStats(team: data.Team): data.TeamStats {
    const statsExtra = get(storedUserConfig).stats_extra as StatsExtra;
    return team[statsExtra];
  }
</script>

<div class="mt-4 px-4 space-y-6">
  <!-- Charts Section -->
  <div class="bg-base-200 rounded-xl border border-base-300 p-6">
    <h3 class="font-semibold text-xl mb-6 text-center">チーム平均比較</h3>
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <div>
        <h4 class="font-medium text-center mb-4">艦成績</h4>
        <BarChart data={shipChartData} />
      </div>
      <div>
        <h4 class="font-medium text-center mb-4">総合成績</h4>
        <BarChart data={overallChartData} />
      </div>
      <div>
        <h4 class="font-medium text-center mb-4">戦力評価(by 178usagi)</h4>
        <BarChart data={threatChartData} />
      </div>
    </div>
  </div>
</div>
