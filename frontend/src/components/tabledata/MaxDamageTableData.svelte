<script lang="ts">
  import type { MaxDamageColumn } from "@libs/columns/MaxDamageColumn";
    import { storedDisplayPref } from "@libs/stores";
  import type { core } from "@wails/go/models";

  export let column: MaxDamageColumn;
  export let player: core.Player;

  const param = column.getDisplayValue(player);
</script>

{#if column.needsShow()}
  <td class="{$storedDisplayPref.showBoarder ? "border border-gray-500" : ""} p-1" style="background-color: {column.getBgColorCode(player)}">
    <div class="text-right">
      <div class="tooltip" data-tip={param.shipInfo?.name}>
        {param.damage}
      </div>
    </div>
  </td>
{/if}
