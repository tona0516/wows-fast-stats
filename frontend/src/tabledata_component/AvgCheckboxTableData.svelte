<script lang="ts">
import { createEventDispatcher } from "svelte";
import { get } from "svelte/store";
import {
  RemoveExcludePlayerID,
  AddExcludePlayerID,
} from "../../wailsjs/go/main/App";
import type { vo } from "../../wailsjs/go/models";
import { storedExcludePlayerIDs } from "../stores";

export let player: vo.Player;

let excludePlayerIDs = get(storedExcludePlayerIDs);
storedExcludePlayerIDs.subscribe((it) => (excludePlayerIDs = it));

$: isChecked = !excludePlayerIDs.includes(player.player_info.id);

const dispatch = createEventDispatcher();

async function onCheck(e: any) {
  const accountID = player.player_info.id;
  if (e.target.checked) {
    await RemoveExcludePlayerID(accountID);
  } else {
    await AddExcludePlayerID(accountID);
  }
  dispatch("CheckPlayer");
}
</script>

<td class="td-checkbox">
  {#if player.player_info.id !== 0}
    <input
      class="form-check-input"
      type="checkbox"
      on:click="{onCheck}"
      checked="{isChecked}"
    />
  {/if}
</td>
