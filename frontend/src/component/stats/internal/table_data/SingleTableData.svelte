<script lang="ts">
  import type { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
  import type { data } from "wailsjs/go/models";

  export let column: AbstractStatsColumn<string>;
  export let player: data.Player;
</script>

{#if column.shouldShow()}
  <td
    class="p-1"
    style="background-color: {column.getBackgroundColorCode(player) ?? ''}"
  >
    <div
      class={column.getCssClass() ?? ""}
      style="color: {column.textColorCode(player)};"
    >
      {column.displayValue(player)}
    </div>
  </td>
{/if}
