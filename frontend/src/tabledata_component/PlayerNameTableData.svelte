<script lang="ts">
  import type { PlayerName } from "src/lib/column/PlayerName";
  import { storedAlertPlayers, storedExcludePlayerIDs } from "src/stores";
  import { createEventDispatcher } from "svelte";
  import {
    AddExcludePlayerID,
    RemoveExcludePlayerID,
  } from "wailsjs/go/main/App";
  import type { domain } from "wailsjs/go/models";
  import { BrowserOpenURL } from "wailsjs/runtime/runtime";

  export let column: PlayerName;
  export let player: domain.Player;

  $: accountID = player.player_info.id;
  $: isNPC = accountID === 0;
  $: isBelongToClan = player.player_info.clan.id !== 0;
  $: isChecked = !$storedExcludePlayerIDs.includes(accountID);
  $: alertPlayer = $storedAlertPlayers.find(
    (it) => it.account_id === accountID,
  );

  const dispatch = createEventDispatcher();

  const onCheck = async (e: any) => {
    if (e.target.checked) {
      await RemoveExcludePlayerID(accountID);
    } else {
      await AddExcludePlayerID(accountID);
    }
    dispatch("CheckPlayer");
  };
</script>

<td>
  {#if column.isShowCheckBox(player)}
    <input
      class="uk-checkbox"
      type="checkbox"
      on:click={onCheck}
      checked={isChecked}
    />
  {/if}
</td>

<td style="background-color: {column.bgColorCode(player)}">
  <div class="td-string omit">
    {#if isNPC}
      {column.displayValue(player)}
    {:else}
      <!-- svelte-ignore a11y-invalid-attribute -->
      <a href="#" uk-tooltip={alertPlayer?.message}>
        {#if alertPlayer}
          <i class="bi {alertPlayer.pattern}" />
        {/if}
        {column.displayValue(player)}
        <span uk-drop-parent-icon></span>
      </a>
    {/if}
  </div>

  {#if !isNPC}
    <div uk-dropdown="mode: click" uk-toggle>
      <ul class="uk-nav uk-dropdown-nav">
        {#if isBelongToClan}
          <li>
            <!-- svelte-ignore a11y-invalid-attribute -->
            <a href="#" on:click={() => BrowserOpenURL(column.clanURL(player))}
              >クラン詳細(WoWS Stats & Numbers)</a
            >
          </li>
        {/if}

        <li>
          <!-- svelte-ignore a11y-invalid-attribute -->
          <a href="#" on:click={() => BrowserOpenURL(column.playerURL(player))}
            >プレイヤー詳細(WoWS Stats & Numbers)</a
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
      </ul>
    </div>
  {/if}
</td>
