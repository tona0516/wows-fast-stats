<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import {
    SearchPlayer,
    AlertPatterns,
    UpdateAlertPlayer,
  } from "wailsjs/go/main/App";
  import type { data } from "wailsjs/go/models";
  import clone from "clone";
  import UIkit from "uikit";
  import UkModal from "src/component/common/uikit/UkModal.svelte";
  import UkIcon from "src/component/common/uikit/UkIcon.svelte";
  import { ModalElementID } from "./ModalElementID";

  export let defaultAlertPlayer: data.AlertPlayer;
  export let maxMemoLength: number;
  const dispatch = createEventDispatcher();

  export const show = () => {
    clean();
    const elem = document.getElementById(ModalElementID.ADD_ALERT_PLAYER);
    UIkit.modal(elem!).show();
  };

  const clean = () => {
    target = clone(defaultAlertPlayer);
    isSearching = false;
    searchInput = "";
    searchPlayers = [];
    searchResult = undefined;
  };

  const add = async (player: data.AlertPlayer) => {
    try {
      player.account_id = searchResult!.account_id;
      player.name = searchResult!.nickname;

      await UpdateAlertPlayer(player);
      dispatch("Success");
    } catch (error) {
      dispatch("Failure", { message: error });
    }
  };

  const searchPlayer = async (input: string) => {
    if (isSearching) {
      return;
    }

    if (input.length < 3) {
      searchPlayers = [];
      return;
    }

    if (input === "") {
      searchPlayers = [];
      return;
    }

    try {
      isSearching = true;
      searchPlayers = await SearchPlayer(input);
    } catch (error) {
      searchPlayers = [];
      return;
    } finally {
      isSearching = false;
    }
  };

  let target: data.AlertPlayer = clone(defaultAlertPlayer);
  let isSearching: boolean = false;
  let searchInput: string = "";
  let searchPlayers: data.WGAccountListData[] = [];
  let searchResult: data.WGAccountListData | undefined = undefined;

  $: disableAddButton =
    searchResult === undefined ||
    searchResult.account_id === 0 ||
    searchResult.nickname === "" ||
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
          {searchResult.nickname}
        </span>
      {/if}

      {#if searchPlayers.length !== 0}
        <div class="uk-dropdown uk-position-absolute">
          <ul class="uk-nav uk-dropdown-nav">
            {#each searchPlayers as player}
              <li>
                <!-- svelte-ignore a11y-invalid-attribute -->
                <a
                  href="#"
                  on:click={() => {
                    searchResult = player;
                    searchPlayers = [];
                  }}>{player.nickname}</a
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
