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
  import AvgTier from "./OverallAvgTier.svelte";
  import ShipTypeRate from "./OverallShipTypeRate.svelte";
  import TierRate from "./OverallTierRate.svelte";
  import OverallBattles from "./OverallBattles.svelte";
  import NoData from "./NoData.svelte";

  export let loadState: LoadState = "standby";
  export let latestHash: string = "";
  export let battle: vo.Battle;
  export let config: vo.UserConfig = Const.DEFAULT_USER_CONFIG;

  const components = {
    basic: {
      player_name: { header: "プレイヤー", component: BasicPlayerName },
      ship_info: { header: "艦", component: BasicShipInfo },
    },
    ship: {
      pr: { header: "PR", component: ShipPr },
      damage: { header: "Dmg", component: ShipDamage },
      win_rate: { header: "勝率", component: ShipWinRate },
      kd_rate: { header: "K/D", component: ShipKdRate },
      win_survived_rate: {
        header: "勝利生存率",
        component: ShipWinSurvivedRate,
      },
      lose_survived_rate: {
        header: "敗北生存率",
        component: ShipLoseSurvivedRate,
      },
      exp: { header: "Exp", component: ShipExp },
      battles: { header: "戦闘数", component: ShipBattles },
    },
    overall: {
      damage: { header: "Dmg", component: OverallDamage },
      win_rate: { header: "勝率", component: OverallWinRate },
      kd_rate: { header: "K/D", component: OverallKdRate },
      win_survived_rate: {
        header: "勝利生存率",
        component: OverallWinSurvivedRate,
      },
      lose_survived_rate: {
        header: "敗北生存率",
        component: OverallLoseSurvivedRate,
      },
      exp: { header: "Exp", component: OverallExp },
      battles: { header: "戦闘数", component: OverallBattles },
      avg_tier: { header: "平均T", component: AvgTier },
      using_ship_type_rate: {
        header: "艦種割合",
        component: ShipTypeRate,
      },
      using_tier_rate: { header: "T別割合", component: TierRate },
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

    if (player.ship_stats.personal_rating === 0) {
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

  function buildTeamSummary(battle: vo.Battle): {
    name: string;
    value1: string;
    value2: string;
    diff: string;
    color_class: string;
  }[] {
    let result: {
      name: string;
      value1: string;
      value2: string;
      diff: string;
      color_class: string;
    }[] = [];
    Object.entries({
      personal_rating: {
        label: "艦:PR",
        digit: 0,
      },
      damage_by_ship: {
        label: "艦:Dmg",
        digit: 0,
      },
      win_rate_by_ship: {
        label: "艦:勝率",
        digit: 1,
      },
      kd_rate_by_ship: {
        label: "艦:K/D",
        digit: 1,
      },
      damage_by_player: {
        label: "総合:Dmg",
        digit: 0,
      },
      win_rate_by_player: {
        label: "総合:勝率",
        digit: 1,
      },
      kd_rate_by_player: {
        label: "総合:K/D",
        digit: 1,
      },
    }).forEach(([k, v]) => {
      const value1 = battle.teams[0].team_average[k];
      const value2 = battle.teams[1].team_average[k];
      const diff = value1 - value2;
      let colorClass = "";
      let sign = "";
      if (diff > 0) {
        sign = "+";
        colorClass = "higher";
      } else if (diff < 0) {
        colorClass = "lower";
      }

      result.push({
        name: v.label,
        value1: value1.toFixed(v.digit),
        value2: value2.toFixed(v.digit),
        diff: sign + diff.toFixed(v.digit),
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
    <table class="table table-sm table-bordered">
      {#each battle.teams as team}
        <thead>
          <tr>
            <th
              colspan={Object.values(config.displays.basic).filter(
                (it) => it === true
              ).length}>基本情報</th
            >
            {#if Object.values(config.displays.ship).filter((it) => it === true).length !== 0}
              <th
                colspan={Object.values(config.displays.ship).filter(
                  (it) => it === true
                ).length}>艦成績</th
              >
            {/if}
            {#if Object.values(config.displays.overall).filter((it) => it === true).length !== 0}
              <th
                colspan={Object.values(config.displays.overall).filter(
                  (it) => it === true
                ).length}>総合成績</th
              >
            {/if}
          </tr>
          <tr>
            {#each Object.entries(components.basic) as [k, v]}
              {#if config.displays.basic[k]}
                <th>{v.header}</th>
              {/if}
            {/each}

            {#each Object.entries(components.ship) as [k, v]}
              {#if config.displays.ship[k]}
                <th>{v.header}</th>
              {/if}
            {/each}

            {#each Object.entries(components.overall) as [k, v]}
              {#if config.displays.overall[k]}
                <th>{v.header}</th>
              {/if}
            {/each}
          </tr>
        </thead>
        <tbody>
          {#each team.players as player}
            {@const displayPattern = decidePlayerDataPattern(player)}
            <tr class={backgroundClass(player.ship_stats.personal_rating)}>
              <!-- basics -->
              {#each Object.entries(components.basic) as [k, v]}
                <svelte:component
                  this={v.component}
                  {config}
                  {player}
                  {displayPattern}
                />
              {/each}

              <NoData {config} {displayPattern} />

              <!-- values -->
              {#each Object.entries(components.ship) as [k, v]}
                <svelte:component
                  this={v.component}
                  {config}
                  {player}
                  {displayPattern}
                />
              {/each}

              <!-- values -->
              {#each Object.entries(components.overall) as [k, v]}
                <svelte:component
                  this={v.component}
                  {config}
                  {player}
                  {displayPattern}
                />
              {/each}
            </tr>
          {/each}
        </tbody>
      {/each}
    </table>
  </div>

  <div class="mt-2 mx-4 d-flex flex-row centerize">
    <div class="mx-2">
      <table class="table table-sm w-auto">
        <tbody>
          <tr>
            <td>開始時刻</td>
            <td>{battle.meta.date}</td>
          </tr>
          <tr>
            <td>マップ</td>
            <td>{battle.meta.arena}</td>
          </tr>
          <tr>
            <td>戦闘タイプ</td>
            <td>{battle.meta.type}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="mx-2">
      <table class="table table-sm w-auto">
        <thead>
          <tr>
            <th />
            <th>{battle.teams[0].name}</th>
            <th>差</th>
            <th>{battle.teams[1].name}</th>
          </tr>
        </thead>
        <tbody>
          {#each buildTeamSummary(battle) as row}
            <tr>
              <td>{row.name}</td>
              <td>{row.value1}</td>
              <td class={row.color_class}>{row.diff}</td>
              <td>{row.value2}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
{/if}
