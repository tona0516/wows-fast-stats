<script lang="ts">
  import type { ComponentOption } from "src/ComponentList";
  import type { StatsCategory } from "src/enums";
  import { colors, values } from "src/util";
  import type { domain } from "wailsjs/go/models";

  export let player: domain.Player;
  export let statsPattern: string;
  export let statsCatetory: StatsCategory;
  export let columnKey: string;
  export let option: ComponentOption;
  export let customColor: domain.CustomColor;
  export let customDigit: domain.CustomDigit;

  $: color = colors(columnKey, value, player, statsCatetory, customColor.skill);
  $: value = values(player, statsPattern, statsCatetory, columnKey);
  $: digit = customDigit[columnKey];
</script>

<td class="td-number" style="color: {color}">
  {#if value !== -1}
    {value.toFixed(digit)}{option.unit}
  {/if}
</td>
