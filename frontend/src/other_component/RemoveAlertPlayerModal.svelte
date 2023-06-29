<script lang="ts">
import { createEventDispatcher } from "svelte";
import { Modal, ModalBody, ModalFooter, Button } from "sveltestrap";
import type { vo } from "../../wailsjs/go/models";
import { RemoveAlertPlayer } from "../../wailsjs/go/main/App";

export const toggle = () => (open = !open);
export const setTarget = (p: vo.AlertPlayer) => (target = p);

const dispatch = createEventDispatcher();

let open = false;
let target: vo.AlertPlayer;

async function remove(accountID: number) {
  try {
    await RemoveAlertPlayer(accountID);
    dispatch("Success");
  } catch (error) {
    dispatch("Failure", { message: error });
    return;
  } finally {
    toggle();
  }
}
</script>

<Modal isOpen="{open}" toggle="{toggle}">
  <ModalBody class="modal-color">
    「{target.name}」を削除しますか？
  </ModalBody>
  <ModalFooter class="modal-color">
    <Button size="sm" color="secondary" on:click="{toggle}">キャンセル</Button>
    <Button
      size="sm"
      color="danger"
      on:click="{() => remove(target.account_id)}">削除</Button
    >
  </ModalFooter>
</Modal>
