<script lang="ts">
  import type { AbstractSingleColumn } from "src/lib/column/intetface/AbstractSingleColumn";
  import type { domain } from "wailsjs/go/models";

  export let column: AbstractSingleColumn<any>;
  export let player: domain.Player;
</script>

<td
  class={column.tdClass(player)}
  style="color: {column.textColorCode(player)}"
>
  {column.displayValue(player)}
</td>
