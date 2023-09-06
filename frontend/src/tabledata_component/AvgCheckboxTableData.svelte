<script lang="ts">
  import type { IsInAvg } from "src/lib/column/IsInAvg";
  import { storedExcludePlayerIDs } from "src/stores";
  import { createEventDispatcher } from "svelte";
  import {
    RemoveExcludePlayerID,
    AddExcludePlayerID,
  } from "wailsjs/go/main/App";
  import type { domain } from "wailsjs/go/models";

  export let column: IsInAvg;
  export let player: domain.Player;

  $: accountID = player.player_info.id;
  $: isChecked = !$storedExcludePlayerIDs.includes(accountID);

  const dispatch = createEventDispatcher();

  const onCheck = async (e: any) => {
    if (e.target.checked) {
      await RemoveExcludePlayerID(accountID);
    } else {
      await AddExcludePlayerID(accountID);
    }
    dispatch("CheckPlayer");
  };
</script>

<td class="td-checkbox">
  {#if column.displayValue(player)}
    <input
      class="form-check-input"
      type="checkbox"
      on:click={onCheck}
      checked={isChecked}
    />
  {/if}
</td>
