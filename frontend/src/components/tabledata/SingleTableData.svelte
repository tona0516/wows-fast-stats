<script lang="ts">
  import type { AbstractStatsColumn } from "@libs/columns/AbstractStatsColumn";
  import { storedDisplayPref } from "@libs/stores";
  import type { core } from "@wails/go/models";

  export let column: AbstractStatsColumn<string>;
  export let player: core.Player;
</script>

{#if column.needsShow()}
  <td
    class="{$storedDisplayPref.isBorderVisible
      ? 'border border-gray-500'
      : ''} px-1 py-0.5"
    style="background-color: {column.getBgColorCode(player)?.raw}"
  >
    <div
      class={column.getCssClass() ?? ""}
      style="color: {column.getTextColorCode(player)?.raw};"
    >
      {column.getDisplayValue(player)}
    </div>
  </td>
{/if}
