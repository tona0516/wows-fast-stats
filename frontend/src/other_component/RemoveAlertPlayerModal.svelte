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
  <ModalBody style="background-color: #2d2c2c;">
    「{target.name}」を削除しますか？
  </ModalBody>
  <ModalFooter style="background-color: #2d2c2c;">
    <Button color="secondary" on:click="{toggle}">キャンセル</Button>
    <Button color="danger" on:click="{() => remove(target.account_id)}"
      >追加</Button
    >
  </ModalFooter>
</Modal>
