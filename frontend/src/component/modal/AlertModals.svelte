<script lang="ts">
  import RemoveAlertPlayerModal from "src/component/modal/RemoveAlertPlayerModal.svelte";
  import EditAlertPlayerModal from "src/component/modal/EditAlertPlayerModal.svelte";
  import AddAlertPlayerModal from "src/component/modal/AddAlertPlayerModal.svelte";
  import type { model } from "wailsjs/go/models";

  const MAX_MEMO_LENGTH = 100;
  const EMPTY: model.AlertPlayer = {
    account_id: 0,
    name: "",
    pattern: "bi-check-circle-fill",
    message: "",
  } as const;
  let addModal: AddAlertPlayerModal;
  let editModal: EditAlertPlayerModal;
  let removeModal: RemoveAlertPlayerModal;

  export const showAdd = () => addModal.show();
  export const showEdit = (target: model.AlertPlayer) => editModal.show(target);
  export const showRemove = (target: model.AlertPlayer) =>
    removeModal.show(target);
</script>

<AddAlertPlayerModal
  bind:this={addModal}
  defaultAlertPlayer={EMPTY}
  maxMemoLength={MAX_MEMO_LENGTH}
/>

<EditAlertPlayerModal
  bind:this={editModal}
  defaultAlertPlayer={EMPTY}
  maxMemoLength={MAX_MEMO_LENGTH}
/>

<RemoveAlertPlayerModal bind:this={removeModal} defaultAlertPlayer={EMPTY} />
