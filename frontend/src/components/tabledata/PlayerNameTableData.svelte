<script lang="ts">
  import type { PlayerNameColumn } from "@libs/columns/PlayerNameColumn";
  import { ModalManager } from "@libs/ModalManager";
  import { storedDisplayPref } from "@libs/stores";
  import type { core } from "@wails/go/models";

  export let column: PlayerNameColumn;
  export let player: core.Player;

  $: clanTag = column.getClanTag(player);
  $: nationFlagClass = column.getNationFlagClass(player);
  $: isNPC = column.isNPC(player);
</script>

<td
  class="{$storedDisplayPref.isBorderVisible
    ? 'border border-gray-500'
    : ''} px-1 py-0.5"
>
  <div class="w-48 flex place-items-center">
    {#if isNPC}
      <div class="truncate">
        {column.getPlayerName(player)}
      </div>
    {:else}
      <button
        class="btn btn-xs bi bi-info-square px-1 py-0.5 mr-1"
        on:click={() => ModalManager.instance.openForPlayerDetail(player)}
      />

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
