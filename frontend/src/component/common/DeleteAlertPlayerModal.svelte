<script lang="ts">
  import { RemoveAlertPlayer, ShowMessageDialog } from "wailsjs/go/main/App";
  import {
    closeModal,
    storedAlertPlayerForm,
    storedIsShowDeleteAlertPlayerModal,
  } from "src/stores";

  async function remove() {
    if ($storedAlertPlayerForm.account_id !== 0) {
      try {
        await RemoveAlertPlayer($storedAlertPlayerForm.account_id);
      } catch (error) {
        ShowMessageDialog("削除に失敗しました");
      }
    }

    closeModal();
  }
</script>

{#if $storedIsShowDeleteAlertPlayerModal}
  <dialog class="modal modal-open">
    <form method="dialog" class="modal-box" on:submit|preventDefault={remove}>
      <h3 class="font-bold text-lg mb-4">本当に削除しますか？</h3>
      <div class="mb-4 text-center">
        <span class="font-bold">{$storedAlertPlayerForm.name}</span>
        <span class="ml-2 text-xs text-gray-500"
          >(ID: {$storedAlertPlayerForm.account_id})</span
        >
      </div>
      <div class="modal-action">
        <button type="button" class="btn" on:click={() => closeModal()}
          >キャンセル</button
        >
        <button type="submit" class="btn btn-error">削除</button>
      </div>
    </form>
  </dialog>
{/if}
