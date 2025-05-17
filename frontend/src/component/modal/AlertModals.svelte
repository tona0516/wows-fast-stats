<script lang="ts">
import AddAlertPlayerModal from "src/component/modal/AddAlertPlayerModal.svelte";
import EditAlertPlayerModal from "src/component/modal/EditAlertPlayerModal.svelte";
import RemoveAlertPlayerModal from "src/component/modal/RemoveAlertPlayerModal.svelte";
import type { data } from "wailsjs/go/models";

const MAX_MEMO_LENGTH = 100;
const EMPTY: data.AlertPlayer = {
  account_id: 0,
  name: "",
  pattern: "bi-check-circle-fill",
  message: "",
} as const;
let addModal: AddAlertPlayerModal;
let editModal: EditAlertPlayerModal;
let removeModal: RemoveAlertPlayerModal;

export const showAdd = () => addModal.show();
export const showEdit = (target: data.AlertPlayer) => editModal.show(target);
export const showRemove = (target: data.AlertPlayer) =>
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
