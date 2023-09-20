<script lang="ts">
  import type { ShipInfo } from "src/lib/column/ShipInfo";
  import ExternalLink from "src/other_component/ExternalLink.svelte";
  import type { domain } from "wailsjs/go/models";

  export let column: ShipInfo;
  export let player: domain.Player;
</script>

<td>
  <!-- svelte-ignore a11y-missing-attribute -->
  <img
    src={column.nationIconPath(player)}
    class="nation-icon uk-preserve-width"
  />
</td>

<td style="background-color: {column.bgColorCode(player)}">
  <!-- svelte-ignore a11y-missing-attribute -->
  <img
    src={column.shipTypeIconPath(player)}
    class="ship-icon uk-preserve-width"
  />
</td>

<td class="td-string omit">
  <ExternalLink url={column.shipURL(player)}
    >{column.displayValue(player)}</ExternalLink
  >
</td>
