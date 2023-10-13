<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import UkModal from "src/component/common/uikit/UkModal.svelte";
  import { ModalElementID } from "./ModalElementID";

  export let message: string;

  const dispatch = createEventDispatcher();
</script>

<UkModal id={ModalElementID.CONFIRM}>
  <div slot="body">
    <p>{message}</p>
  </div>

  <div slot="footer">
    <button
      class="uk-button uk-button-primary uk-modal-close"
      type="button"
      on:click={() => dispatch("Confirmed")}>はい</button
    >
  </div>
</UkModal>
