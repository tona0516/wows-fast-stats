<script lang="ts">
  import type { PlayerNameColumn } from "@libs/columns/PlayerNameColumn";
  import { ModalManager } from "@libs/ModalManager";
  import { storedBlackList } from "@libs/stores";
  import type { data } from "@wails/go/models";

  export let column: PlayerNameColumn;
  export let player: data.Player;

  $: blackListItem = $storedBlackList.find(
    (it) => it.account_id === player.player_info.id,
  );
  $: clanTag = column.getClanTag(player);
  $: nationFlagClass = column.getNationFlagClass(player);
  $: isNPC = column.isNPC(player);
</script>

<td class="p-1">
  <div class="w-48 flex place-items-center">
    {#if isNPC}
      <div class="truncate">
        {column.getPlayerName(player)}
      </div>
    {:else}
      <button
        class="btn btn-xs bi bi-info-square p-1 mr-1"
        on:click={() => ModalManager.instance.openForPlayerDetail(player)}
      />
      {#if blackListItem}
        <span class="bi {blackListItem.pattern}"></span>
      {/if}

      {#if clanTag}
        {#if nationFlagClass}
          <span class={nationFlagClass}></span>
        {/if}
        <span style="color: {column.getClanColorCode(player)?.raw}">
          {clanTag}
        </span>
      {/if}
      <div
        class="truncate"
        style="color: {column.getTextColorCode(player)?.raw};"
      >
        {column.getPlayerName(player)}
      </div>
    {/if}
  </div>
</td>
