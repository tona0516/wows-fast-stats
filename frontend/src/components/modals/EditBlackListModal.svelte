<script lang="ts">
  import {
    showToast,
    storedBlackList,
    storedEditBlackListItem,
  } from "@libs/stores";
  import { SearchPlayer, UpdateBlackList } from "@wails/go/main/App";
  import type { data } from "@wails/go/models";
  import ModalCommon from "./ModalCommon.svelte";
  import { ModalManager } from "@libs/ModalManager";

  const MAX_MESSAGE_LENGTH = 100;
  const BLACKLIST_PATTERNS = [
    "bi-check-circle-fill",
    "bi-exclamation-triangle-fill",
    "bi-patch-question-fill",
    "bi-1-square-fill",
    "bi-2-square-fill",
    "bi-3-square-fill",
  ];
  const Z_VALUE = 51;

  let suggestedPlayers: data.WGAccountListData[] = [];

  function includesBlackList(accountID: number): boolean {
    for (const item of $storedBlackList) {
      if (item.account_id === accountID) {
        return true;
      }
    }
    return false;
  }

  async function searchPlayer(e: Event) {
    const input = (e.target as HTMLInputElement).value;
    if (input.length < 2) {
      suggestedPlayers = [];
      return;
    }

    try {
      const result = await SearchPlayer(input);
      suggestedPlayers = result.filter(
        (player) => !includesBlackList(player.account_id),
      );
    } catch (error) {
      suggestedPlayers = [];
    }
  }

  function selectPlayer(player: data.WGAccountListData) {
    if (!$storedEditBlackListItem) {
      return;
    }

    $storedEditBlackListItem.form.account_id = player.account_id;
    $storedEditBlackListItem.form.name = player.nickname;
    suggestedPlayers = [];
  }

  async function save() {
    if (!$storedEditBlackListItem) {
      ModalManager.instance.closeForEdit();
      return;
    }

    try {
      await UpdateBlackList($storedEditBlackListItem.form);
      showToast("保存しました");
    } catch (error) {
      showToast("保存に失敗しました");
    } finally {
      ModalManager.instance.closeForEdit();
    }
  }
</script>

{#if $storedEditBlackListItem}
  <ModalCommon zValue={Z_VALUE} close={ModalManager.instance.closeForEdit}>
    <h3 class="font-bold text-lg mb-4">
      アラートプレイヤー{$storedEditBlackListItem.mode === "edit"
        ? "編集"
        : "追加"}
    </h3>

    {#if $storedEditBlackListItem.mode === "create"}
      <fieldset class="fieldset">
        <legend class="fieldset-legend">プレイヤー名</legend>
        <div class="dropdown">
          <label class="input">
            <i class="bi bi-search"></i>
            <input
              id="player-search"
              class="grow"
              type="text"
              bind:value={$storedEditBlackListItem.form.name}
              on:input={searchPlayer}
              autocomplete="off"
            />
          </label>
          {#if suggestedPlayers.length > 0}
            <ul
              class="dropdown-content menu bg-base-100 rounded-md border-1 border-neutral-300 shadow-lg mt-1 w-full z-{Z_VALUE} max-h-60 overflow-y-auto"
            >
              {#each suggestedPlayers as player}
                <li>
                  <button
                    type="button"
                    class="w-full text-left"
                    on:click={() => selectPlayer(player)}
                  >
                    {player.nickname}
                    <span class="text-xs text-gray-400"
                      >({player.account_id})</span
                    >
                  </button>
                </li>
              {/each}
            </ul>
          {/if}
        </div>
      </fieldset>
    {:else}
      <fieldset class="fieldset">
        <legend class="fieldset-legend">プレイヤー名</legend>
        <input
          id="player-name"
          class="input input-bordered"
          type="text"
          value={$storedEditBlackListItem.form.name}
          readonly
        />
      </fieldset>
    {/if}

    <fieldset class="fieldset">
      <legend class="fieldset-legend">アイコン</legend>
      <div class="flex flex-wrap gap-4">
        {#each BLACKLIST_PATTERNS as pattern}
          <label class="flex items-center gap-1">
            <input
              type="radio"
              class="radio"
              value={pattern}
              bind:group={$storedEditBlackListItem.form.pattern}
            />
            <i class="{`bi ${pattern}`} text-lg"></i>
          </label>
        {/each}
      </div>
    </fieldset>

    <fieldset class="fieldset">
      <legend class="fieldset-legend">メモ(任意)</legend>
      <input
        id="message"
        class="input input-bordered"
        type="text"
        bind:value={$storedEditBlackListItem.form.message}
        maxlength={MAX_MESSAGE_LENGTH}
      />
      <div class="text-xs text-gray-500 mt-1">
        {$storedEditBlackListItem.form.message.length}/{MAX_MESSAGE_LENGTH}文字
      </div>
    </fieldset>

    <div class="modal-action">
      <button
        type="button"
        class="btn btn-primary"
        disabled={$storedEditBlackListItem.form.account_id === 0}
        on:click={save}>保存</button
      >
    </div>
  </ModalCommon>
{/if}
