<script lang="ts">
  import type { ISingleColumn } from "src/lib/column/intetface/ISingleColumn";
  import type { model } from "wailsjs/go/models";

  export let column: ISingleColumn;
  export let player: model.Player;
</script>

<td
  class={column.getTdClass(player)}
  style="color: {column.getTextColorCode(player)}"
>
  {column.getDisplayValue(player)}
</td>
