<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { Modal, ModalBody, ModalFooter, Button } from "sveltestrap";

  export const toggle = () => (open = !open);
  export let message: string;

  let open = false;

  const dispatch = createEventDispatcher();
</script>

<Modal isOpen={open} {toggle}>
  <ModalBody class="modal-color">
    {message}
  </ModalBody>
  <ModalFooter class="modal-color">
    <Button size="sm" color="secondary" on:click={toggle}>キャンセル</Button>
    <Button
      size="sm"
      color="primary"
      on:click={() => {
        toggle();
        dispatch("onConfirmed", null);
      }}>はい</Button
    >
  </ModalFooter>
</Modal>
