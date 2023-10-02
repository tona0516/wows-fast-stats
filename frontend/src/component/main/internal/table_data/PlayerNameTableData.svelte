<script lang="ts">
  import ExternalLink from "src/component/common/ExternalLink.svelte";
  import UkDowndown from "src/component/common/uikit/UkDowndown.svelte";
  import UkIcon from "src/component/common/uikit/UkIcon.svelte";
  import UkTooltip from "src/component/common/uikit/UkTooltip.svelte";
  import type { PlayerName } from "src/lib/column/PlayerName";
  import { CssClass } from "src/lib/types";
  import { storedAlertPlayers, storedExcludePlayerIDs } from "src/stores";
  import { createEventDispatcher } from "svelte";
  import {
    AddExcludePlayerID,
    RemoveExcludePlayerID,
  } from "wailsjs/go/main/App";
  import type { domain } from "wailsjs/go/models";

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
  <div class="{CssClass.TD_STR} {CssClass.OMIT}">
    {#if isNPC}
      {column.displayValue(player)}
    {:else}
      <UkTooltip tooltip={alertPlayer?.message}>
        <!-- svelte-ignore a11y-invalid-attribute -->
        <a href="#">
          {#if alertPlayer}
            <i class="bi {alertPlayer.pattern}" />
          {/if}
          {column.displayValue(player)}
          <UkIcon name="chevron-down" />
        </a>
      </UkTooltip>
    {/if}
  </div>

  {#if !isNPC}
    <UkDowndown>
      <ul class="uk-nav uk-dropdown-nav">
        {#if isBelongToClan}
          <li>
            <ExternalLink url={column.clanURL(player)}
              >クラン詳細(WoWS Stats & Numbers)</ExternalLink
            >
          </li>
        {/if}

        <li>
          <ExternalLink url={column.playerURL(player)}
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
      </ul>
    </UkDowndown>
  {/if}
</td>