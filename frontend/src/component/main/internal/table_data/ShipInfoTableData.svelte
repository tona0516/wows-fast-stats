<script lang="ts">
  import ExternalLink from "src/component/common/ExternalLink.svelte";
  import type { data } from "wailsjs/go/models";
  import { CssClass } from "src/lib/CssClass";
  import { NumbersURL } from "src/lib/NumbersURL";
  import type { ShipInfo } from "src/lib/column/model/ShipInfo";

  export let column: ShipInfo;
  export let player: data.Player;
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

<td class={CssClass.TD_STR}>
  <ExternalLink url={NumbersURL.ship(player.ship_info.id)}
    ><div class="uk-text-truncate">
      {column.displayValue(player)}
    </div></ExternalLink
  >
</td>

<style>
  :global(.ship-icon) {
    width: 1.25em;
    display: flex;
    align-items: center;
  }
</style>
