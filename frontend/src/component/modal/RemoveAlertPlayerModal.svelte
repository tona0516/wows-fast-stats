<script lang="ts">
  import {
    EMPTY_ALERT_PLAYER,
    REMOVE_ALERT_PLAYER_MODAL_ID,
  } from "src/lib/types";
  import { createEventDispatcher } from "svelte";
  import { RemoveAlertPlayer } from "wailsjs/go/main/App";
  import type { domain } from "wailsjs/go/models";
  import UkModal from "src/component/common/uikit/UkModal.svelte";
  import UIkit from "uikit";
  import clone from "clone";

  let target: domain.AlertPlayer = clone(EMPTY_ALERT_PLAYER);

  export const show = (_target: domain.AlertPlayer) => {
    target = _target;

    const elem = document.getElementById(REMOVE_ALERT_PLAYER_MODAL_ID);
    UIkit.modal(elem!).show();
  };

  const dispatch = createEventDispatcher();

  const remove = async () => {
    try {
      await RemoveAlertPlayer(target.account_id);
      dispatch("Success");
    } catch (error) {
      dispatch("Failure", { message: error });
    }
  };
</script>

<UkModal id={REMOVE_ALERT_PLAYER_MODAL_ID}>
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
