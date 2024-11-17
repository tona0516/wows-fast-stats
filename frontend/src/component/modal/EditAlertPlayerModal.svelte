<script lang="ts">
  import { AlertPatterns, UpdateAlertPlayer } from "wailsjs/go/main/App";
  import type { data } from "wailsjs/go/models";
  import UkModal from "src/component/common/uikit/UkModal.svelte";
  import UIkit from "uikit";
  import { ModalElementID } from "./ModalElementID";
  import clone from "clone";
  import { Notifier } from "src/lib/Notifier";

  export let defaultAlertPlayer: data.AlertPlayer;
  export let maxMemoLength: number;

  let target: data.AlertPlayer = clone(defaultAlertPlayer);
  $: disableUpdateButton =
    target.account_id === 0 || target.name === "" || target.pattern === "";

  export const show = (_target: data.AlertPlayer) => {
    target = _target;
    const elem = document.getElementById(ModalElementID.EDIT_ALERT_PLAYER);
    UIkit.modal(elem!).show();
  };

  const update = async () => {
    try {
      await UpdateAlertPlayer(target);
    } catch (error) {
      Notifier.failure(error);
    }
  };
</script>

<UkModal id={ModalElementID.EDIT_ALERT_PLAYER}>
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
        maxlength={maxMemoLength}
        bind:value={target.message}
      />
      <div>{target.message.length}/{maxMemoLength}</div>
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
