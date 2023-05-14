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
import { Average, type AverageFactor } from "./Average";
import { ExcludePlayerIDs } from "../wailsjs/go/main/App.js";
import type { DisplayPattern } from "./DisplayPattern";
import ShipHitRate from "./ShipHitRate.svelte";

export let battle: vo.Battle;
export let config: vo.UserConfig = Const.DEFAULT_USER_CONFIG;
export let averageFactors: AverageFactor;
export let excludePlayerIDs: number[] = [];

const components = {
  basic: {
    player_name: BasicPlayerName,
    ship_info: BasicShipInfo,
  },
  ship: {
    pr: ShipPr,
    damage: ShipDamage,
    win_rate: ShipWinRate,
    kd_rate: ShipKdRate,
    exp: ShipExp,
    battles: ShipBattles,
    survived_rate: ShipSurvivedRate,
    hit_rate: ShipHitRate,
  },
  overall: {
    damage: OverallDamage,
    win_rate: OverallWinRate,
    kd_rate: OverallKdRate,
    exp: OverallExp,
    battles: OverallBattles,
    survived_rate: OverallSurvivedRate,
    avg_tier: OverallAvgTier,
    using_ship_type_rate: OverallShipTypeRate,
    using_tier_rate: OverallTierRate,
  },
};

function decidePlayerDataPattern(player: vo.Player): DisplayPattern {
  if (player.player_info.is_hidden) {
    return "private";
  }

  if (player.player_info.id === 0 || player.overall_stats.battles == 0) {
    return "nodata";
  }

  if (player.ship_stats.battles === 0) {
    return "noshipstats";
  }

  if (player.ship_stats.pr === -1) {
    return "nopr";
  }

  return "full";
}

function onCheckPlayer() {
  ExcludePlayerIDs().then((result) => {
    excludePlayerIDs = result;
    const average = new Average(battle);
    averageFactors = average.calc(result);
  });
}
</script>

{#if battle}
  <div class="m-2">
    <table class="table table-sm table-bordered table-text-color">
      {#each battle.teams as team}
        <thead>
          <tr>
            <th
              colspan="{Object.values(config.displays.basic).filter(
                (it) => it === true
              ).length}">基本情報</th
            >
            {#if Object.values(config.displays.ship).filter((it) => it === true).length !== 0}
              <th
                colspan="{Object.values(config.displays.ship).filter(
                  (it) => it === true
                ).length}">艦成績</th
              >
            {/if}
            {#if Object.values(config.displays.overall).filter((it) => it === true).length !== 0}
              <th
                colspan="{Object.values(config.displays.overall).filter(
                  (it) => it === true
                ).length}">総合成績</th
              >
            {/if}
          </tr>
          <tr>
            <th></th>
            {#each Object.keys(components.basic) as k}
              {#if config.displays.basic[k]}
                <th>{Const.COLUMN_NAMES[k].min}</th>
              {/if}
            {/each}

            {#each Object.keys(components.ship) as k}
              {#if config.displays.ship[k]}
                <th>{Const.COLUMN_NAMES[k].min}</th>
              {/if}
            {/each}

            {#each Object.keys(components.overall) as k}
              {#if config.displays.overall[k]}
                <th>{Const.COLUMN_NAMES[k].min}</th>
              {/if}
            {/each}
          </tr>
        </thead>
        <tbody>
          {#each team.players as player}
            {@const displayPattern = decidePlayerDataPattern(player)}
            <tr>
              <BasicIsInAvg
                player="{player}"
                excludePlayerIDs="{excludePlayerIDs}"
                displayPattern="{displayPattern}"
                on:onCheck="{onCheckPlayer}"
              />

              <!-- basics -->
              {#each Object.values(components.basic) as _v}
                <svelte:component this="{_v}" player="{player}" />
              {/each}

              <NoData config="{config}" displayPattern="{displayPattern}" />

              <!-- values -->
              {#each Object.values(components.ship) as _v}
                <svelte:component
                  this="{_v}"
                  config="{config}"
                  player="{player}"
                  displayPattern="{displayPattern}"
                />
              {/each}

              <!-- values -->
              {#each Object.values(components.overall) as _v}
                <svelte:component
                  this="{_v}"
                  config="{config}"
                  player="{player}"
                  displayPattern="{displayPattern}"
                />
              {/each}
            </tr>
          {/each}
        </tbody>
      {/each}
    </table>
  </div>

  {#if averageFactors}
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
            <th colspan="{averageFactors.shipStatsCount}">艦成績</th>
            <th colspan="{averageFactors.overallCount}">総合成績</th>
          </tr>
          <tr>
            <th></th>
            {#each averageFactors.labels as label}
              <th>{label}</th>
            {/each}
          </tr>
        </thead>
        <tbody>
          <tr>
            <td class="td-string">{battle.teams[0].name}</td>
            {#each averageFactors.friends as friend}
              <td class="td-number">{friend}</td>
            {/each}
          </tr>

          <tr>
            <td class="td-string">{battle.teams[1].name}</td>
            {#each averageFactors.enemies as enemy}
              <td class="td-number">{enemy}</td>
            {/each}
          </tr>

          <tr>
            <td class="td-string">差</td>
            {#each averageFactors.diffs as diff}
              <td class="td-number {diff.colorClass}">{diff.value}</td>
            {/each}
          </tr>
        </tbody>
      </table>
    </div>
  {/if}
{/if}
