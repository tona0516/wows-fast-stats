<script lang="ts">
import type { vo } from "wailsjs/go/models";
import {
  AddExcludePlayerID,
  RemoveExcludePlayerID,
} from "../wailsjs/go/main/App.js";
import { createEventDispatcher } from "svelte";
import type { DisplayPattern } from "./DisplayPattern.js";

export let player: vo.Player;
export let excludePlayerIDs: number[];
export let displayPattern: DisplayPattern;

$: isChecked = !excludePlayerIDs.includes(player.player_info.id);

const dispatch = createEventDispatcher();

async function onCheck(e: any) {
  if (e.target.checked) {
    await RemoveExcludePlayerID(player.player_info.id);
    dispatch("onCheck", null);
  } else {
    await AddExcludePlayerID(player.player_info.id);
    dispatch("onCheck", null);
  }
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
