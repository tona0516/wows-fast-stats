<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import {
    SearchPlayer,
    AlertPatterns,
    UpdateAlertPlayer,
  } from "wailsjs/go/main/App";
  import type { domain } from "wailsjs/go/models";
  import {
    ADD_ALERT_PLAYER_MODAL_ID,
    EMPTY_ALERT_PLAYER,
    MAX_MEMO_LENGTH,
  } from "src/lib/types";
  import clone from "clone";
  import ModalCommon from "./ModalCommon.svelte";
  import UIkit from "uikit";

  let target: domain.AlertPlayer = clone(EMPTY_ALERT_PLAYER);
  let isSearching: boolean = false;
  let searchInput: string = "";
  let searchPlayers: domain.WGAccountListData[] = [];
  let searchResult: domain.WGAccountListData | undefined = undefined;

  $: disableAddButton =
    searchResult === undefined ||
    searchResult.account_id === 0 ||
    searchResult.nickname === "" ||
    target.pattern === "";

  const dispatch = createEventDispatcher();

  export const show = () => {
    clean();

    const elem = document.getElementById(ADD_ALERT_PLAYER_MODAL_ID);
    UIkit.modal(elem!).show();
  };

  const clean = () => {
    target = clone(EMPTY_ALERT_PLAYER);
    isSearching = false;
    searchInput = "";
    searchPlayers = [];
    searchResult = undefined;
  };

  const add = async (player: domain.AlertPlayer) => {
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
</script>

<ModalCommon id={ADD_ALERT_PLAYER_MODAL_ID}>
  <div slot="body">
    <div class="uk-margin-small">
      <form class="uk-search uk-search-default">
        <span uk-search-icon></span>
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
          <span uk-icon="check"></span>
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
      disabled={disableAddButton}
      on:click={() => add(target)}>追加</button
    >
  </div>
</ModalCommon>
