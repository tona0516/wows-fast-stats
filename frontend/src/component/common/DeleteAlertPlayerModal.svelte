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
    <form method="dialog" class="modal-box" on:submit|preventDefault={remove}>
      <h3 class="font-bold text-lg mb-4">本当に削除しますか？</h3>
      <div class="mb-4 text-center">
        <span class="font-bold">{$storedDeleteAlertPlayer.name}</span>
        <span class="ml-2 text-xs text-gray-500"
          >(ID: {$storedDeleteAlertPlayer.account_id})</span
        >
      </div>
      <div class="modal-action">
        <button
          type="button"
          class="btn"
          on:click={() => DeleteAlertPlayerModal.close()}>キャンセル</button
        >
        <button type="submit" class="btn btn-error">削除</button>
      </div>
    </form>
  </dialog>
{/if}
