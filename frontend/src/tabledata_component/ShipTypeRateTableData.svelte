<script lang="ts">
  import type { domain } from "../../wailsjs/go/models";
  import { Const } from "../Const";
  import type { StackedBarGraphParam } from "../other_component/stacked_bar/StackedBarGraphParam";
  import StackedBarGraph from "../other_component/stacked_bar/StackedBarGraph.svelte";
  import type { StatsCategory } from "../enums";
  import { values } from "../util";

  export let player: domain.Player;
  export let statsPattern: string;
  export let statsCatetory: StatsCategory;
  export let columnKey: string;
  export let customColor: domain.CustomColor;
  export let customDigit: domain.CustomDigit;

  function toParam(
    player: domain.Player,
    shipTypeGroup: { [key: string]: number },
    customColor: domain.CustomColor,
    customDigit: domain.CustomDigit
  ): StackedBarGraphParam {
    const digit = customDigit[columnKey];

    if (shipTypeGroup === undefined) {
      return { digit: digit, items: [] };
    }

    const items = Object.keys(shipTypeGroup).map((key) => {
      const label = Const.SHIP_TYPE_LABELS[key] ?? "";
      const color =
        key === player.ship_info.type
          ? customColor.ship_type.own[key]
          : customColor.ship_type.other[key];
      const value = shipTypeGroup[key];

      return {
        label: label,
        color: color,
        value: value,
      };
    });

    return {
      digit: digit,
      items: items,
    };
  }

  $: shipTypeGroup = values(player, statsPattern, statsCatetory, columnKey) as {
    [key: string]: number;
  };

  $: param = toParam(player, shipTypeGroup, customColor, customDigit);
</script>

<td class="td-graph">
  <StackedBarGraph {param} />
</td>
