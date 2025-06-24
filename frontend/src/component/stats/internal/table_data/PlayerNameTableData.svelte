<script lang="ts">
  import ExternalLink from "src/component/common/ExternalLink.svelte";
  import { Notifier } from "src/lib/Notifier";
  import { NumbersURL } from "src/lib/NumbersURL";
  import type { PlayerName } from "src/lib/column/model/PlayerName";
  import {
    showDeleteAlertPlayerModal,
    showUpdateAlertPlayerModal,
    storedAlertPlayers,
  } from "src/stores";
  import type { data } from "wailsjs/go/models";
  import { ClipboardSetText } from "wailsjs/runtime/runtime";

  export let column: PlayerName;
  export let player: data.Player;

  $: accountID = player.player_info.id;
  $: alertPlayer = $storedAlertPlayers.find(
    (it) => it.account_id === accountID,
  );
  $: clanTag = column.clanTag(player);
  $: isNPC = column.isNPC(player);

  const setPlayerNameToClipboard = async () => {
    const isSuccess = await ClipboardSetText(player.player_info.name);

    isSuccess
      ? Notifier.success("コピーしました！")
      : Notifier.failure("コピーに失敗しました");
  };
</script>

<div>
  {#if isNPC}
    {column.playerName(player)}
  {:else}
    <div class="flex items-center">
      <div class="dropdown">
        <div class="tooltip" data-tip={alertPlayer?.message}>
          <!-- svelte-ignore a11y-no-noninteractive-tabindex -->
          <div tabindex="0">
            <!-- svelte-ignore a11y-invalid-attribute -->
            <a href="#">
              <div class="flex">
                {#if alertPlayer}
                  <i class="bi {alertPlayer.pattern}" />
                {/if}
                {#if clanTag}
                  {#if column.clanFlagIconClass(player)}
                    <span class={column.clanFlagIconClass(player)}></span>
                  {/if}
                  <span style="color: {column.clanColorCode(player)}">
                    {clanTag}
                  </span>
                {/if}
                <div style="color: {column.textColorCode(player)}">
                  {column.playerName(player)}
                </div>
              </div>
            </a>
          </div>
        </div>

        <!-- svelte-ignore a11y-no-noninteractive-tabindex -->
        <ul
          tabindex="0"
          class="dropdown-content menu bg-base-100 rounded-md border-1 border-neutral-500 shadow-lg p-2"
        >
          {#if clanTag}
            <li>
              <ExternalLink url={NumbersURL.clan(player.player_info.clan.id)}>
                <span class="text-nowrap">クラン詳細(WoWS Stats & Numbers)</span
                >
              </ExternalLink>
            </li>
          {/if}

          <li>
            <ExternalLink
              url={NumbersURL.player(
                player.player_info.id,
                player.player_info.name,
              )}
            >
              <span class="text-nowrap"
                >プレイヤー詳細(WoWS Stats & Numbers)</span
              >
            </ExternalLink>
          </li>

          <li>
            <!-- svelte-ignore a11y-invalid-attribute -->
            <a
              href="#"
              on:click={() => {
                showUpdateAlertPlayerModal(
                  player.player_info.id,
                  player.player_info.name,
                );
                if (document.activeElement instanceof HTMLElement) {
                  document.activeElement.blur();
                }
              }}
            >
              <span class="text-nowrap">アラートプレイヤーへ追加する</span>
            </a>
          </li>

          {#if alertPlayer}
            <li>
              <!-- svelte-ignore a11y-invalid-attribute -->
              <a
                href="#"
                on:click={() => {
                  showDeleteAlertPlayerModal(alertPlayer.account_id);
                  if (document.activeElement instanceof HTMLElement) {
                    document.activeElement.blur();
                  }
                }}
              >
                <span class="text-nowrap">アラートプレイヤーから削除する</span>
              </a>
            </li>
          {/if}

          <li>
            <!-- svelte-ignore a11y-invalid-attribute -->
            <a
              href="#"
              on:click={() => {
                setPlayerNameToClipboard();
                if (document.activeElement instanceof HTMLElement) {
                  document.activeElement.blur();
                }
              }}
            >
              <span class="text-nowrap"
                >プレイヤー名をクリップボードにコピーする</span
              >
            </a>
          </li>
        </ul>
      </div>
    </div>
  {/if}
</div>
