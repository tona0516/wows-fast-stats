<script lang="ts">
  import type { MaxDamage } from "src/lib/column/model/MaxDamage";
  import type { data } from "wailsjs/go/models";

  export let column: MaxDamage;
  export let player: data.Player;

  $: param = column.displayValue(player);
</script>

<div class="text-center tooltip" data-tip={param.shipInfo?.name}>
  {param.damage}
</div>
