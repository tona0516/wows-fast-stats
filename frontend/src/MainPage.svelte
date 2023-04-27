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
  import IconButton from "@smui/icon-button";
  import { createEventDispatcher } from "svelte";
  import { WindowReloadApp } from "../wailsjs/runtime/runtime";
  import DataTable, { Head, Body, Row, Cell } from "@smui/data-table";

  export let loadState: LoadState = "standby";
  export let latestHash: string = "";
  export let battle: vo.Battle;
  export let config: vo.UserConfig = Const.DEFAULT_USER_CONFIG;

  const dispatch = createEventDispatcher();

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

  function renderTdForTeamSummary(
    battle: vo.Battle,
    key: string,
    digits: number
  ): string {
    const value1Round =
      +battle["teams"][0]["team_average"][key].toFixed(digits);
    const value2Round =
      +battle["teams"][1]["team_average"][key].toFixed(digits);

    let value1Class = "";
    let value2Class = "";
    if (value1Round > value2Round) {
      value1Class = "higher";
      value2Class = "lower";
    }
    if (value1Round < value2Round) {
      value1Class = "lower";
      value2Class = "higher";
    }

    return `<Cell class="${value1Class}">${value1Round.toFixed(
      digits
    )}</Cell><Cell class="${value2Class}">${value2Round.toFixed(
      digits
    )}</Cell>`;
  }
</script>

<div style="display: flex; align-items: right;">
  <IconButton
    class="material-icons"
    on:click={() => dispatch("onClickScreenshot", null)}
    disabled={battle === undefined || loadState === "fetching"}
    >photo_camera</IconButton
  >
  <IconButton class="material-icons" on:click={() => WindowReloadApp()}
    >restart_alt</IconButton
  >
</div>

{#if latestHash !== ""}
  <DataTable class="table" stickyHeader>
    {#each battle.teams as team}
      <Head>
        <Row>
          <Cell
            colspan={Object.values(config.displays.basic).filter(
              (it) => it === true
            ).length}>基本情報</Cell
          >
          {#if Object.values(config.displays.ship).filter((it) => it === true).length !== 0}
            <Cell
              colspan={Object.values(config.displays.ship).filter(
                (it) => it === true
              ).length}>艦成績</Cell
            >
          {/if}
          {#if Object.values(config.displays.overall).filter((it) => it === true).length !== 0}
            <Cell
              colspan={Object.values(config.displays.overall).filter(
                (it) => it === true
              ).length}>総合成績</Cell
            >
          {/if}
        </Row>
        <Row>
          {#each Object.entries(components.basic) as [k, v]}
            {#if config.displays.basic[k]}
              <Cell>{v.header}</Cell>
            {/if}
          {/each}

          {#each Object.entries(components.ship) as [k, v]}
            {#if config.displays.ship[k]}
              <Cell>{v.header}</Cell>
            {/if}
          {/each}

          {#each Object.entries(components.overall) as [k, v]}
            {#if config.displays.overall[k]}
              <Cell>{v.header}</Cell>
            {/if}
          {/each}
        </Row>
      </Head>
      <Body>
        {#each team.players as player}
          {@const displayPattern = decidePlayerDataPattern(player)}
          <Row class={backgroundClass(player.ship_stats.personal_rating)}>
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
          </Row>
        {/each}
      </Body>
    {/each}
  </DataTable>

  <div style="display: flex; align-items: center;">
    <DataTable>
      <Body>
        <Row>
          <Cell>開始時刻</Cell>
          <Cell>{battle.meta.date}</Cell>
        </Row>
        <Row>
          <Cell>マップ</Cell>
          <Cell>{battle.meta.arena}</Cell>
        </Row>
        <Row>
          <Cell>戦闘タイプ</Cell>
          <Cell>{battle.meta.type}</Cell>
        </Row>
      </Body>
    </DataTable>
    <DataTable>
      <Head>
        <Row>
          <Cell />
          <Cell>{battle.teams[0].name}</Cell>
          <Cell>{battle.teams[1].name}</Cell>
        </Row>
      </Head>
      <Body>
        <Row>
          <Cell>艦:PR</Cell>
          <Cell>{battle.teams[0].team_average.personal_rating.toFixed(0)}</Cell>
          <Cell>{battle.teams[1].team_average.personal_rating.toFixed(0)}</Cell>
        </Row>
        <Row>
          <Cell>艦:Dmg</Cell>
          <Cell>{battle.teams[0].team_average.damage_by_ship.toFixed(0)}</Cell>
          <Cell>{battle.teams[1].team_average.damage_by_ship.toFixed(0)}</Cell>
        </Row>
        <Row>
          <Cell>艦:勝率</Cell>
          <Cell>{battle.teams[0].team_average.win_rate_by_ship.toFixed(1)}</Cell
          >
          <Cell>{battle.teams[1].team_average.win_rate_by_ship.toFixed(1)}</Cell
          >
        </Row>
        <Row>
          <Cell>艦:K/D</Cell>
          <Cell>{battle.teams[0].team_average.kd_rate_by_ship.toFixed(1)}</Cell>
          <Cell>{battle.teams[1].team_average.kd_rate_by_ship.toFixed(1)}</Cell>
        </Row>
        <Row>
          <Cell>総合:Dmg</Cell>
          <Cell>{battle.teams[0].team_average.damage_by_player.toFixed(0)}</Cell
          >
          <Cell>{battle.teams[1].team_average.damage_by_player.toFixed(0)}</Cell
          >
        </Row>
        <Row>
          <Cell>総合:勝率</Cell>
          <Cell
            >{battle.teams[0].team_average.win_rate_by_player.toFixed(1)}</Cell
          >
          <Cell
            >{battle.teams[1].team_average.win_rate_by_player.toFixed(1)}</Cell
          >
        </Row>
        <Row>
          <Cell>総合:K/D</Cell>
          <Cell
            >{battle.teams[0].team_average.kd_rate_by_player.toFixed(1)}</Cell
          >
          <Cell
            >{battle.teams[1].team_average.kd_rate_by_player.toFixed(1)}</Cell
          >
        </Row>
      </Body>
    </DataTable>
  </div>
{/if}
