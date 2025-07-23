<script lang="ts">
  import {
    DeleteAlertPlayerModal,
    showToast,
    storedDeleteAlertPlayer,
  } from "src/stores";
  import { RemoveAlertPlayer } from "wailsjs/go/main/App";
  import ModalCommon from "./ModalCommon.svelte";

  async function remove() {
    if (!$storedDeleteAlertPlayer) {
      DeleteAlertPlayerModal.close();
      return;
    }

    try {
      await RemoveAlertPlayer($storedDeleteAlertPlayer.account_id);
      showToast("削除しました");
    } catch (error) {
      showToast("削除に失敗しました");
    } finally {
      DeleteAlertPlayerModal.close();
    }
  }
</script>

{#if $storedDeleteAlertPlayer}
  <ModalCommon zValue={51} close={DeleteAlertPlayerModal.close}>
    <h3 class="font-bold text-lg mb-4">本当に削除しますか？</h3>
    <div class="mb-4 text-center">
      <span class="font-bold">{$storedDeleteAlertPlayer.name}</span>
      <span class="ml-2 text-xs text-gray-500"
        >(ID: {$storedDeleteAlertPlayer.account_id})</span
      >
    </div>
    <form class="modal-action">
      <button type="button" class="btn btn-error" on:click={remove}>削除</button
      >
    </form>
  </ModalCommon>
{/if}
