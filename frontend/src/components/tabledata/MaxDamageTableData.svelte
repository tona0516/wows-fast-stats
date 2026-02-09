<script lang="ts">
  import ExternalLink from "@components/ExternalLink.svelte";
  import type { MaxDamageColumn } from "@libs/columns/MaxDamageColumn";
  import { storedDisplayPref } from "@libs/stores";
  import type { core } from "@wails/go/models";

  export let column: MaxDamageColumn;
  export let player: core.Player;

  const param = column.getDisplayValue(player);
</script>

{#if column.needsShow()}
  <td
    class="{$storedDisplayPref.isBorderVisible
      ? 'border border-gray-500'
      : ''} px-1 py-0.5"
    style="background-color: {column.getBgColorCode(player)}"
  >
    <div class="flex flex-row">
      <div>{param.damage}</div>
      {#if param.shipInfo}
        <div>
          <span>, </span>
          <ExternalLink url={param.shipInfo.url}>
            {param.shipInfo.name}
          </ExternalLink>
        </div>
      {/if}
    </div>
  </td>
{/if}
