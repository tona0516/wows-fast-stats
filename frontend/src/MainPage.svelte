<script lang="ts">
  import type { vo } from "wailsjs/go/models";
  import Const from "./Const";
  import Pr from "./Pr.svelte";
  import PlayerName from "./PlayerName.svelte";
  import ShipInfo from "./ShipInfo.svelte";
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
  import AvgTier from "./AvgTier.svelte";
  import ShipTypeRate from "./ShipTypeRate.svelte";
  import TierRate from "./TierRate.svelte";
  import OverallBattles from "./OverallBattles.svelte";
  import NoData from "./NoData.svelte";

  export let loadState: LoadState = "standby";
  export let latestHash: string = "";
  export let battle: vo.Battle;
  export let config: vo.UserConfig = Const.DEFAULT_USER_CONFIG;

  const basicComponents: { [key: string]: { header: string; component: any } } =
    {
      player_name: { header: "プレイヤー", component: PlayerName },
      ship_info: { header: "艦", component: ShipInfo },
    };

  const components: { [key: string]: { header: string; component: any } } = {
    pr: { header: "PR", component: Pr },
    ship_damage: { header: "S:Dmg", component: ShipDamage },
    ship_win_rate: { header: "S:勝率", component: ShipWinRate },
    ship_kd_rate: { header: "S:K/D", component: ShipKdRate },
    ship_win_survived_rate: {
      header: "S:勝利生存率",
      component: ShipWinSurvivedRate,
    },
    ship_lose_survived_rate: {
      header: "S:敗北生存率",
      component: ShipLoseSurvivedRate,
    },
    ship_exp: { header: "S:Exp", component: ShipExp },
    ship_battles: { header: "S:戦闘数", component: ShipBattles },
    player_damage: { header: "O:Dmg", component: OverallDamage },
    player_win_rate: { header: "O:勝率", component: OverallWinRate },
    player_kd_rate: { header: "O:K/D", component: OverallKdRate },
    player_win_survived_rate: {
      header: "O:勝利生存率",
      component: OverallWinSurvivedRate,
    },
    player_lose_survived_rate: {
      header: "O:敗北生存率",
      component: OverallLoseSurvivedRate,
    },
    player_exp: { header: "O:Exp", component: OverallExp },
    player_battles: { header: "O:戦闘数", component: OverallBattles },
    player_avg_tier: { header: "O:平均T", component: AvgTier },
    player_using_ship_type_rate: {
      header: "O:艦種割合",
      component: ShipTypeRate,
    },
    player_using_tier_rate: { header: "O:T別割合", component: TierRate },
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

    return `<td class="${value1Class}">${value1Round.toFixed(
      digits
    )}</td><td class="${value2Class}">${value2Round.toFixed(digits)}</td>`;
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
    {#each battle.teams as team}
      <table class="table table-sm">
        <thead>
          <tr>
            {#each Object.entries(basicComponents) as [k, v]}
              {#if config.displays[k]}
                <th>{v.header}</th>
              {/if}
            {/each}

            {#each Object.entries(components) as [k, v]}
              {#if config.displays[k]}
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
              {#each Object.entries(basicComponents) as [k, v]}
                <svelte:component
                  this={v.component}
                  {config}
                  {player}
                  {displayPattern}
                />
              {/each}

              <NoData {config} {displayPattern} />

              <!-- values -->
              {#each Object.entries(components) as [k, v]}
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
      </table>
    {/each}
  </div>

  <div class="mt-2 mx-4 d-flex flex-row">
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
            {#each battle.teams as team}
              <th>{team.name}</th>
            {/each}
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>PR</td>
            {@html renderTdForTeamSummary(battle, "personal_rating", 0)}
          </tr>
          <tr>
            <td>S:Dmg</td>
            {@html renderTdForTeamSummary(battle, "damage_by_ship", 0)}
          </tr>
          <tr>
            <td>S:勝率</td>
            {@html renderTdForTeamSummary(battle, "win_rate_by_ship", 1)}
          </tr>
          <tr>
            <td>S:K/D</td>
            {@html renderTdForTeamSummary(battle, "kd_rate_by_ship", 1)}
          </tr>
          <tr>
            <td>P:Dmg</td>
            {@html renderTdForTeamSummary(battle, "damage_by_player", 0)}
          </tr>
          <tr>
            <td>P:勝率</td>
            {@html renderTdForTeamSummary(battle, "win_rate_by_player", 1)}
          </tr>
          <tr>
            <td>P:K/D</td>
            {@html renderTdForTeamSummary(battle, "kd_rate_by_player", 1)}
          </tr>
        </tbody>
      </table>
    </div>
  </div>
{/if}
