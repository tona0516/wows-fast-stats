<script lang="ts">
  import { CONFIRM_MODAL_ID } from "src/lib/types";
  import { createEventDispatcher } from "svelte";
  import UkModal from "src/component/common/uikit/UkModal.svelte";

  export let message: string;

  const dispatch = createEventDispatcher();
</script>

<UkModal id={CONFIRM_MODAL_ID}>
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
