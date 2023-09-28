<script lang="ts">
  import type { MaxDamage } from "src/lib/column/MaxDamage";
  import ExternalLink from "src/other_component/ExternalLink.svelte";
  import type { domain } from "wailsjs/go/models";

  export let column: MaxDamage;
  export let player: domain.Player;
</script>

<td class="td-number">
  {column.damage(player)}
</td>

{#if column.countInnerColumn() > 1}
  {@const [url, text] = column.shipInfo(player)}
  <td class="td-string omit">
    <ExternalLink {url}>{text}</ExternalLink>
  </td>
{/if}
