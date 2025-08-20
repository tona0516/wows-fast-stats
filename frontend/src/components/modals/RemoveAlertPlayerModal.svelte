<script lang="ts">
  import { showToast, storedDeleteAlertPlayer } from "@libs/stores";
  import { RemoveAlertPlayer } from "@wails/go/main/App";
  import ModalCommon from "./ModalCommon.svelte";
  import { ModalManager } from "@libs/ModalManager";

  async function execute() {
    if (!$storedDeleteAlertPlayer) {
      ModalManager.instance.closeForDelete();
      return;
    }

    try {
      await RemoveAlertPlayer($storedDeleteAlertPlayer.account_id);
      showToast("削除しました");
    } catch (error) {
      showToast("削除に失敗しました");
    } finally {
      ModalManager.instance.closeForDelete();
    }
  }
</script>

{#if $storedDeleteAlertPlayer}
  <ModalCommon zValue={51} close={ModalManager.instance.closeForDelete}>
    <h3 class="font-bold text-lg mb-4">本当に削除しますか？</h3>
    <div class="mb-4 text-center">
      <span class="font-bold">{$storedDeleteAlertPlayer.name}</span>
      <span class="ml-2 text-xs text-gray-500"
        >(ID: {$storedDeleteAlertPlayer.account_id})</span
      >
    </div>
    <form class="modal-action">
      <button type="button" class="btn btn-error" on:click={execute}
        >削除</button
      >
    </form>
  </ModalCommon>
{/if}
