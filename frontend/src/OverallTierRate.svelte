<script lang="ts">
  import type { vo } from "wailsjs/go/models";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;

  function roundup(value: number, digit: number = 0): string {
    const digitVal = Math.pow(10, digit);
    const result = String(Math.ceil(value * digitVal) / digitVal);
    return result;
  }
</script>

<!-- using tier rate -->
{#if config.displays.overall.using_tier_rate}
  {#if displayPattern === "full" || displayPattern === "nopr" || displayPattern === "noshipstats"}
    <td class="td-graph">
      <div class="progress">
        <div
          class="progress-bar progress-bar-striped bar-low"
          role="progressbar"
          style="width: {roundup(player.player_stats.using_tier_rate.low)}%"
        />
        <div
          class="progress-bar progress-bar-striped bar-middle"
          role="progressbar"
          style="width: {roundup(player.player_stats.using_tier_rate.middle)}%"
        />
        <div
          class="progress-bar progress-bar-striped bar-high"
          role="progressbar"
          style="width: {roundup(player.player_stats.using_tier_rate.high)}%"
        />
      </div>
    </td>
  {/if}
{/if}
