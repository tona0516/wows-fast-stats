<script lang="ts">
  import { onMount } from "svelte";
  import { data } from "wailsjs/go/models";
  import {
    AlertPatterns,
    UpdateAlertPlayer,
    RemoveAlertPlayer,
    SearchPlayer,
  } from "wailsjs/go/main/App";
  import { Notifier } from "src/lib/Notifier";
  import { storedAlertPlayers } from "src/stores";

  let patterns: string[] = [];
  let showModal = false;
  let isEdit = false;
  let editIndex: number | null = null;
  let form = new data.AlertPlayer({
    account_id: 0,
    name: "",
    pattern: "",
    message: "",
  });

  let playerSearch = "";
  let playerSuggestions: { id: number; name: string }[] = [];

  onMount(async () => {
    patterns = await AlertPatterns();
  });

  async function onPlayerSearchInput(e: Event) {
    const value = (e.target as HTMLInputElement).value;
    playerSearch = value;
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
    form.account_id = s.id;
    form.name = s.name;
    playerSearch = s.name;
    playerSuggestions = [];
  }

  function openAddModal() {
    isEdit = false;
    form = new data.AlertPlayer({
      account_id: 0,
      name: "",
      pattern: patterns[0] ?? "",
      message: "",
    });
    playerSearch = "";
    playerSuggestions = [];
    showModal = true;
  }

  function openEditModal(index: number) {
    isEdit = true;
    editIndex = index;
    form = new data.AlertPlayer({ ...$storedAlertPlayers[index] });
    playerSearch = form.name;
    playerSuggestions = [];
    showModal = true;
  }

  function closeModal() {
    showModal = false;
    editIndex = null;
  }

  async function saveAlertPlayer() {
    try {
      await UpdateAlertPlayer(form);
      Notifier.success(isEdit ? "編集しました" : "追加しました");
      closeModal();
    } catch (e) {
      Notifier.failure("保存に失敗しました");
    }
  }

  let showDeleteModal = false;
  let confirmDeleteIndex: number | null = null;

  async function deleteAlertPlayer(index: number) {
    confirmDeleteIndex = index;
    showDeleteModal = true;
  }

  async function confirmDelete() {
    if (confirmDeleteIndex !== null) {
      try {
        await RemoveAlertPlayer(
          $storedAlertPlayers[confirmDeleteIndex].account_id,
        );
        Notifier.success("削除しました");
      } catch (e) {
        Notifier.failure("削除に失敗しました");
      }
    }
    showDeleteModal = false;
    confirmDeleteIndex = null;
  }

  function cancelDelete() {
    showDeleteModal = false;
    confirmDeleteIndex = null;
  }
</script>

<div class="mb-4 flex justify-end">
  <button class="btn btn-primary" on:click={openAddModal}>追加</button>
</div>

<table class="table w-full text-center">
  <thead>
    <tr>
      <th class="text-center">ID</th>
      <th class="text-center">プレイヤー名</th>
      <th class="text-center">アイコン</th>
      <th class="text-center">メッセージ</th>
      <th class="text-center">操作</th>
    </tr>
  </thead>
  <tbody>
    {#each $storedAlertPlayers as player, i}
      <tr>
        <td class="text-center">{player.account_id}</td>
        <td class="text-center">{player.name}</td>
        <td class="text-center"><i class={`bi ${player.pattern}`}></i></td>
        <td class="text-center">{player.message}</td>
        <td class="text-center">
          <button
            class="btn btn-sm btn-info mr-2"
            on:click={() => openEditModal(i)}>編集</button
          >
          <button
            class="btn btn-sm btn-error"
            on:click={() => deleteAlertPlayer(i)}>削除</button
          >
        </td>
      </tr>
    {/each}
  </tbody>
</table>

{#if showModal}
  <dialog class="modal modal-open">
    <form
      method="dialog"
      class="modal-box"
      on:submit|preventDefault={saveAlertPlayer}
    >
      <h3 class="font-bold text-lg mb-4">
        {isEdit ? "編集" : "追加"}
      </h3>
      {#if !isEdit}
        <div class="form-control mb-2">
          <label class="label" for="player-search">プレイヤー名</label>
          <div class="dropdown w-full">
            <input
              id="player-search"
              class="input input-bordered w-full"
              type="text"
              bind:value={playerSearch}
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
      {#if isEdit}
        <div class="form-control mb-2">
          <label class="label" for="account-id">ID</label>
          <input
            id="account-id"
            class="input input-bordered"
            type="number"
            value={form.account_id}
            readonly
          />
        </div>
        <div class="form-control mb-2">
          <label class="label" for="player-name">名前</label>
          <input
            id="player-name"
            class="input input-bordered"
            type="text"
            value={form.name}
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
                bind:group={form.pattern}
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
          bind:value={form.message}
          maxlength="30"
        />
        <div class="text-right text-xs text-gray-500 mt-1">
          {form.message.length}/30
        </div>
      </div>
      <div class="modal-action">
        <button type="button" class="btn" on:click={closeModal}
          >キャンセル</button
        >
        <button
          type="submit"
          class="btn btn-primary"
          disabled={!isEdit && (!form.account_id || !form.name)}>保存</button
        >
      </div>
    </form>
  </dialog>
{/if}

{#if showDeleteModal}
  <dialog class="modal modal-open">
    <form
      method="dialog"
      class="modal-box"
      on:submit|preventDefault={confirmDelete}
    >
      <h3 class="font-bold text-lg mb-4">本当に削除しますか？</h3>
      {#if confirmDeleteIndex !== null}
        <div class="mb-4 text-center">
          <!-- 削除時にクラッシュ -->
          <span class="font-bold"
            >{$storedAlertPlayers[confirmDeleteIndex].name}</span
          >
          <span class="ml-2 text-xs text-gray-500"
            >(ID: {$storedAlertPlayers[confirmDeleteIndex].account_id})</span
          >
        </div>
      {/if}
      <div class="modal-action">
        <button type="button" class="btn" on:click={cancelDelete}
          >キャンセル</button
        >
        <button type="submit" class="btn btn-error">削除</button>
      </div>
    </form>
  </dialog>
{/if}
