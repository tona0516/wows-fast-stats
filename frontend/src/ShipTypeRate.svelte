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

{#if config.displays.player_using_ship_type_rate}
  {#if displayPattern === "full" || displayPattern === "nopr" || displayPattern === "noshipstats"}
    <td class="using_ship_type_rate">
      <div class="progress">
        <div
          class="progress-bar progress-bar-striped bg-primary"
          role="progressbar"
          style="width: {roundup(player.player_stats.using_ship_type_rate.ss)}%"
        />
        <div
          class="progress-bar progress-bar-striped bg-info"
          role="progressbar"
          style="width: {roundup(player.player_stats.using_ship_type_rate.dd)}%"
        />
        <div
          class="progress-bar progress-bar-striped bg-success"
          role="progressbar"
          style="width: {roundup(player.player_stats.using_ship_type_rate.cl)}%"
        />
        <div
          class="progress-bar progress-bar-striped bg-warning"
          role="progressbar"
          style="width: {roundup(player.player_stats.using_ship_type_rate.bb)}%"
        />
        <div
          class="progress-bar progress-bar-striped bg-danger"
          role="progressbar"
          style="width: {roundup(player.player_stats.using_ship_type_rate.cv)}%"
        />
      </div>
    </td>
  {/if}
{/if}
