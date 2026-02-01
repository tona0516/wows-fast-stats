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

  let allChartData: {
    title: string;
    data: ChartData[];
  }[] = [];

  $: if (friendTeam) {
    friendTeamStats = getTeamStats(friendTeam);
  }
  $: if (enemyTeam) {
    enemyTeamStats = getTeamStats(enemyTeam);
  }

  $: if (friendTeamStats && enemyTeamStats) {
    allChartData = [
      {
        title: "PR",
        data: [
          {
            label: "艦成績",
            friendValue: friendTeamStats.teamAverageStats.ship_pr,
            enemyValue: enemyTeamStats.teamAverageStats.ship_pr,
          },
          {
            label: "総合成績",
            friendValue: friendTeamStats.teamAverageStats.overall_pr,
            enemyValue: enemyTeamStats.teamAverageStats.overall_pr,
          },
        ],
      },
      {
        title: "ダメージ",
        data: [
          {
            label: "艦成績",
            friendValue: friendTeamStats.teamAverageStats.ship_damage,
            enemyValue: enemyTeamStats.teamAverageStats.ship_damage,
          },
          {
            label: "総合成績",
            friendValue: friendTeamStats.teamAverageStats.overall_damage,
            enemyValue: enemyTeamStats.teamAverageStats.overall_damage,
          },
        ],
      },
      {
        title: "勝率",
        data: [
          {
            label: "艦成績",
            friendValue: friendTeamStats.teamAverageStats.ship_win_rate,
            enemyValue: enemyTeamStats.teamAverageStats.ship_win_rate,
          },
          {
            label: "総合成績",
            friendValue: friendTeamStats.teamAverageStats.overall_win_rate,
            enemyValue: enemyTeamStats.teamAverageStats.overall_win_rate,
          },
        ],
      },
      {
        title: "戦闘数",
        data: [
          {
            label: "艦成績",
            friendValue: friendTeamStats.teamAverageStats.ship_battles,
            enemyValue: enemyTeamStats.teamAverageStats.ship_battles,
          },
          {
            label: "総合成績",
            friendValue: friendTeamStats.teamAverageStats.overall_battles,
            enemyValue: enemyTeamStats.teamAverageStats.overall_battles,
          },
        ],
      },
      {
        title: "戦力評価 (by 178usagi)",
        data: [
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
        ],
      },
    ];
  }

  function getTeamStats(team: core.Team): core.TeamStats {
    const statsExtra = get(storedDisplayPref).statsExtra;
    return team[statsExtra];
  }
</script>

<div class="mt-4">
  <section
    class="bg-base-200/70 border border-base-300/80 p-2 shadow-sm rounded-md"
  >
    <div class="flex items-center justify-between mb-2">
      <span class="w-full text-center font-bold">チーム平均値の比較</span>
      <div class="hidden md:flex items-center gap-2 whitespace-nowrap">
        <span class="inline-flex h-2 w-2 rounded-full bg-ally"></span>
        味方
        <span class="inline-flex h-2 w-2 rounded-full bg-enemy"></span>
        敵
      </div>
    </div>

    <div class="grid gap-2 grid-cols-5">
      {#each allChartData as chartData}
        <div
          class="relative overflow-hidden border border-base-300 bg-base-100/70 p-4 shadow-sm rounded-md"
        >
          <div class="flex items-center justify-between mb-2">
            <h4 class="w-full text-center font-semibold line-clamp-1">
              {chartData.title}
            </h4>
          </div>
          <BarChart data={chartData.data} />
        </div>
      {/each}
    </div>
  </section>
</div>
