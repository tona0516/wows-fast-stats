<script lang="ts">
  import { CONFIRM_MODAL_ID } from "src/lib/types";
  import { createEventDispatcher } from "svelte";
  import ModalCommon from "./ModalCommon.svelte";

  export let message: string;

  const dispatch = createEventDispatcher();
</script>

<ModalCommon id={CONFIRM_MODAL_ID}>
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
</ModalCommon>
