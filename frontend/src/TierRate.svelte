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
{#if config.displays.player_using_tier_rate}
  {#if displayPattern === "full" || displayPattern === "nopr" || displayPattern === "noshipstats"}
    <td class="using_tier_rate">
      <div class="progress">
        <div
          class="progress-bar progress-bar-striped bg-primary"
          role="progressbar"
          style="width: {roundup(player.player_stats.using_tier_rate.low)}%"
        />
        <div
          class="progress-bar progress-bar-striped bg-info"
          role="progressbar"
          style="width: {roundup(player.player_stats.using_tier_rate.middle)}%"
        />
        <div
          class="progress-bar progress-bar-striped bg-success"
          role="progressbar"
          style="width: {roundup(player.player_stats.using_tier_rate.high)}%"
        />
      </div>
    </td>
  {/if}
{/if}
