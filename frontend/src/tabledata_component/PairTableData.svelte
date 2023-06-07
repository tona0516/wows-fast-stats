<script lang="ts">
import type { vo } from "../../wailsjs/go/models";
import type { ComponentOption } from "../ComponentList";
import { Const } from "../Const";
import type { StatsCategory } from "../enums";
import { values } from "../util";

export let player: vo.Player;
export let statsPattern: string;
export let statsCatetory: StatsCategory;
export let key: string;
export let option: ComponentOption;

const digitKey1 = Const.DIGITS[key];
const digitKey2 = Const.DIGITS[key];

$: value1 = values(player, statsPattern, statsCatetory, option.key1);
$: value2 = values(player, statsPattern, statsCatetory, option.key2);
</script>

<td class="td-multiple">
  {value1.toFixed(digitKey1)}{option.unit} | {value2.toFixed(
    digitKey2
  )}{option.unit}
</td>
