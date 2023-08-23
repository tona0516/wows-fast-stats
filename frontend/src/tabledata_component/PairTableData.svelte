<script lang="ts">
  import type { ComponentOption } from "src/ComponentList";
  import type { StatsCategory } from "src/enums";
  import { values } from "src/util";
  import type { domain } from "wailsjs/go/models";

  export let player: domain.Player;
  export let statsPattern: string;
  export let statsCatetory: StatsCategory;
  export let columnKey: string;
  export let option: ComponentOption;
  export let customDigit: domain.CustomDigit;

  $: value1 = values(player, statsPattern, statsCatetory, option.key1);
  $: value2 = values(player, statsPattern, statsCatetory, option.key2);
  $: digit = customDigit[columnKey];
</script>

<td class="td-multiple">
  {value1.toFixed(digit)}{option.unit} | {value2.toFixed(digit)}{option.unit}
</td>
