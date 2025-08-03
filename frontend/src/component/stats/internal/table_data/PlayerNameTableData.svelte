<script lang="ts">
  import type { PlayerName } from "src/lib/column/model/PlayerName";
  import { PlayerDetailModal, storedAlertPlayers } from "src/stores";
  import type { data } from "wailsjs/go/models";

  export let column: PlayerName;
  export let player: data.Player;

  $: alertPlayer = $storedAlertPlayers.find(
    (it) => it.account_id === player.player_info.id,
  );
  $: clanTag = column.clanTag(player);
  $: nationFlagClass = column.clanFlagIconClass(player);
  $: isNPC = column.isNPC(player);
</script>

<td
  class="p-1"
  style="background-color: {column.getBackgroundColorCode(player) ?? ''}"
>
  <div class="w-48 flex place-items-center">
    {#if isNPC}
      <div class="truncate">
        {column.playerName(player)}
      </div>
    {:else}
      <button
        class="btn btn-xs mr-1 bi bi-info-square"
        on:click={() => PlayerDetailModal.open(player)}
      />
      <div class="truncate">
        {#if alertPlayer}
          <span class="bi {alertPlayer.pattern}"></span>
        {/if}
        {#if clanTag}
          {#if nationFlagClass}
            <span class={nationFlagClass}></span>
          {/if}
          <span style="color: {column.clanColorCode(player)}">
            {clanTag}
          </span>
        {/if}
        {column.playerName(player)}
      </div>
    {/if}
  </div>
</td>
