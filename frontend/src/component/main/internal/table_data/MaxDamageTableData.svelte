<script lang="ts">
  import ExternalLink from "src/component/common/ExternalLink.svelte";
  import type { data } from "wailsjs/go/models";
  import { CssClass } from "src/lib/CssClass";
  import type { MaxDamage } from "src/lib/column/model/MaxDamage";

  export let column: MaxDamage;
  export let player: data.Player;

  $: param = column.displayValue(player);
</script>

<td class={CssClass.TD_NUM}>
  {param.damage}
</td>

{#if param.warship}
  <td class={CssClass.TD_STR}>
    <ExternalLink url={param.warship.url}
      ><div class="uk-text-truncate">{param.warship.name}</div></ExternalLink
    >
  </td>
{/if}
