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

function onCheck(e: any) {
  if (e.target.checked) {
    RemoveExcludePlayerID(player.player_info.id).then(() => {
      dispatch("onCheck", null);
    });
  } else {
    AddExcludePlayerID(player.player_info.id).then(() => {
      dispatch("onCheck", null);
    });
  }
}
</script>

<td class="p-0">
  {#if displayPattern === "noshipstats" || displayPattern === "full" || displayPattern === "nopr"}
    <div class="form-check centerize">
      <input
        class="form-check-input"
        type="checkbox"
        on:click="{onCheck}"
        checked="{isChecked}"
      />
    </div>
  {/if}
</td>
