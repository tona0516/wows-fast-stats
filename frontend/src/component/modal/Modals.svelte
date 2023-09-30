<script lang="ts">
  import RemoveAlertPlayerModal from "src/component/modal/RemoveAlertPlayerModal.svelte";
  import EditAlertPlayerModal from "src/component/modal/EditAlertPlayerModal.svelte";
  import { createEventDispatcher } from "svelte";
  import AddAlertPlayerModal from "./AddAlertPlayerModal.svelte";
  import type { domain } from "wailsjs/go/models";

  const dispatch = createEventDispatcher();
  let addAlertPlayerModal: AddAlertPlayerModal;
  let editAlertPlayerModal: EditAlertPlayerModal;
  let removeAlertPlayerModal: RemoveAlertPlayerModal;

  export const showAddAlertPlayer = () => {
    addAlertPlayerModal.show();
  };

  export const showEditAlertPlayer = (target: domain.AlertPlayer) => {
    editAlertPlayerModal.show(target);
  };

  export const showRemoveAlertPlayer = (target: domain.AlertPlayer) => {
    removeAlertPlayerModal.show(target);
  };
</script>

<AddAlertPlayerModal
  bind:this={addAlertPlayerModal}
  on:Success={() => dispatch("AlertPlayerUpdated")}
  on:Failure
/>

<EditAlertPlayerModal
  bind:this={editAlertPlayerModal}
  on:Success={() => dispatch("AlertPlayerUpdated")}
  on:Failure
/>

<RemoveAlertPlayerModal
  bind:this={removeAlertPlayerModal}
  on:Success={() => dispatch("AlertPlayerUpdated")}
  on:Failure
/>
