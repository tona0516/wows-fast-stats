<script lang="ts">
  import type { core } from "@wails/go/models";
  import BarChart from "./BarChart.svelte";
  import { get } from "svelte/store";
  import { storedDisplayPref } from "@libs/stores";

  export let friendTeam: core.Team;
  export let enemyTeam: core.Team;

  type ChartData = {
    label: string;
    friendValue: number;
    enemyValue: number;
  };

  let friendTeamStats: core.TeamStats;
  let enemyTeamStats: core.TeamStats;

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
        friendValue: friendTeamStats.teamAverageStats.ship_pr,
        enemyValue: enemyTeamStats.teamAverageStats.ship_pr,
      },
      {
        label: "Dmg",
        friendValue: friendTeamStats.teamAverageStats.ship_damage,
        enemyValue: enemyTeamStats.teamAverageStats.ship_damage,
      },
      {
        label: "勝率",
        friendValue: friendTeamStats.teamAverageStats.ship_win_rate,
        enemyValue: enemyTeamStats.teamAverageStats.ship_win_rate,
      },
      {
        label: "戦闘数",
        friendValue: friendTeamStats.teamAverageStats.ship_battles,
        enemyValue: enemyTeamStats.teamAverageStats.ship_battles,
      },
    ];

    overallChartData = [
      {
        label: "PR",
        friendValue: friendTeamStats.teamAverageStats.overall_pr,
        enemyValue: enemyTeamStats.teamAverageStats.overall_pr,
      },
      {
        label: "Dmg",
        friendValue: friendTeamStats.teamAverageStats.overall_damage,
        enemyValue: enemyTeamStats.teamAverageStats.overall_damage,
      },
      {
        label: "勝率",
        friendValue: friendTeamStats.teamAverageStats.overall_win_rate,
        enemyValue: enemyTeamStats.teamAverageStats.overall_win_rate,
      },
      {
        label: "戦闘数",
        friendValue: friendTeamStats.teamAverageStats.overall_battles,
        enemyValue: enemyTeamStats.teamAverageStats.overall_battles,
      },
    ];

    threatChartData = [
      {
        label: "脅威度",
        friendValue: friendTeam.pvpAll.teamThreatLevel.average,
        enemyValue: enemyTeam.pvpAll.teamThreatLevel.average,
      },
      {
        label: "確度",
        friendValue: friendTeam.pvpAll.teamThreatLevel.accuracy,
        enemyValue: enemyTeam.pvpAll.teamThreatLevel.accuracy,
      },
      {
        label: "介護指数",
        friendValue: friendTeam.pvpAll.teamThreatLevel.dissociationDegree,
        enemyValue: enemyTeam.pvpAll.teamThreatLevel.dissociationDegree,
      },
    ];
  }

  function getTeamStats(team: core.Team): core.TeamStats {
    const statsExtra = get(storedDisplayPref).statsExtra;
    return team[statsExtra];
  }
</script>

<div class="mt-4 px-4 space-y-6">
  <div class="bg-base-200 border border-base-300 p-4 space-y-4">
    <h3 class="font-bold text-lg text-center">チーム平均</h3>
    <div class="grid grid-cols-3 gap-5">
      <div class="bg-base-100 p-4">
        <h4 class="font-semibold text-center mb-4 text-sm">艦成績</h4>
        <BarChart data={shipChartData} />
      </div>
      <div class="bg-base-100 p-4">
        <h4 class="font-semibold text-center mb-4 text-sm">総合成績</h4>
        <BarChart data={overallChartData} />
      </div>
      <div class="bg-base-100 p-4">
        <h4 class="font-semibold text-center mb-4 text-sm">戦力評価(by 178usagi)</h4>
        <BarChart data={threatChartData} />
      </div>
    </div>
  </div>
</div>
