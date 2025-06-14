<script lang="ts">
import ExternalLink from "src/component/common/ExternalLink.svelte";
import type { MaxDamage } from "src/lib/column/model/MaxDamage";
import type { data } from "wailsjs/go/models";

export let column: MaxDamage;
export let player: data.Player;

$: param = column.displayValue(player);
</script>

<td>
  {param.damage}
</td>

{#if param.shipInfo}
  <td>
    <ExternalLink url={param.shipInfo.url}>{param.shipInfo.name}</ExternalLink>
  </td>
{/if}
