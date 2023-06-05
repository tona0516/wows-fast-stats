<script lang="ts">
import type { vo } from "wailsjs/go/models";
import Const from "./Const";
import ShipPr from "./ShipPr.svelte";
import BasicPlayerName from "./BasicPlayerName.svelte";
import BasicShipInfo from "./BasicShipInfo.svelte";
import ShipDamage from "./ShipDamage.svelte";
import ShipWinRate from "./ShipWinRate.svelte";
import ShipKdRate from "./ShipKdRate.svelte";
import ShipSurvivedRate from "./ShipSurvivedRate.svelte";
import ShipExp from "./ShipExp.svelte";
import ShipBattles from "./ShipBattles.svelte";
import OverallDamage from "./OverallDamage.svelte";
import OverallWinRate from "./OverallWinRate.svelte";
import OverallKdRate from "./OverallKdRate.svelte";
import OverallSurvivedRate from "./OverallSurvivedRate.svelte";
import OverallExp from "./OverallExp.svelte";
import OverallAvgTier from "./OverallAvgTier.svelte";
import OverallShipTypeRate from "./OverallShipTypeRate.svelte";
import OverallTierRate from "./OverallTierRate.svelte";
import OverallBattles from "./OverallBattles.svelte";
import NoData from "./NoData.svelte";
import BasicIsInAvg from "./BasicIsInAvg.svelte";
import type { DisplayPattern } from "./DisplayPattern";
import ShipHitRate from "./ShipHitRate.svelte";
import { storedBattle, storedSummaryResult, storedUserConfig } from "./stores";
import { get } from "svelte/store";
import type { StatsCategory } from "./StatsCategory";
import { LogDebug } from "../wailsjs/runtime/runtime";

let battle = get(storedBattle);
storedBattle.subscribe((it) => (battle = it));

let userConfig = get(storedUserConfig);
storedUserConfig.subscribe((it) => (userConfig = it));

let summaryResult = get(storedSummaryResult);
storedSummaryResult.subscribe((it) => {
  summaryResult = it;
});

type ComponentInfo = {
  category: StatsCategory;
  name: string;
  component: any;
  column: number;
};

const components: ComponentInfo[] = [
  // basic
  {
    category: "basic",
    name: "player_name",
    component: BasicPlayerName,
    column: 1,
  },
  { category: "basic", name: "ship_info", component: BasicShipInfo, column: 2 },
  // ship
  { category: "ship", name: "pr", component: ShipPr, column: 1 },
  { category: "ship", name: "damage", component: ShipDamage, column: 1 },
  { category: "ship", name: "win_rate", component: ShipWinRate, column: 1 },
  { category: "ship", name: "kd_rate", component: ShipKdRate, column: 1 },
  { category: "ship", name: "exp", component: ShipExp, column: 1 },
  { category: "ship", name: "battles", component: ShipBattles, column: 1 },
  {
    category: "ship",
    name: "survived_rate",
    component: ShipSurvivedRate,
    column: 1,
  },
  { category: "ship", name: "hit_rate", component: ShipHitRate, column: 1 },
  // overall
  { category: "overall", name: "damage", component: OverallDamage, column: 1 },
  {
    category: "overall",
    name: "win_rate",
    component: OverallWinRate,
    column: 1,
  },
  { category: "overall", name: "kd_rate", component: OverallKdRate, column: 1 },
  { category: "overall", name: "exp", component: OverallExp, column: 1 },
  {
    category: "overall",
    name: "battles",
    component: OverallBattles,
    column: 1,
  },
  {
    category: "overall",
    name: "survived_rate",
    component: OverallSurvivedRate,
    column: 1,
  },
  {
    category: "overall",
    name: "avg_tier",
    component: OverallAvgTier,
    column: 1,
  },
  {
    category: "overall",
    name: "using_ship_type_rate",
    component: OverallShipTypeRate,
    column: 1,
  },
  {
    category: "overall",
    name: "using_tier_rate",
    component: OverallTierRate,
    column: 1,
  },
];

function decidePlayerDataPattern(
  player: vo.Player,
  statsPattern: string
): DisplayPattern {
  if (player.player_info.is_hidden) {
    return "private";
  }

  const isNoOverallBattle =
    (statsPattern === "pvp_all" && player.overall_stats.battles === 0) ||
    (userConfig.stats_pattern === "pvp_solo" &&
      player.overall_stats_solo.battles === 0);

  if (player.player_info.id === 0 || isNoOverallBattle) {
    return "nodata";
  }

  const isNoShipBattle =
    (statsPattern === "pvp_all" && player.ship_stats.battles === 0) ||
    (userConfig.stats_pattern === "pvp_solo" &&
      player.ship_stats_solo.battles === 0);

  if (isNoShipBattle) {
    return "noshipstats";
  }

  return "full";
}
</script>

{#if battle}
  <div class="m-2">
    <table class="table table-sm table-bordered table-text-color">
      {#each battle.teams as team}
        {@const basicColspan = components
          .filter(
            (it) =>
              it.category === "basic" && userConfig.displays.basic[it.name]
          )
          .reduce((a, it) => a + it.column, 1)}
        {@const shipColspan = components
          .filter(
            (it) => it.category === "ship" && userConfig.displays.ship[it.name]
          )
          .reduce((a, it) => a + it.column, 0)}
        {@const overallColspan = components
          .filter(
            (it) =>
              it.category === "overall" && userConfig.displays.overall[it.name]
          )
          .reduce((a, it) => a + it.column, 0)}
        <thead>
          <tr>
            <th colspan="{basicColspan}">基本情報</th>
            {#if shipColspan > 0}
              <th colspan="{shipColspan}">艦成績</th>
            {/if}
            {#if overallColspan > 0}
              <th colspan="{overallColspan}">総合成績</th>
            {/if}
          </tr>
          <tr>
            <th></th>

            {#each components as c}
              {#if userConfig.displays[c.category][c.name] === true}
                <th colspan="{c.column}">{Const.COLUMN_NAMES[c.name].min}</th>
              {/if}
            {/each}
          </tr>
        </thead>
        <tbody>
          {#each team.players as player}
            {@const displayPattern = decidePlayerDataPattern(
              player,
              userConfig.stats_pattern
            )}
            <tr>
              <!-- basics -->
              <BasicIsInAvg
                player="{player}"
                displayPattern="{displayPattern}"
                on:CheckPlayer
              />
              {#each components.filter((it) => it.category === "basic") as c}
                <svelte:component
                  this="{c.component}"
                  player="{player}"
                  displayPattern="{displayPattern}"
                  statsPattern="{userConfig.stats_pattern}"
                  on:UpdateAlertPlayer
                  on:RemoveAlertPlayer
                />
              {/each}

              <NoData
                shipColspan="{shipColspan}"
                overallColspan="{overallColspan}"
                displayPattern="{displayPattern}"
              />

              <!-- values -->
              {#each components.filter((it) => it.category === "ship" || it.category === "overall") as c}
                {#if userConfig.displays[c.category][c.name] === true}
                  <svelte:component
                    this="{c.component}"
                    player="{player}"
                    displayPattern="{displayPattern}"
                    statsPattern="{userConfig.stats_pattern}"
                    statsCatetory="{c.category}"
                  />
                {/if}
              {/each}
            </tr>
          {/each}
        </tbody>
      {/each}
    </table>
  </div>

  {#if summaryResult}
    <div class="mx-4 d-flex flex-row centerize">
      <table class="mx-2 table table-sm table-text-color w-auto">
        <tbody>
          <tr>
            <td class="td-string">日時</td>
            <td class="td-string">{battle.meta.date}</td>
          </tr>

          <tr>
            <td class="td-string">戦闘タイプ</td>
            <td class="td-string">{battle.meta.type}</td>
          </tr>

          <tr>
            <td class="td-string">マップ</td>
            <td class="td-string">{battle.meta.arena}</td>
          </tr>
        </tbody>
      </table>

      <table class="mx-2 table table-sm table-text-color w-auto">
        <thead>
          <tr>
            <th></th>
            <th colspan="{summaryResult.shipStatsCount}">艦成績</th>
            <th colspan="{summaryResult.overallStatsCount}">総合成績</th>
          </tr>
          <tr>
            <th></th>
            {#each summaryResult.labels as label}
              <th>{label}</th>
            {/each}
          </tr>
        </thead>
        <tbody>
          <tr>
            <td class="td-string">{battle.teams[0].name}</td>
            {#each summaryResult.friends as friend}
              <td class="td-number">{friend}</td>
            {/each}
          </tr>

          <tr>
            <td class="td-string">{battle.teams[1].name}</td>
            {#each summaryResult.enemies as enemy}
              <td class="td-number">{enemy}</td>
            {/each}
          </tr>

          <tr>
            <td class="td-string">差</td>
            {#each summaryResult.diffs as diff}
              <td class="td-number {diff.colorClass}">{diff.value}</td>
            {/each}
          </tr>
        </tbody>
      </table>
    </div>
  {/if}
{/if}
