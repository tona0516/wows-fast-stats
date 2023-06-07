<script lang="ts">
import type { vo } from "../../wailsjs/go/models";
import type { ComponentOption } from "../ComponentList";
import { Const } from "../Const";
import type { StatsCategory } from "../enums";
import { colors, values } from "../util";

export let player: vo.Player;
export let statsPattern: string;
export let statsCatetory: StatsCategory;
export let key: string;
export let option: ComponentOption;

const digit = Const.DIGITS[key];

$: color = colors(key, value, player, statsCatetory);
$: value = values(player, statsPattern, statsCatetory, key);
</script>

<td class="td-number" style="color: {color}">
  {#if value !== -1}
    {value.toFixed(digit)}{option.unit}
  {/if}
</td>
