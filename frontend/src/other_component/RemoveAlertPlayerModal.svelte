<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { Modal, ModalBody, ModalFooter, Button } from "sveltestrap";
  import { RemoveAlertPlayer } from "wailsjs/go/main/App";
  import type { domain } from "wailsjs/go/models";

  export const open = (p: domain.AlertPlayer) => {
    target = p;
    isOpen = true;
  };

  let isOpen = false;
  let target: domain.AlertPlayer;

  const dispatch = createEventDispatcher();

  const remove = async (accountID: number) => {
    try {
      await RemoveAlertPlayer(accountID);
      dispatch("Success");
    } catch (error) {
      dispatch("Failure", { message: error });
      return;
    } finally {
      isOpen = false;
    }
  };
</script>

<Modal {isOpen}>
  <ModalBody class="modal-color">
    「{target.name}」を削除しますか？
  </ModalBody>
  <ModalFooter class="modal-color">
    <Button size="sm" color="secondary" on:click={() => (isOpen = false)}
      >キャンセル</Button
    >
    <Button size="sm" color="danger" on:click={() => remove(target.account_id)}
      >削除</Button
    >
  </ModalFooter>
</Modal>
