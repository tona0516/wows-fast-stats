<script lang="ts">
  import iconCV from "./assets/images/icon-cv.png";
  import iconBB from "./assets/images/icon-bb.png";
  import iconCL from "./assets/images/icon-cl.png";
  import iconDD from "./assets/images/icon-dd.png";
  import iconSS from "./assets/images/icon-ss.png";
  import iconNone from "./assets/images/icon-none.png";
  import { BrowserOpenURL } from "../wailsjs/runtime";

  type StatsType = "ship" | "player";
  type State = "standby" | "fetching" | "error";

  export let state: State = "standby";
  export let latestHash = "";
  export let teams = [];

  function isValidStatsValue(player: any, statsType: StatsType) {
    let battles: number;
    switch (statsType) {
      case "ship":
        battles = player.player_ship_stats.battles;
        break;
      case "player":
        battles = player.player_player_stats.battles;
        break;
      default:
        battles = 0;
        break;
    }
    return !player.player_player_info.is_hidden && battles > 0;
  }

  function backgroundClassForTr(personalRating: number): string {
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

  function numberForDisplay(value: number): string {
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

  function shipIconForDisplay(shipType: string): string {
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
</script>

{#if state === "fetching"}
  <div class="d-flex justify-content-center m-3">
    <div class="spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
  </div>
{/if}

{#if latestHash !== ""}
  <div class="mt-1 mx-3">
    {#each teams as team}
      <span>
        {team.name} 艦別:{team.win_rate_by_ship.toFixed(1)}% 全体:{team.win_rate_by_player.toFixed(
          1
        )}%
      </span>
      <table class="table table-sm">
        <thead>
          <tr>
            <th>プレイヤー</th>
            <th>艦</th>
            <th>PR</th>
            <th>Dmg(艦)</th>
            <th>勝率(艦)</th>
            <th>K/D(艦)</th>
            <th>戦闘数(艦)</th>
            <th>Dmg</th>
            <th>勝率</th>
            <th>K/D</th>
            <th>戦闘数</th>
            <th>平均T</th>
          </tr>
        </thead>
        <tbody>
          {#each team.players as player}
            <tr
              class={backgroundClassForTr(
                player.player_ship_stats.personal_rating
              )}
            >
              <td class="name omit">
                <a
                  href="#"
                  on:click={BrowserOpenURL(player.player_player_info.stats_url)}
                >
                  {#if player.player_player_info.clan}
                    [{player.player_player_info.clan}]{player.player_player_info
                      .name}
                  {:else}
                    {player.player_player_info.name}
                  {/if}
                </a>
              </td>

              <td class="name omit">
                <a
                  href="#"
                  on:click={BrowserOpenURL(player.player_ship_info.stats_url)}
                >
                  <div class="aligner">
                    <img
                      alt=""
                      src={shipIconForDisplay(player.player_ship_info.type)}
                      class="icon-scale"
                    />
                    <div class="omit">
                      {numberForDisplay(player.player_ship_info.tier)}
                      {player.player_ship_info.name}
                    </div>
                  </div>
                </a>
              </td>

              {#if isValidStatsValue(player, "ship")}
                <td class="pr"
                  >{player.player_ship_stats.personal_rating.toFixed(0)}</td
                >
              {:else}
                <td class="pr" />
              {/if}

              {#if isValidStatsValue(player, "ship")}
                <td class="damage"
                  >{player.player_ship_stats.avg_damage.toFixed(0)}</td
                >
              {:else}
                <td class="damage" />
              {/if}

              {#if isValidStatsValue(player, "ship")}
                <td class="win"
                  >{player.player_ship_stats.win_rate.toFixed(1)}</td
                >
              {:else}
                <td class="win" />
              {/if}

              {#if isValidStatsValue(player, "ship")}
                <td class="kd">{player.player_ship_stats.kd_rate.toFixed(1)}</td
                >
              {:else}
                <td class="kd" />
              {/if}

              {#if isValidStatsValue(player, "ship")}
                <td class="battles">{player.player_ship_stats.battles}</td>
              {:else}
                <td class="battles" />
              {/if}

              {#if isValidStatsValue(player, "player")}
                <td class="damage"
                  >{player.player_player_stats.avg_damage.toFixed(0)}</td
                >
              {:else}
                <td class="damage" />
              {/if}
              {#if isValidStatsValue(player, "player")}
                <td class="win"
                  >{player.player_player_stats.win_rate.toFixed(1)}</td
                >
              {:else}
                <td class="win" />
              {/if}
              {#if isValidStatsValue(player, "player")}
                <td class="kd"
                  >{player.player_player_stats.kd_rate.toFixed(1)}</td
                >
              {:else}
                <td class="kd" />
              {/if}
              {#if isValidStatsValue(player, "player")}
                <td class="battles">{player.player_player_stats.battles}</td>
              {:else}
                <td class="battles" />
              {/if}
              {#if isValidStatsValue(player, "player")}
                <td class="avg-tier"
                  >{player.player_player_stats.avg_tier.toFixed(1)}</td
                >
              {:else}
                <td class="avg-tier" />
              {/if}
            </tr>
          {/each}
        </tbody>
      </table>
    {/each}
  </div>
{/if}
