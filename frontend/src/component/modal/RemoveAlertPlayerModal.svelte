<script lang="ts">
  import { RemoveAlertPlayer } from "wailsjs/go/main/App";
  import type { model } from "wailsjs/go/models";
  import UkModal from "src/component/common/uikit/UkModal.svelte";
  import UIkit from "uikit";
  import { ModalElementID } from "./ModalElementID";
  import clone from "clone";
  import { Notifier } from "src/lib/Notifier";

  export let defaultAlertPlayer: model.AlertPlayer;

  export const show = (_target: model.AlertPlayer) => {
    target = _target;

    const elem = document.getElementById(ModalElementID.REMOVE_ALERT_PLAYER);
    UIkit.modal(elem!).show();
  };

  const remove = async () => {
    try {
      await RemoveAlertPlayer(target.account_id);
    } catch (error) {
      Notifier.failure(error);
    }
  };

  let target: model.AlertPlayer = clone(defaultAlertPlayer);
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
