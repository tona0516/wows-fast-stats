<script lang="ts">
  import type { MaxDamage } from "src/lib/column/model/MaxDamage";
  import type { data } from "wailsjs/go/models";

  export let column: MaxDamage;
  export let player: data.Player;

  const param = column.displayValue(player);
</script>

{#if column.shouldShow()}
  <td
    class="p-1"
    style="background-color: {column.getBackgroundColorCode(player) ?? ''}"
  >
    <div class="text-right">
      <div class="tooltip" data-tip={param.shipInfo?.name}>
        {param.damage}
      </div>
    </div>
  </td>
{/if}
