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

<div class="w-48 flex place-items-center">
  {#if isNPC}
    <div class="truncate">
      {column.playerName(player)}
    </div>
  {:else}
    <button
      class="btn btn-xs mr-1"
      on:click={() => PlayerDetailModal.open(player.player_info)}
    >
      <i class="bi bi-box-arrow-in-up-right"></i>
    </button>
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
      <span style="color: {column.textColorCode(player)}">
        {column.playerName(player)}
      </span>
    </div>
  {/if}
</div>
