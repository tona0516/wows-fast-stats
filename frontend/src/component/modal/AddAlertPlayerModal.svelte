<script lang="ts">
  import {
    SearchPlayer,
    AlertPatterns,
    UpdateAlertPlayer,
  } from "wailsjs/go/main/App";
  import type { model } from "wailsjs/go/models";
  import clone from "clone";
  import UIkit from "uikit";
  import UkModal from "src/component/common/uikit/UkModal.svelte";
  import UkIcon from "src/component/common/uikit/UkIcon.svelte";
  import { ModalElementID } from "./ModalElementID";
  import { Notifier } from "src/lib/Notifier";
  import UkSpinner from "../common/uikit/UkSpinner.svelte";

  export let defaultAlertPlayer: model.AlertPlayer;
  export let maxMemoLength: number;

  export const show = () => {
    clean();
    const elem = document.getElementById(ModalElementID.ADD_ALERT_PLAYER);
    UIkit.modal(elem!).show();
  };

  const MIN_SEARCH_LENGTH = 3;

  const clean = () => {
    target = clone(defaultAlertPlayer);
    isSearching = false;
    searchInput = "";
    searchPlayers = {};
    searchResult = undefined;
  };

  const add = async (player: model.AlertPlayer) => {
    try {
      player.name = searchResult![0];
      player.account_id = searchResult![1];

      await UpdateAlertPlayer(player);
    } catch (error) {
      Notifier.failure(error);
    }
  };

  const searchPlayer = async (input: string) => {
    if (isSearching) {
      return;
    }

    if (input.length < MIN_SEARCH_LENGTH) {
      searchPlayers = {};
      return;
    }

    try {
      isSearching = true;
      searchPlayers = await SearchPlayer(input);
    } catch (error) {
      searchPlayers = {};
    } finally {
      isSearching = false;
    }
  };

  let target: model.AlertPlayer = clone(defaultAlertPlayer);
  let isSearching: boolean = false;
  let searchInput: string = "";
  let searchPlayers: { [key: string]: number } = {};
  let searchResult: [string, number] | undefined = undefined;

  $: disableAddButton =
    searchResult === undefined ||
    searchResult[0] === "" ||
    searchResult[1] === 0 ||
    target.pattern === "";
</script>

<UkModal id={ModalElementID.ADD_ALERT_PLAYER}>
  <div slot="body">
    <div class="uk-margin-small">
      <form class="uk-search uk-search-default">
        <input
          class="uk-search-input"
          type="search"
          placeholder="プレイヤー名"
          bind:value={searchInput}
          on:input={() => searchPlayer(searchInput)}
        />
      </form>

      {#if searchResult}
        <span class="uk-margin-small">
          <UkIcon name="check" />
          {searchResult[0]}
        </span>
      {/if}

      {#if searchInput.length >= MIN_SEARCH_LENGTH}
        <div class="uk-dropdown uk-position-absolute">
          {#if isSearching}
            <div class="uk-flex uk-flex-center">
              <UkSpinner />
            </div>
          {/if}

          <ul class="uk-nav uk-dropdown-nav">
            {#each Object.entries(searchPlayers) as player}
              <li>
                <!-- svelte-ignore a11y-invalid-attribute -->
                <a
                  href="#"
                  on:click={() => {
                    searchInput = "";
                    searchResult = player;
                    searchPlayers = {};
                  }}>{player[0]}</a
                >
              </li>
            {/each}
          </ul>
        </div>
      {/if}
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
      disabled={disableAddButton}
      on:click={() => add(target)}>追加</button
    >
  </div>
</UkModal>
