<script lang="ts">
import type { vo } from "../../wailsjs/go/models";
import type { ComponentOption } from "../ComponentList";
import type { StatsCategory } from "../enums";
import { values } from "../util";

export let player: vo.Player;
export let statsPattern: string;
export let statsCatetory: StatsCategory;
export let key: string;
export let option: ComponentOption;
export let customDigit: vo.CustomDigit;

$: value1 = values(player, statsPattern, statsCatetory, option.key1);
$: value2 = values(player, statsPattern, statsCatetory, option.key2);
$: digit = customDigit[key];
</script>

<td class="td-multiple">
  {value1.toFixed(digit)}{option.unit} | {value2.toFixed(digit)}{option.unit}
</td>
