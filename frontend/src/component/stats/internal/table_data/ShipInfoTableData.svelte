<script lang="ts">
  import type { ShipInfo } from "src/lib/column/model/ShipInfo";
  import { ShipDetailModal } from "src/stores";
  import type { data } from "wailsjs/go/models";

  export let column: ShipInfo;
  export let player: data.Player;
</script>

<div class="w-40 flex place-items-center">
  <button
    class="btn btn-xs mr-1 bi bi-box-arrow-in-up-right"
    on:click={() => {
      ShipDetailModal.open(player);
    }}
  />
  <img class="w-icon" src={column.nationIconPath(player)} alt="" />
  <img class="w-icon" src={column.shipTypeIconPath(player)} alt="" />
  <div class="truncate">
    {column.displayValue(player)}
  </div>
</div>

<style>
  .w-icon {
    width: 1.25rem;
    margin-left: 1px;
    margin-right: 1px;
  }
</style>
