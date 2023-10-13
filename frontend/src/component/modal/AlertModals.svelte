<script lang="ts">
  import RemoveAlertPlayerModal from "src/component/modal/RemoveAlertPlayerModal.svelte";
  import EditAlertPlayerModal from "src/component/modal/EditAlertPlayerModal.svelte";
  import { createEventDispatcher } from "svelte";
  import AddAlertPlayerModal from "src/component/modal/AddAlertPlayerModal.svelte";
  import type { domain } from "wailsjs/go/models";

  const MAX_MEMO_LENGTH = 100;
  const EMPTY: domain.AlertPlayer = {
    account_id: 0,
    name: "",
    pattern: "bi-check-circle-fill",
    message: "",
  } as const;
  const dispatch = createEventDispatcher();
  let addModal: AddAlertPlayerModal;
  let editModal: EditAlertPlayerModal;
  let removeModal: RemoveAlertPlayerModal;

  export const showAdd = () => addModal.show();
  export const showEdit = (target: domain.AlertPlayer) =>
    editModal.show(target);
  export const showRemove = (target: domain.AlertPlayer) =>
    removeModal.show(target);
</script>

<AddAlertPlayerModal
  bind:this={addModal}
  defaultAlertPlayer={EMPTY}
  maxMemoLength={MAX_MEMO_LENGTH}
  on:Success={() => dispatch("AlertPlayerUpdated")}
  on:Failure
/>

<EditAlertPlayerModal
  bind:this={editModal}
  defaultAlertPlayer={EMPTY}
  maxMemoLength={MAX_MEMO_LENGTH}
  on:Success={() => dispatch("AlertPlayerUpdated")}
  on:Failure
/>

<RemoveAlertPlayerModal
  bind:this={removeModal}
  defaultAlertPlayer={EMPTY}
  on:Success={() => dispatch("AlertPlayerUpdated")}
  on:Failure
/>
