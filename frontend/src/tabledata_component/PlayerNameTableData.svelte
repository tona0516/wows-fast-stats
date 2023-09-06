<script lang="ts">
  import clone from "clone";
  import type { PlayerName } from "src/lib/column/PlayerName";
  import { storedAlertPlayers } from "src/stores";
  import { createEventDispatcher } from "svelte";
  import type { domain } from "wailsjs/go/models";
  import { BrowserOpenURL } from "wailsjs/runtime/runtime";

  export let column: PlayerName;
  export let player: domain.Player;

  const dispatch = createEventDispatcher();

  $: alertPlayer = $storedAlertPlayers.find(
    (it) => it.account_id === player.player_info.id,
  );
</script>

<td
  class="td-string omit"
  style="background-color: {column.bgColorCode(player)}"
>
  {#if player.player_info.id === 0}
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
            on:click={() =>
              dispatch("RemoveAlertPlayer", { target: clone(alertPlayer) })}
            >プレイヤーリストから削除する</a
          >
        {:else}
          <a
            class="dropdown-item"
            href="#"
            on:click={() => {
              dispatch("UpdateAlertPlayer", {
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
      {#if alertPlayer}
        <li>
          <div class="dropdown-item">メモ: {alertPlayer.message}</div>
        </li>
      {/if}
    </ul>
  {/if}
</td>
