<script lang="ts">
import { createEventDispatcher } from "svelte";
import { Modal, ModalBody, ModalFooter, Button } from "sveltestrap";

export const toggle = () => (open = !open);
export let message: string;

const dispatch = createEventDispatcher();

let open = false;
</script>

<Modal isOpen="{open}" toggle="{toggle}">
  <ModalBody class="modal-color">
    {message}
  </ModalBody>
  <ModalFooter class="modal-color">
    <Button size="sm" color="secondary" on:click="{toggle}">キャンセル</Button>
    <Button
      size="sm"
      color="primary"
      on:click="{() => {
        toggle();
        dispatch('onConfirmed', null);
      }}">はい</Button
    >
  </ModalFooter>
</Modal>
