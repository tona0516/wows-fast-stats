<script lang="ts">
  import type { ISingleColumn } from "src/lib/column/intetface/ISingleColumn";
  import type { domain } from "wailsjs/go/models";

  export let column: ISingleColumn;
  export let player: domain.Player;
</script>

<td
  class={column.getTdClass(player)}
  style="color: {column.getTextColorCode(player)}"
>
  {column.getDisplayValue(player)}
</td>
