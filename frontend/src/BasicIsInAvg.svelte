<script lang="ts">
import type { vo } from "wailsjs/go/models";
import {
  AddExcludePlayerID,
  RemoveExcludePlayerID,
} from "../wailsjs/go/main/App.js";
import { createEventDispatcher } from "svelte";
import type { DisplayPattern } from "./DisplayPattern.js";
import { storedExcludePlayerIDs } from "./stores.js";
import { get } from "svelte/store";

export let player: vo.Player;
export let displayPattern: DisplayPattern;

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

{#if displayPattern === "noshipstats" || displayPattern === "full" || displayPattern === "nopr"}
  <td class="td-checkbox">
    <input
      class="form-check-input"
      type="checkbox"
      on:click="{onCheck}"
      checked="{isChecked}"
    />
  </td>
{:else}
  <td class="td-checkbox"></td>
{/if}
