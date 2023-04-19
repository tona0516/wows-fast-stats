<script lang="ts">
  import iconCV from "./assets/images/icon-cv.png";
  import iconBB from "./assets/images/icon-bb.png";
  import iconCL from "./assets/images/icon-cl.png";
  import iconDD from "./assets/images/icon-dd.png";
  import iconSS from "./assets/images/icon-ss.png";
  import iconNone from "./assets/images/icon-none.png";
  import { BrowserOpenURL } from "../wailsjs/runtime/runtime";
  import type { vo } from "wailsjs/go/models";

  export let loadState: LoadState = "standby";
  export let latestHash: string = "";
  export let teams: vo.Team[] = [];

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
</script>

{#if loadState === "fetching"}
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
              <td class="name omit">
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

              {#if dataPattern === "private"}
                <td colspan="10">PRIVATE</td>
              {:else if dataPattern === "nodata"}
                <td colspan="10">NO DATA</td>
              {/if}

              <!-- personal rating -->
              {#if dataPattern === "full"}
                <td class="pr">
                  {player.player_ship_stats.personal_rating.toFixed(0)}
                </td>
              {:else if dataPattern === "noshipstats" || dataPattern === "nopr"}
                <td />
              {/if}

              <!-- ship avg damage -->
              {#if dataPattern === "full" || dataPattern === "nopr"}
                <td class="damage">
                  {player.player_ship_stats.avg_damage.toFixed(0)}
                </td>
              {:else if dataPattern === "noshipstats"}
                <td />
              {/if}

              <!-- ship win rate -->
              {#if dataPattern === "full" || dataPattern === "nopr"}
                <td class="win">
                  {player.player_ship_stats.win_rate.toFixed(1)}
                </td>
              {:else if dataPattern === "noshipstats"}
                <td />
              {/if}

              <!-- ship kd rate -->
              {#if dataPattern === "full" || dataPattern === "nopr"}
                <td class="kd">
                  {player.player_ship_stats.kd_rate.toFixed(1)}
                </td>
              {:else if dataPattern === "noshipstats"}
                <td />
              {/if}

              <!-- ship battles -->
              {#if dataPattern === "full" || dataPattern === "nopr"}
                <td class="battles">
                  {player.player_ship_stats.battles}
                </td>
              {:else if dataPattern === "noshipstats"}
                <td />
              {/if}

              <!-- player avg damage -->
              {#if dataPattern === "noshipstats" || dataPattern === "full" || dataPattern === "nopr"}
                <td class="damage">
                  {player.player_player_stats.avg_damage.toFixed(0)}
                </td>
              {/if}

              <!-- player win rate -->
              {#if dataPattern === "noshipstats" || dataPattern === "full" || dataPattern === "nopr"}
                <td class="win">
                  {player.player_player_stats.win_rate.toFixed(1)}
                </td>
              {/if}

              <!-- player kd rate -->
              {#if dataPattern === "noshipstats" || dataPattern === "full" || dataPattern === "nopr"}
                <td class="kd">
                  {player.player_player_stats.kd_rate.toFixed(1)}
                </td>
              {/if}

              <!-- player battles -->
              {#if dataPattern === "noshipstats" || dataPattern === "full" || dataPattern === "nopr"}
                <td class="battles">
                  {player.player_player_stats.battles}
                </td>
              {/if}

              <!-- avg tier -->
              {#if dataPattern === "noshipstats" || dataPattern === "full" || dataPattern === "nopr"}
                <td class="avg-tier">
                  {player.player_player_stats.avg_tier.toFixed(1)}
                </td>
              {/if}
            </tr>
          {/each}
        </tbody>
      </table>
    {/each}
  </div>
{/if}
