<script lang="ts">
  import {
    storedAlertPlayers,
    storedPlayerDetail,
    showToast,
  } from "@libs/stores";
  import { ClipboardSetText, BrowserOpenURL } from "@wails/runtime/runtime";
  import ModalCommon from "./ModalCommon.svelte";
  import { clanNumbersURL, playerNumbersURL } from "@libs/utils";
  import { ModalManager } from "@libs/ModalManager";

  $: alertPlayer = $storedAlertPlayers.find(
    (ap) => ap.account_id === $storedPlayerDetail?.player_info.id,
  );
</script>

{#if $storedPlayerDetail}
  {@const accountID = $storedPlayerDetail.player_info.id}
  {@const playerName = $storedPlayerDetail.player_info.name}
  {@const clan = $storedPlayerDetail.player_info.clan}
  <ModalCommon zValue={50} close={ModalManager.instance.closeForPlayerDetail}>
    <h2 class="text-lg font-bold">
      <span>
        {#if clan.tag !== ""}
          [{clan.tag}]
        {/if}
        {playerName}
      </span>
      <button
        class="btn btn-circle btn-ghost"
        on:click={() => {
          ClipboardSetText(playerName);
          showToast("クリップボードにコピーしました！");
        }}
      >
        <i class="bi bi-clipboard"></i>
      </button>
    </h2>

    <div class="grid xl:grid-cols-1 gap-4 mt-2">
      {#if alertPlayer}
        {#if alertPlayer.message !== ""}
          <p>
            <i class="bi {alertPlayer.pattern}"></i>
            {alertPlayer.message}
          </p>
        {/if}

        <button
          class="btn btn-primary"
          on:click={() => ModalManager.instance.openForEdit(alertPlayer)}
          >アラートプレイヤーの編集</button
        >

        <button
          class="btn btn-error"
          on:click={() => {
            ModalManager.instance.openForRemove(alertPlayer);
          }}
        >
          アラートプレイヤーの削除
        </button>
      {:else}
        <button
          class="btn btn-primary"
          on:click={() =>
            ModalManager.instance.openForSpecify(accountID, playerName)}
          >アラートプレイヤーの編集</button
        >
      {/if}

      <button
        class="btn"
        on:click={() => BrowserOpenURL(playerNumbersURL(accountID, playerName))}
      >
        プレイヤーページ(wows-numbers.com)<i class="bi bi-box-arrow-in-up-right"
        ></i>
      </button>

      {#if clan.id !== 0}
        <button
          class="btn"
          on:click={() => BrowserOpenURL(clanNumbersURL(clan.id))}
        >
          クランページ(wows-numbers.com)<i class="bi bi-box-arrow-in-up-right"
          ></i>
        </button>
      {/if}
    </div>
  </ModalCommon>
{/if}
