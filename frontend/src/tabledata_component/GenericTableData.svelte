<script lang="ts">
import type { vo } from "../../wailsjs/go/models";
import type { ComponentOption } from "../ComponentList";
import type { StatsCategory } from "../enums";
import { colors, values } from "../util";

export let player: vo.Player;
export let statsPattern: string;
export let statsCatetory: StatsCategory;
export let columnKey: string;
export let option: ComponentOption;
export let customColor: vo.CustomColor;
export let customDigit: vo.CustomDigit;

$: color = colors(columnKey, value, player, statsCatetory, customColor.skill);
$: value = values(player, statsPattern, statsCatetory, columnKey);
$: digit = customDigit[columnKey];
</script>

<td class="td-number" style="color: {color}">
  {#if value !== -1}
    {value.toFixed(digit)}{option.unit}
  {/if}
</td>
