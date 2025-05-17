<script lang="ts">
import clone from "clone";
import UkModal from "src/component/common/uikit/UkModal.svelte";
import { Notifier } from "src/lib/Notifier";
import UIkit from "uikit";
import { RemoveAlertPlayer } from "wailsjs/go/main/App";
import type { data } from "wailsjs/go/models";
import { ModalElementID } from "./ModalElementID";

export let defaultAlertPlayer: data.AlertPlayer;

export const show = (_target: data.AlertPlayer) => {
  target = _target;

  const elem = document.getElementById(ModalElementID.REMOVE_ALERT_PLAYER);
  if (elem) {
    UIkit.modal(elem).show();
  }
};

const remove = async () => {
  try {
    await RemoveAlertPlayer(target.account_id);
  } catch (error) {
    Notifier.failure(error);
  }
};

let target: data.AlertPlayer = clone(defaultAlertPlayer);
</script>

<UkModal id={ModalElementID.REMOVE_ALERT_PLAYER}>
  <div slot="body">
    <p>「{target.name}」を削除しますか？</p>
  </div>

  <div slot="footer">
    <button
      class="uk-button uk-button-danger uk-modal-close"
      type="button"
      on:click={() => remove()}>削除</button
    >
  </div>
</UkModal>
