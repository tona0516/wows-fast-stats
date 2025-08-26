<script lang="ts">
  import type { ShipInfoColumn } from "@libs/columns/ShipInfoColumn";
  import { ModalManager } from "@libs/ModalManager";
  import type { data } from "@wails/go/models";

  export let column: ShipInfoColumn;
  export let player: data.Player;
</script>

<td class="p-1" style="background-color: {column.getBgColorCode(player)?.raw}">
  <div class="w-48 flex place-items-center">
    <button
      class="btn btn-xs bi bi-info-square p-1 mr-1"
      on:click={() => ModalManager.instance.openForShipDetail(player)}
    />
    <img class="w-icon" src={column.getNationIconPath(player)} alt="" />
    <img class="w-icon" src={column.getShipIconPath(player)} alt="" />
    <div class="truncate">
      {column.getDisplayValue(player)}
    </div>
  </div>
</td>

<style>
  .w-icon {
    width: 1.25rem;
    margin-left: 1px;
    margin-right: 1px;
  }
</style>
