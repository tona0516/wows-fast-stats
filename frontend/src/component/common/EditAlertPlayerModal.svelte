<script lang="ts">
  import {
    EditAlertPlayerModal,
    showToast,
    storedAlertPlayers,
    storedEditAlertPlayer,
  } from "src/stores";
  import { onMount } from "svelte";
  import {
    AlertPatterns,
    UpdateAlertPlayer,
    SearchPlayer,
  } from "wailsjs/go/main/App";

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

      const alertPlayerIDs = $storedAlertPlayers.map((ap) => ap.account_id);

      playerSuggestions = result
        .filter((p) => !alertPlayerIDs.includes(p.account_id))
        .map((p) => {
          return { id: p.account_id, name: p.nickname };
        });
    });
  }

  function selectPlayerSuggestion(s: { id: number; name: string }) {
    if (!$storedEditAlertPlayer) {
      return;
    }

    $storedEditAlertPlayer.form.account_id = s.id;
    $storedEditAlertPlayer.form.name = s.name;
    playerSuggestions = [];
  }

  async function save() {
    if (!$storedEditAlertPlayer) {
      return;
    }

    try {
      await UpdateAlertPlayer($storedEditAlertPlayer.form);
    } catch (error) {
      showToast("保存に失敗しました");
    }

    EditAlertPlayerModal.close();
  }
</script>

{#if $storedEditAlertPlayer}
  <dialog class="modal modal-open z-51">
    <form method="dialog" class="modal-box" on:submit|preventDefault={save}>
      <h3 class="font-bold text-lg mb-4">
        {$storedEditAlertPlayer.mode === "edit" ? "編集" : "追加"}
      </h3>

      {#if $storedEditAlertPlayer.mode === "create"}
        <div class="form-control mb-2">
          <label class="label" for="player-search">プレイヤー名</label>
          <div class="dropdown w-full">
            <input
              id="player-search"
              class="input input-bordered w-full"
              type="text"
              bind:value={$storedEditAlertPlayer.form.name}
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
      {:else}
        <div class="form-control mb-2">
          <label class="label" for="account-id">ID</label>
          <input
            id="account-id"
            class="input input-bordered"
            type="number"
            value={$storedEditAlertPlayer.form.account_id}
            readonly={$storedEditAlertPlayer.mode === "specify"}
          />
        </div>

        <div class="form-control mb-2">
          <label class="label" for="player-name">プレイヤー名</label>
          <input
            id="player-name"
            class="input input-bordered"
            type="text"
            value={$storedEditAlertPlayer.form.name}
            readonly={$storedEditAlertPlayer.mode === "specify"}
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
                bind:group={$storedEditAlertPlayer.form.pattern}
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
          bind:value={$storedEditAlertPlayer.form.message}
          maxlength={MAX_MESSAGE_LENGTH}
        />
        <div class="text-right text-xs text-gray-500 mt-1">
          {$storedEditAlertPlayer.form.message.length}/{MAX_MESSAGE_LENGTH}文字
        </div>
      </div>

      <div class="modal-action">
        <button
          type="button"
          class="btn"
          on:click={() => EditAlertPlayerModal.close()}>キャンセル</button
        >
        <button
          type="submit"
          class="btn btn-primary"
          disabled={$storedEditAlertPlayer.form.account_id === 0 ||
            $storedEditAlertPlayer.form.name.length === 0}>保存</button
        >
      </div>
    </form>
  </dialog>
{/if}
