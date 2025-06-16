<script lang="ts">
import ExternalLink from "src/component/common/ExternalLink.svelte";
import { NumbersURL } from "src/lib/NumbersURL";
import type { ShipInfo } from "src/lib/column/model/ShipInfo";
import type { data } from "wailsjs/go/models";

export let column: ShipInfo;
export let player: data.Player;
</script>

<div>
  <ExternalLink url={NumbersURL.ship(player.ship_info.id)}>
    <div
      class="flex place-items-center tooltip"
      data-tip={column.tooltipText(player)}
    >
      <img class="w-6" src={column.nationIconPath(player)} alt="" />
      <img class="w-6" src={column.shipTypeIconPath(player)} alt="" />
      <div class="truncate">
        {column.displayValue(player)}
      </div>
    </div>
  </ExternalLink>
</div>
