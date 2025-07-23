<script lang="ts">
  import {
    DeleteAlertPlayerModal,
    showToast,
    storedDeleteAlertPlayer,
  } from "src/stores";
  import { RemoveAlertPlayer } from "wailsjs/go/main/App";

  async function remove() {
    if (!$storedDeleteAlertPlayer) {
      return;
    }

    try {
      await RemoveAlertPlayer($storedDeleteAlertPlayer.account_id);
    } catch (error) {
      showToast("削除に失敗しました");
    } finally {
      DeleteAlertPlayerModal.close();
    }
  }
</script>

{#if $storedDeleteAlertPlayer}
  <dialog class="modal modal-open z-51">
    <div class="modal-box">
      <form method="dialog">
        <button
          class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2"
          on:click={DeleteAlertPlayerModal.close}
          ><i class="bi bi-x-lg"></i>
        </button>
      </form>
      <h3 class="font-bold text-lg mb-4">本当に削除しますか？</h3>
      <div class="mb-4 text-center">
        <span class="font-bold">{$storedDeleteAlertPlayer.name}</span>
        <span class="ml-2 text-xs text-gray-500"
          >(ID: {$storedDeleteAlertPlayer.account_id})</span
        >
      </div>
      <form class="modal-action">
        <button type="button" class="btn btn-error" on:click={remove}
          >削除</button
        >
      </form>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button type="button" on:click={DeleteAlertPlayerModal.close}
        >close</button
      >
    </form>
  </dialog>
{/if}
