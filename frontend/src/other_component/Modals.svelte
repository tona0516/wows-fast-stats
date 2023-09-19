<script lang="ts">
  import RemoveAlertPlayerModal from "src/other_component/RemoveAlertPlayerModal.svelte";
  import UpdateAlertPlayerModal from "src/other_component/UpdateAlertPlayerModal.svelte";
  import { createEventDispatcher } from "svelte";
  import type { domain } from "wailsjs/go/models";

  let updateAlertPlayerModal: UpdateAlertPlayerModal;
  let removeAlertPlayerModal: RemoveAlertPlayerModal;

  const dispatch = createEventDispatcher();

  export const addAlertPlayer = () => {
    updateAlertPlayerModal.openAsAdd();
  };

  export const editAlertPlayer = (target: domain.AlertPlayer) => {
    updateAlertPlayerModal.openAsEdit(target);
  };

  export const removeAlertPlayer = (target: domain.AlertPlayer) => {
    removeAlertPlayerModal.open(target);
  };
</script>

<UpdateAlertPlayerModal
  bind:this={updateAlertPlayerModal}
  on:Success={() => dispatch("AlertPlayerUpdated")}
  on:Failure
/>

<RemoveAlertPlayerModal
  bind:this={removeAlertPlayerModal}
  on:Success={() => dispatch("AlertPlayerUpdated")}
  on:Failure
/>
