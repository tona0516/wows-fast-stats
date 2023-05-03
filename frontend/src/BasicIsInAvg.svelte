<script lang="ts">
import type { vo } from "wailsjs/go/models";
import {
  AddExcludePlayerID,
  RemoveExcludePlayerID,
} from "../wailsjs/go/main/App.js";
import { createEventDispatcher } from "svelte";

export let player: vo.Player;
export let excludePlayerIDs: number[];

$: isChecked = !excludePlayerIDs.includes(player.player_info.id);

const dispatch = createEventDispatcher();

function onCheck(event) {
  if (event.target.checked) {
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

<td>
  {#if player.player_info.id !== 0}
    <div class="form-check m-0 centerize">
      <input
        class="form-check-input"
        type="checkbox"
        on:click="{onCheck}"
        checked="{isChecked}"
      />
    </div>
  {/if}
</td>
