<script lang="ts">
  import ExternalLink from "src/component/common/ExternalLink.svelte";
  import type { domain } from "wailsjs/go/models";
  import { CssClass } from "src/lib/CssClass";
  import type { MaxDamage } from "src/lib/column/model/MaxDamage";

  export let column: MaxDamage;
  export let player: domain.Player;
</script>

<td class={CssClass.TD_NUM}>
  {column.damage(player)}
</td>

{#if column.innerColumnNumber > 1}
  {@const [url, text] = column.shipInfo(player)}
  <td class="{CssClass.TD_STR} {CssClass.OMIT}">
    <ExternalLink {url}>{text}</ExternalLink>
  </td>
{/if}
