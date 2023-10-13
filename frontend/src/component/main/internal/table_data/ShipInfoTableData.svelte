<script lang="ts">
  import type { ShipInfo } from "src/lib/column/ShipInfo";
  import ExternalLink from "src/component/common/ExternalLink.svelte";
  import type { domain } from "wailsjs/go/models";
  import { CssClass } from "src/lib/types";
  import { NumbersURL } from "src/lib/NumbersURL";

  export let column: ShipInfo;
  export let player: domain.Player;
</script>

<td>
  <img
    alt=""
    src={column.nationIconPath(player)}
    class="nation-icon uk-preserve-width"
  />
</td>

<td style="background-color: {column.bgColorCode(player)}">
  <img
    alt=""
    src={column.shipTypeIconPath(player)}
    class="ship-icon uk-preserve-width"
  />
</td>

<td class="{CssClass.TD_STR} {CssClass.OMIT}">
  <ExternalLink
    url={NumbersURL.ship(player.ship_info.id, player.ship_info.name)}
    >{column.displayValue(player)}</ExternalLink
  >
</td>

<style>
  :global(.nation-icon) {
    width: 1.5em;
    display: flex;
    align-items: center;
  }

  :global(.ship-icon) {
    width: 1.25em;
    display: flex;
    align-items: center;
  }
</style>
