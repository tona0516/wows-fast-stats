<script lang="ts">
  import iconCV from "./assets/images/icon-cv.png";
  import iconBB from "./assets/images/icon-bb.png";
  import iconCL from "./assets/images/icon-cl.png";
  import iconDD from "./assets/images/icon-dd.png";
  import iconSS from "./assets/images/icon-ss.png";
  import iconNone from "./assets/images/icon-none.png";
  import { BrowserOpenURL } from "../wailsjs/runtime/runtime";
  import type { vo } from "wailsjs/go/models";
  import Const from "./Const";

  export let loadState: LoadState = "standby";
  export let latestHash: string = "";
  export let battle: vo.Battle;
  export let config: vo.UserConfig = Const.DEFAULT_USER_CONFIG;

  /**
   * private: hidden player.
   * nodata: invalid player(bot/deleted account) or 0 for all random battle.
   * noshipstats: 0 in for random battle with the ship.
   * nopr: not exist expected value in numbers api.
   * full: all values exists.
   */
  type PlayerDataPattern =
    | "private"
    | "nodata"
    | "noshipstats"
    | "nopr"
    | "full";

  function decidePlayerDataPattern(player: vo.Player): PlayerDataPattern {
    if (player.player_player_info.is_hidden) {
      return "private";
    }

    if (
      player.player_player_info.id === 0 ||
      player.player_player_stats.battles == 0
    ) {
      return "nodata";
    }

    if (player.player_ship_stats.battles === 0) {
      return "noshipstats";
    }

    if (player.player_ship_stats.personal_rating === 0) {
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

  function tierString(value: number): string {
    if (value === 11) return "★";

    const decimal = [10, 9, 5, 4, 1];
    const romanNumeral = ["X", "IX", "V", "IV", "I"];

    let romanized = "";

    for (var i = 0; i < decimal.length; i++) {
      while (decimal[i] <= value) {
        romanized += romanNumeral[i];
        value -= decimal[i];
      }
    }
    return romanized;
  }

  function shipIcon(shipType: string): string {
    switch (shipType) {
      case "AirCarrier":
        return iconCV;
      case "Battleship":
        return iconBB;
      case "Cruiser":
        return iconCL;
      case "Destroyer":
        return iconDD;
      case "Submarine":
        return iconSS;
      default:
        return iconNone;
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

  function countDisplays(config: vo.UserConfig): number {
    return Object.values(config.displays).filter((it) => it === true).length;
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
    <span> プレイヤー統計 </span>
    {#each battle.teams as team}
      <table class="table table-sm">
        <thead>
          <tr>
            <th>プレイヤー</th>
            <th class="border-right">艦</th>
            {#if config.displays.pr}
              <th>PR</th>
            {/if}
            {#if config.displays.ship_damage}
              <th>S:Dmg</th>
            {/if}
            {#if config.displays.ship_win_rate}
              <th>S:勝率</th>
            {/if}
            {#if config.displays.ship_kd_rate}
              <th>S:K/D</th>
            {/if}
            {#if config.displays.ship_win_survived_rate}
              <th>S:勝利生存率</th>
            {/if}
            {#if config.displays.ship_lose_survived_rate}
              <th>S:敗北生存率</th>
            {/if}
            {#if config.displays.ship_exp}
              <th>S:Exp</th>
            {/if}
            {#if config.displays.ship_battles}
              <th>S:戦闘数</th>
            {/if}
            {#if config.displays.player_damage}
              <th>P:Dmg</th>
            {/if}
            {#if config.displays.player_win_rate}
              <th>P:勝率</th>
            {/if}
            {#if config.displays.player_kd_rate}
              <th>P:K/D</th>
            {/if}
            {#if config.displays.player_win_survived_rate}
              <th>P:勝利生存率</th>
            {/if}
            {#if config.displays.player_lose_survived_rate}
              <th>P:敗北生存率</th>
            {/if}
            {#if config.displays.player_exp}
              <th>P:Exp</th>
            {/if}
            {#if config.displays.player_battles}
              <th>P:戦闘数</th>
            {/if}
            {#if config.displays.player_avg_tier}
              <th>P:平均T</th>
            {/if}
          </tr>
        </thead>
        <tbody>
          {#each team.players as player}
            {@const dataPattern = decidePlayerDataPattern(player)}
            <tr
              class={backgroundClass(player.player_ship_stats.personal_rating)}
            >
              <!-- player name -->
              <td class="name omit">
                <!-- svelte-ignore a11y-invalid-attribute -->
                <a
                  href="#"
                  on:click={() =>
                    BrowserOpenURL(player.player_player_info.stats_url)}
                >
                  {#if player.player_player_info.clan}
                    [{player.player_player_info.clan}]{player.player_player_info
                      .name}
                  {:else}
                    {player.player_player_info.name}
                  {/if}
                </a>
              </td>

              <!-- ship info -->
              <td class="name omit border-right">
                <!-- svelte-ignore a11y-invalid-attribute -->
                <a
                  href="#"
                  on:click={() =>
                    BrowserOpenURL(player.player_ship_info.stats_url)}
                >
                  <div class="horizontal">
                    <img
                      alt=""
                      src={shipIcon(player.player_ship_info.type)}
                      class="icon-scale"
                    />
                    <div class="omit">
                      {tierString(player.player_ship_info.tier)}
                      {player.player_ship_info.name}
                    </div>
                  </div>
                </a>
              </td>

              <!-- join column case -->
              {#if countDisplays(config) > 0}
                {#if dataPattern === "private"}
                  <td colspan={countDisplays(config)}>PRIVATE</td>
                {:else if dataPattern === "nodata"}
                  <td colspan={countDisplays(config)}>NO DATA</td>
                {/if}
              {/if}

              <!-- personal rating -->
              {#if config.displays.pr}
                {#if dataPattern === "full"}
                  <td class="pr">
                    {player.player_ship_stats.personal_rating.toFixed(0)}
                  </td>
                {:else if dataPattern === "noshipstats" || dataPattern === "nopr"}
                  <td class="pr" />
                {/if}
              {/if}

              <!-- ship avg damage -->
              {#if config.displays.ship_damage}
                {#if dataPattern === "full" || dataPattern === "nopr"}
                  <td class="damage">
                    {player.player_ship_stats.avg_damage.toFixed(0)}
                  </td>
                {:else if dataPattern === "noshipstats"}
                  <td class="damage" />
                {/if}
              {/if}

              <!-- ship win rate -->
              {#if config.displays.ship_win_rate}
                {#if dataPattern === "full" || dataPattern === "nopr"}
                  <td class="win">
                    {player.player_ship_stats.win_rate.toFixed(1)}
                  </td>
                {:else if dataPattern === "noshipstats"}
                  <td class="win" />
                {/if}
              {/if}

              <!-- ship kd rate -->
              {#if config.displays.ship_kd_rate}
                {#if dataPattern === "full" || dataPattern === "nopr"}
                  <td class="kd">
                    {player.player_ship_stats.kd_rate.toFixed(1)}
                  </td>
                {:else if dataPattern === "noshipstats"}
                  <td class="kd" />
                {/if}
              {/if}

              <!-- ship win survived rate -->
              {#if config.displays.ship_win_survived_rate}
                {#if dataPattern === "full" || dataPattern === "nopr"}
                  <td class="survived-rate">
                    {player.player_ship_stats.win_survived_rate.toFixed(1)}
                  </td>
                {:else if dataPattern === "noshipstats"}
                  <td class="survived-rate" />
                {/if}
              {/if}

              <!-- ship lose survived rate -->
              {#if config.displays.ship_lose_survived_rate}
                {#if dataPattern === "full" || dataPattern === "nopr"}
                  <td class="survived-rate">
                    {player.player_ship_stats.lose_survived_rate.toFixed(1)}
                  </td>
                {:else if dataPattern === "noshipstats"}
                  <td class="survived-rate" />
                {/if}
              {/if}

              <!-- ship exp -->
              {#if config.displays.ship_exp}
                {#if dataPattern === "full" || dataPattern === "nopr"}
                  <td class="exp">
                    {player.player_ship_stats.exp.toFixed(0)}
                  </td>
                {:else if dataPattern === "noshipstats"}
                  <td class="exp" />
                {/if}
              {/if}

              <!-- ship battles -->
              {#if config.displays.ship_battles}
                {#if dataPattern === "full" || dataPattern === "nopr"}
                  <td class="battles">
                    {player.player_ship_stats.battles}
                  </td>
                {:else if dataPattern === "noshipstats"}
                  <td class="battles" />
                {/if}
              {/if}

              <!-- player avg damage -->
              {#if config.displays.player_damage}
                {#if dataPattern === "noshipstats" || dataPattern === "full" || dataPattern === "nopr"}
                  <td class="damage">
                    {player.player_player_stats.avg_damage.toFixed(0)}
                  </td>
                {/if}
              {/if}

              <!-- player win rate -->
              {#if config.displays.player_win_rate}
                {#if dataPattern === "noshipstats" || dataPattern === "full" || dataPattern === "nopr"}
                  <td class="win">
                    {player.player_player_stats.win_rate.toFixed(1)}
                  </td>
                {/if}
              {/if}

              <!-- player kd rate -->
              {#if config.displays.player_kd_rate}
                {#if dataPattern === "noshipstats" || dataPattern === "full" || dataPattern === "nopr"}
                  <td class="kd">
                    {player.player_player_stats.kd_rate.toFixed(1)}
                  </td>
                {/if}
              {/if}

              <!-- player win survived rate -->
              {#if config.displays.player_win_survived_rate}
                {#if dataPattern === "full" || dataPattern === "nopr"}
                  <td class="survived-rate">
                    {player.player_player_stats.win_survived_rate.toFixed(1)}
                  </td>
                {:else if dataPattern === "noshipstats"}
                  <td class="survived-rate" />
                {/if}
              {/if}

              <!-- player lose survived rate -->
              {#if config.displays.player_lose_survived_rate}
                {#if dataPattern === "full" || dataPattern === "nopr"}
                  <td class="survived-rate">
                    {player.player_player_stats.lose_survived_rate.toFixed(1)}
                  </td>
                {:else if dataPattern === "noshipstats"}
                  <td class="survived-rate" />
                {/if}
              {/if}

              <!-- player exp -->
              {#if config.displays.player_exp}
                {#if dataPattern === "full" || dataPattern === "nopr"}
                  <td class="exp">
                    {player.player_player_stats.exp.toFixed(0)}
                  </td>
                {:else if dataPattern === "noshipstats"}
                  <td class="exp" />
                {/if}
              {/if}

              <!-- player battles -->
              {#if config.displays.player_battles}
                {#if dataPattern === "noshipstats" || dataPattern === "full" || dataPattern === "nopr"}
                  <td class="battles">
                    {player.player_player_stats.battles}
                  </td>
                {/if}
              {/if}

              <!-- avg tier -->
              {#if config.displays.player_avg_tier}
                {#if dataPattern === "noshipstats" || dataPattern === "full" || dataPattern === "nopr"}
                  <td class="avg-tier">
                    {player.player_player_stats.avg_tier.toFixed(1)}
                  </td>
                {/if}
              {/if}
            </tr>
          {/each}
        </tbody>
      </table>
    {/each}
  </div>

  <div class="mt-2 row">
    <div class="col">
      <span> 戦闘サマリー </span>
      <table class="table table-sm w-auto" align="center">
        <tr>
          <th>項目</th>
          <th>値</th>
        </tr>
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
    <div class="col">
      <span> チームサマリー </span>
      <table class="table table-sm w-auto" align="center">
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
