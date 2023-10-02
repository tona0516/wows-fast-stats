<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { AlertPatterns, UpdateAlertPlayer } from "wailsjs/go/main/App";
  import type { domain } from "wailsjs/go/models";
  import {
    MAX_MEMO_LENGTH,
    EDIT_ALERT_PLAYER_MODAL_ID,
    EMPTY_ALERT_PLAYER,
  } from "src/lib/types";
  import UkModal from "src/component/common/uikit/UkModal.svelte";
  import UIkit from "uikit";
  import clone from "clone";

  $: disableUpdateButton =
    target.account_id === 0 || target.name === "" || target.pattern === "";

  let target: domain.AlertPlayer = clone(EMPTY_ALERT_PLAYER);

  export const show = (_target: domain.AlertPlayer) => {
    target = _target;

    const elem = document.getElementById(EDIT_ALERT_PLAYER_MODAL_ID);
    UIkit.modal(elem!).show();
  };

  const dispatch = createEventDispatcher();

  const update = async () => {
    try {
      await UpdateAlertPlayer(target);
      dispatch("Success");
    } catch (error) {
      dispatch("Failure", { message: error });
    }
  };
</script>

<UkModal id={EDIT_ALERT_PLAYER_MODAL_ID}>
  <div slot="body">
    <div class="uk-margin-small">
      <input
        class="uk-input"
        type="text"
        placeholder="プレイヤー名"
        bind:value={target.name}
        disabled
      />
    </div>

    <div class="uk-margin-small">
      <div>アイコン</div>
      {#await AlertPatterns() then alertPatterns}
        <div class="uk-grid-small uk-child-width-auto uk-grid">
          {#each alertPatterns as pattern}
            <label
              ><input
                class="uk-radio"
                type="radio"
                bind:group={target.pattern}
                value={pattern}
              /> <i class="bi {pattern}" /></label
            >
          {/each}
        </div>
      {/await}
    </div>

    <div class="uk-margin-small">
      <textarea
        class="uk-textarea"
        placeholder="メモ(任意)"
        maxlength={MAX_MEMO_LENGTH}
        bind:value={target.message}
      />
      <div>{target.message.length}/{MAX_MEMO_LENGTH}</div>
    </div>
  </div>

  <div slot="footer">
    <button
      class="uk-button uk-button-primary uk-modal-close"
      type="button"
      disabled={disableUpdateButton}
      on:click={() => update()}>更新</button
    >
  </div>
</UkModal>
