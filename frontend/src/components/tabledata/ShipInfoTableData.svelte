<script lang="ts">
  import type { WarshipColumn } from "@libs/columns/WarshipColumn";
  import { ModalManager } from "@libs/ModalManager";
    import { storedPref } from "@libs/stores";
  import type { core } from "@wails/go/models";

  export let column: WarshipColumn;
  export let player: core.Player;

  const nationIconPath = column.getNationIconPath(player);
</script>

<td class="p-1" style="background-color: {column.getBgColorCode(player)?.raw}">
  <div class="w-48 flex place-items-center">
    <button
      class="btn btn-xs bi bi-info-square p-1 mr-1"
      on:click={() => ModalManager.instance.openForShipDetail(player)}
    />
    {#if nationIconPath}
      <img
        class="w-icon"
        style="width: {(1.25 * $storedPref.zoomRate) / 100}rem"
        src={nationIconPath}
        alt=""
      />
    {/if}
    <img
      class="w-icon"
      style="width: {(1.25 * $storedPref.zoomRate) / 100}rem"
      src={column.getShipIconPath(player)}
      alt=""
    />
    <div class="truncate">
      {column.getDisplayValue(player)}
    </div>
  </div>
</td>

<style>
  .w-icon {
    margin-left: 1px;
    margin-right: 1px;
  }
</style>
