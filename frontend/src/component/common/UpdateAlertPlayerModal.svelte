<script lang="ts">
  import { onMount } from "svelte";
  import {
    AlertPatterns,
    UpdateAlertPlayer,
    SearchPlayer,
    ShowMessageDialog,
  } from "wailsjs/go/main/App";
  import {
    closeAlertPlayerModal,
    storedAlertPlayerForm,
    storedIsEditAlertPlayer,
    storedIsShowUpdateAlertPlayerModal,
  } from "src/stores";

  const MAX_MESSAGE_LENGTH = 100;

  let patterns: string[] = [];

  let playerSuggestions: { id: number; name: string }[] = [];

  onMount(async () => {
    patterns = await AlertPatterns();
  });

  async function onPlayerSearchInput(e: Event) {
    const value = (e.target as HTMLInputElement).value;
    if (value.length < 2) {
      playerSuggestions = [];
      return;
    }

    SearchPlayer(value).then((result) => {
      if (!result) {
        playerSuggestions = [];
        return;
      }

      playerSuggestions = result.map((player) => {
        return { id: player.account_id, name: player.nickname };
      });
    });
  }

  function selectPlayerSuggestion(s: { id: number; name: string }) {
    $storedAlertPlayerForm.account_id = s.id;
    $storedAlertPlayerForm.name = s.name;
    playerSuggestions = [];
  }

  async function save() {
    try {
      await UpdateAlertPlayer($storedAlertPlayerForm);
    } catch (error) {
      ShowMessageDialog("保存に失敗しました");
    }

    closeAlertPlayerModal();
  }
</script>

{#if $storedIsShowUpdateAlertPlayerModal}
  <dialog class="modal modal-open">
    <form method="dialog" class="modal-box" on:submit|preventDefault={save}>
      <h3 class="font-bold text-lg mb-4">
        {$storedIsEditAlertPlayer ? "編集" : "追加"}
      </h3>
      {#if !$storedIsEditAlertPlayer}
        <div class="form-control mb-2">
          <label class="label" for="player-search">プレイヤー名</label>
          <div class="dropdown w-full">
            <input
              id="player-search"
              class="input input-bordered w-full"
              type="text"
              bind:value={$storedAlertPlayerForm.name}
              on:input={onPlayerSearchInput}
              autocomplete="off"
              required
            />
            {#if playerSuggestions.length > 0}
              <ul
                class="dropdown-content menu bg-base-100 rounded-md border-1 border-neutral-300 shadow-lg mt-1 w-full z-50 max-h-60 overflow-y-auto"
              >
                {#each playerSuggestions as s}
                  <li>
                    <button
                      type="button"
                      class="w-full text-left"
                      on:click={() => selectPlayerSuggestion(s)}
                    >
                      {s.name}
                      <span class="text-xs text-gray-400">({s.id})</span>
                    </button>
                  </li>
                {/each}
              </ul>
            {/if}
          </div>
        </div>
      {/if}
      {#if $storedIsEditAlertPlayer}
        <div class="form-control mb-2">
          <label class="label" for="account-id">ID</label>
          <input
            id="account-id"
            class="input input-bordered"
            type="number"
            value={$storedAlertPlayerForm.account_id}
            readonly
          />
        </div>
        <div class="form-control mb-2">
          <label class="label" for="player-name">名前</label>
          <input
            id="player-name"
            class="input input-bordered"
            type="text"
            value={$storedAlertPlayerForm.name}
            readonly
          />
        </div>
      {/if}
      <div class="form-control mb-2">
        <label class="label" for="pattern">アイコン</label>
        <div class="flex flex-wrap gap-2">
          {#each patterns as p}
            <label class="cursor-pointer flex items-center gap-1">
              <input
                type="radio"
                name="pattern"
                class="radio radio-primary"
                value={p}
                bind:group={$storedAlertPlayerForm.pattern}
              />
              <i class={`bi ${p}`}></i>
            </label>
          {/each}
        </div>
      </div>
      <div class="form-control mb-2">
        <label class="label" for="message">メッセージ</label>
        <input
          id="message"
          class="input input-bordered"
          type="text"
          bind:value={$storedAlertPlayerForm.message}
          maxlength={MAX_MESSAGE_LENGTH}
        />
        <div class="text-right text-xs text-gray-500 mt-1">
          {$storedAlertPlayerForm.message.length}/{MAX_MESSAGE_LENGTH}文字
        </div>
      </div>
      <div class="modal-action">
        <button
          type="button"
          class="btn"
          on:click={() => closeAlertPlayerModal()}>キャンセル</button
        >
        <button
          type="submit"
          class="btn btn-primary"
          disabled={!$storedIsEditAlertPlayer &&
            (!$storedAlertPlayerForm.account_id ||
              !$storedAlertPlayerForm.name)}>保存</button
        >
      </div>
    </form>
  </dialog>
{/if}
