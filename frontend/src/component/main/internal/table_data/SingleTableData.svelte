<script lang="ts">
  import type { AbstractStatsColumn } from "src/lib/column/intetface/AbstractStatsColumn";
  import type { model } from "wailsjs/go/models";

  export let column: AbstractStatsColumn<string>;
  export let player: model.Player;
</script>

<td
  class={column.tdClass(player)}
  style="color: {column.textColorCode(
    player,
  )}; background-color: {column.bgColorCode(player)}"
>
  {column.displayValue(player)}
</td>
