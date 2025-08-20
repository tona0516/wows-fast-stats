<script lang="ts">
  import type { AbstractStatsColumn } from "@libs/columns/AbstractStatsColumn";
  import type { data } from "@wails/go/models";

  export let column: AbstractStatsColumn<string>;
  export let player: data.Player;
</script>

{#if column.needsShow()}
  <td
    class="p-1"
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
