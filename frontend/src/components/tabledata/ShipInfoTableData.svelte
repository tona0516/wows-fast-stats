<script lang="ts">
  import type { ShipInfoColumn } from "@libs/columns/ShipInfoColumn";
  import { ModalManager } from "@libs/ModalManager";
  import { storedZoomRate } from "@libs/stores";
  import type { data } from "@wails/go/models";

  export let column: ShipInfoColumn;
  export let player: data.Player;

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
        style="width: {(1.25 * $storedZoomRate) / 100}rem"
        src={nationIconPath}
        alt=""
      />
    {/if}
    <img
      class="w-icon"
      style="width: {(1.25 * $storedZoomRate) / 100}rem"
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
