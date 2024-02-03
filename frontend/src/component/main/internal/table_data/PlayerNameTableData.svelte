<script lang="ts">
  import ExternalLink from "src/component/common/ExternalLink.svelte";
  import UkDowndown from "src/component/common/uikit/UkDowndown.svelte";
  import UkTooltip from "src/component/common/uikit/UkTooltip.svelte";
  import { NumbersURL } from "src/lib/NumbersURL";
  import { CssClass } from "src/lib/CssClass";
  import { storedAlertPlayers, storedExcludedPlayers } from "src/stores";
  import { createEventDispatcher } from "svelte";
  import {
    AddExcludePlayerID,
    RemoveExcludePlayerID,
  } from "wailsjs/go/main/App";
  import type { model } from "wailsjs/go/models";
  import type { PlayerName } from "src/lib/column/model/PlayerName";
  import { ClipboardSetText } from "wailsjs/runtime/runtime";
  import { Notifier } from "src/lib/Notifier";

  export let column: PlayerName;
  export let player: model.Player;

  $: accountID = player.player_info.id;
  $: isChecked = !$storedExcludedPlayers.includes(accountID);
  $: alertPlayer = $storedAlertPlayers.find(
    (it) => it.account_id === accountID,
  );
  $: clanTag = column.clanTag(player);
  $: isNPC = column.isNPC(player);

  const dispatch = createEventDispatcher();

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const onCheck = async (e: any) => {
    if (e.target.checked) {
      await RemoveExcludePlayerID(accountID);
    } else {
      await AddExcludePlayerID(accountID);
    }
    dispatch("CheckPlayer");
  };

  const setPlayerNameToClipboard = async () => {
    const isSuccess = await ClipboardSetText(player.player_info.name);

    isSuccess
      ? Notifier.success("コピーしました！")
      : Notifier.failure("コピーに失敗しました");
  };
</script>

<td>
  {#if !isNPC}
    <input
      class="uk-checkbox"
      type="checkbox"
      on:click={onCheck}
      checked={isChecked}
    />
  {/if}
</td>

<td class={CssClass.TD_STR}>
  {#if !isNPC}
    <UkTooltip tooltip={alertPlayer?.message}>
      <!-- svelte-ignore a11y-invalid-attribute -->
      <a href="#">
        <div class="uk-flex">
          {#if alertPlayer}
            <i class="bi {alertPlayer.pattern}" />
          {/if}
          {#if clanTag}
            <div class="clan-tag" style="color: {column.clanColorCode(player)}">
              {clanTag}
            </div>
          {/if}
          <div
            class="uk-text-truncate"
            style="color: {column.textColorCode(player)}"
          >
            {column.playerName(player)}
          </div>
        </div>
      </a>
    </UkTooltip>

    <UkDowndown>
      <ul class="uk-nav uk-dropdown-nav">
        {#if clanTag}
          <li>
            <ExternalLink url={NumbersURL.clan(player.player_info.clan.id)}
              >クラン詳細(WoWS Stats & Numbers)</ExternalLink
            >
          </li>
        {/if}

        <li>
          <ExternalLink url={NumbersURL.player(player.player_info.id)}
            >プレイヤー詳細(WoWS Stats & Numbers)</ExternalLink
          >
        </li>

        {#if alertPlayer}
          {@const target = alertPlayer}
          <li>
            <!-- svelte-ignore a11y-invalid-attribute -->
            <a
              href="#"
              on:click={() => dispatch("EditAlertPlayer", { target: target })}
              >プレイヤーリストを編集する</a
            >
          </li>
          <li>
            <!-- svelte-ignore a11y-invalid-attribute -->
            <a
              href="#"
              on:click={() => dispatch("RemoveAlertPlayer", { target: target })}
              >プレイヤーリストから削除する</a
            >
          </li>
        {:else}
          <li>
            <!-- svelte-ignore a11y-invalid-attribute -->
            <a
              href="#"
              on:click={() =>
                dispatch("EditAlertPlayer", {
                  target: {
                    account_id: accountID,
                    name: player.player_info.name,
                    pattern: "bi-check-circle-fill",
                    message: "",
                  },
                })}>プレイヤーリストへ追加する</a
            >
          </li>
        {/if}
        <li>
          <!-- svelte-ignore a11y-invalid-attribute -->
          <a href="#" on:click={setPlayerNameToClipboard}
            >プレイヤー名をクリップボードにコピーする</a
          >
        </li>
      </ul>
    </UkDowndown>
  {:else}
    <div class="uk-text-truncate">{column.playerName(player)}</div>
  {/if}
</td>

<style>
  :global(.clan-tag) {
    margin-right: 2px;
  }
</style>
