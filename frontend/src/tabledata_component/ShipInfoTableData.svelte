<script lang="ts">
  import type { ShipInfo } from "src/lib/column/ShipInfo";
  import ExternalLink from "src/other_component/ExternalLink.svelte";
  import type { domain } from "wailsjs/go/models";

  export let column: ShipInfo;
  export let player: domain.Player;
</script>

<td class="td-icon">
  <img alt="" src={column.nationIconPath(player)} class="nation-icon" />
</td>

<td class="td-icon" style="background-color: {column.bgColorCode(player)}">
  <img alt="" src={column.shipTypeIconPath(player)} class="ship-icon" />
</td>

<td class="td-string omit">
  <ExternalLink
    url={column.shipURL(player)}
    text={column.displayValue(player)}
  />
</td>
