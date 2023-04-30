<script lang="ts">
import type { vo } from "wailsjs/go/models";
import Const from "./Const";
import ShipPr from "./ShipPr.svelte";
import BasicPlayerName from "./BasicPlayerName.svelte";
import BasicShipInfo from "./BasicShipInfo.svelte";
import ShipDamage from "./ShipDamage.svelte";
import ShipWinRate from "./ShipWinRate.svelte";
import ShipKdRate from "./ShipKdRate.svelte";
import ShipWinSurvivedRate from "./ShipWinSurvivedRate.svelte";
import ShipLoseSurvivedRate from "./ShipLoseSurvivedRate.svelte";
import ShipExp from "./ShipExp.svelte";
import ShipBattles from "./ShipBattles.svelte";
import OverallDamage from "./OverallDamage.svelte";
import OverallWinRate from "./OverallWinRate.svelte";
import OverallKdRate from "./OverallKdRate.svelte";
import OverallWinSurvivedRate from "./OverallWinSurvivedRate.svelte";
import OverallLoseSurvivedRate from "./OverallLoseSurvivedRate.svelte";
import OverallExp from "./OverallExp.svelte";
import OverallAvgTier from "./OverallAvgTier.svelte";
import OverallShipTypeRate from "./OverallShipTypeRate.svelte";
import OverallTierRate from "./OverallTierRate.svelte";
import OverallBattles from "./OverallBattles.svelte";
import NoData from "./NoData.svelte";
import { LogDebug } from "../wailsjs/runtime/runtime";

export let loadState: LoadState = "standby";
export let latestHash: string = "";
export let battle: vo.Battle;
export let config: vo.UserConfig = Const.DEFAULT_USER_CONFIG;

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
    win_survived_rate: ShipWinSurvivedRate,
    lose_survived_rate: ShipLoseSurvivedRate,
    exp: ShipExp,
    battles: ShipBattles,
  },
  overall: {
    damage: OverallDamage,
    win_rate: OverallWinRate,
    kd_rate: OverallKdRate,
    win_survived_rate: OverallWinSurvivedRate,
    lose_survived_rate: OverallLoseSurvivedRate,
    exp: OverallExp,
    battles: OverallBattles,
    avg_tier: OverallAvgTier,
    using_ship_type_rate: OverallShipTypeRate,
    using_tier_rate: OverallTierRate,
  },
};

function decidePlayerDataPattern(player: vo.Player): DisplayPattern {
  if (player.player_info.is_hidden) {
    return "private";
  }

  if (player.player_info.id === 0 || player.player_stats.battles == 0) {
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

function backgroundClass(personalRating: number): string {
  switch (true) {
    case personalRating == 0:
      return "";
    case personalRating < 750:
      return "bad";
    case personalRating < 1100:
      return "below-average";
    case personalRating < 1350:
      return "average";
    case personalRating < 1550:
      return "good";
    case personalRating < 1750:
      return "very-good";
    case personalRating < 2100:
      return "great";
    case personalRating < 2450:
      return "unicum";
    case personalRating >= 2450:
      return "super-unicum";
    default:
      return "";
  }
}

function buildTeamSummary(comp: vo.Comparision): {
  label: string;
  friend: string;
  enemy: string;
  diff: string;
  color_class: string;
}[] {
  let result: {
    label: string;
    friend: string;
    enemy: string;
    diff: string;
    color_class: string;
  }[] = [];

  let resultKeys: { key1: string; key2: string }[] = [];

  Object.keys(comp.ship).forEach((it) => {
    resultKeys.push({ key1: "ship", key2: it });
  });

  Object.keys(comp.overall).forEach((it) => {
    resultKeys.push({ key1: "overall", key2: it });
  });

  resultKeys.forEach((it) => {
    const between: vo.Between = comp[it.key1][it.key2];

    let colorClass = "";
    let sign = "";
    if (between.diff > 0) {
      sign = "+";
      colorClass = "higher";
    } else if (between.diff < 0) {
      colorClass = "lower";
    }

    result.push({
      label:
        Const.COLUMN_NAMES[it.key1].min + ":" + Const.COLUMN_NAMES[it.key2].min,
      friend: between.friend.toFixed(Const.DIGITS[it.key2]),
      enemy: between.enemy.toFixed(Const.DIGITS[it.key2]),
      diff: sign + between.diff.toFixed(Const.DIGITS[it.key2]),
      color_class: colorClass,
    });
  });

  return result;
}
</script>

{#if loadState === "fetching"}
  <div class="d-flex justify-content-center m-3">
    <div class="spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
  </div>
{/if}

{#if latestHash !== ""}
  <div class="mt-2 mx-4">
    <span>
      {battle.meta.date}
      {battle.meta.arena}
      {battle.meta.type}
    </span>
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
            {#each Object.entries(components.basic) as [k, v]}
              {#if config.displays.basic[k]}
                <th>{Const.COLUMN_NAMES[k].min}</th>
              {/if}
            {/each}

            {#each Object.entries(components.ship) as [k, v]}
              {#if config.displays.ship[k]}
                <th>{Const.COLUMN_NAMES[k].min}</th>
              {/if}
            {/each}

            {#each Object.entries(components.overall) as [k, v]}
              {#if config.displays.overall[k]}
                <th>{Const.COLUMN_NAMES[k].min}</th>
              {/if}
            {/each}
          </tr>
        </thead>
        <tbody>
          {#each team.players as player}
            {@const displayPattern = decidePlayerDataPattern(player)}
            <tr class="{backgroundClass(player.ship_stats.pr)}">
              <!-- basics -->
              {#each Object.entries(components.basic) as [k, v]}
                <svelte:component
                  this="{v}"
                  config="{config}"
                  player="{player}"
                  displayPattern="{displayPattern}"
                />
              {/each}

              <NoData config="{config}" displayPattern="{displayPattern}" />

              <!-- values -->
              {#each Object.entries(components.ship) as [k, v]}
                <svelte:component
                  this="{v}"
                  config="{config}"
                  player="{player}"
                  displayPattern="{displayPattern}"
                />
              {/each}

              <!-- values -->
              {#each Object.entries(components.overall) as [k, v]}
                <svelte:component
                  this="{v}"
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

  <div class="mt-2 mx-4 d-flex flex-row centerize">
    <table class="table table-sm table-text-color w-auto">
      <thead>
        <tr>
          <th></th>
          <th>{battle.teams[0].name}</th>
          <th>差</th>
          <th>{battle.teams[1].name}</th>
        </tr>
      </thead>
      <tbody>
        {#each buildTeamSummary(battle.comparision) as row}
          <tr>
            <td class="td-string">{row.label}</td>
            <td class="text-center td-number">{row.friend}</td>
            <td class="text-center td-number {row.color_class}">{row.diff}</td>
            <td class="text-center td-number">{row.enemy}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
{/if}
