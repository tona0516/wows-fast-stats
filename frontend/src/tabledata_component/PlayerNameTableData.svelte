<script lang="ts">
  import clone from "clone";
  import type { PlayerName } from "src/lib/column/PlayerName";
  import { storedAlertPlayers, storedExcludePlayerIDs } from "src/stores";
  import { createEventDispatcher } from "svelte";
  import { Tooltip } from "sveltestrap";
  import {
    AddExcludePlayerID,
    RemoveExcludePlayerID,
  } from "wailsjs/go/main/App";
  import type { domain } from "wailsjs/go/models";
  import { BrowserOpenURL } from "wailsjs/runtime/runtime";

  export let column: PlayerName;
  export let player: domain.Player;

  $: accountID = player.player_info.id;
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

<td class="td-checkbox">
  {#if column.isShowCheckBox(player)}
    <input
      class="form-check-input"
      type="checkbox"
      on:click={onCheck}
      checked={isChecked}
    />
  {/if}
</td>

<td style="background-color: {column.bgColorCode(player)}">
  <div id={`player-${player.player_info.name}`} class="td-string omit">
    {#if accountID === 0}
      {column.displayValue(player)}
    {:else}
      {#if alertPlayer}
        <i class="bi {alertPlayer.pattern}" />
      {/if}

      <!-- svelte-ignore a11y-invalid-attribute -->
      <a
        class="td-link dropdown-toggle"
        href="#"
        id="dropdownMenuLink"
        data-bs-toggle="dropdown"
      >
        {column.displayValue(player)}
      </a>

      <ul class="dropdown-menu" aria-labelledby="dropdownMenuLink">
        {#if player.player_info.clan.id}
          <!-- svelte-ignore a11y-invalid-attribute -->
          <li>
            <a
              class="dropdown-item"
              href="#"
              on:click={() => BrowserOpenURL(column.clanURL(player))}
              >クラン詳細(WoWS Stats & Numbers)</a
            >
          </li>
        {/if}
        <!-- svelte-ignore a11y-invalid-attribute -->
        <li>
          <a
            class="dropdown-item"
            href="#"
            on:click={() => BrowserOpenURL(column.playerURL(player))}
            >プレイヤー詳細(WoWS Stats & Numbers)</a
          >
        </li>
        <!-- svelte-ignore a11y-invalid-attribute -->
        <li>
          {#if alertPlayer}
            <a
              class="dropdown-item"
              href="#"
              on:click={() => {
                dispatch("EditAlertPlayer", { target: clone(alertPlayer) });
              }}>プレイヤーリストを編集する</a
            >
            <a
              class="dropdown-item"
              href="#"
              on:click={() =>
                dispatch("RemoveAlertPlayer", { target: clone(alertPlayer) })}
              >プレイヤーリストから削除する</a
            >
          {:else}
            <a
              class="dropdown-item"
              href="#"
              on:click={() => {
                dispatch("EditAlertPlayer", {
                  target: {
                    account_id: player.player_info.id,
                    name: player.player_info.name,
                    pattern: "bi-check-circle-fill",
                    message: "",
                  },
                });
              }}>プレイヤーリストへ追加する</a
            >
          {/if}
        </li>
      </ul>
    {/if}
  </div>
  {#if alertPlayer?.message}
    <Tooltip target={`player-${player.player_info.name}`} placement="top">
      {alertPlayer.message}
    </Tooltip>
  {/if}
</td>
