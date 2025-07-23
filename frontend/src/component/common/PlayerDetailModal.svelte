<script lang="ts">
  import { NumbersURL } from "src/lib/NumbersURL";
  import {
    storedAlertPlayers,
    storedPlayerDetail,
    showToast,
    EditAlertPlayerModal,
    DeleteAlertPlayerModal,
    PlayerDetailModal,
  } from "src/stores";
  import { BrowserOpenURL, ClipboardSetText } from "wailsjs/runtime/runtime";
  import ModalCommon from "./ModalCommon.svelte";

  $: alertPlayer = $storedAlertPlayers.find(
    (ap) => ap.account_id === $storedPlayerDetail?.id,
  );
</script>

{#if $storedPlayerDetail}
  {@const accountID = $storedPlayerDetail.id}
  {@const playerName = $storedPlayerDetail.name}
  {@const clan = $storedPlayerDetail.clan}
  <ModalCommon zValue={50} close={PlayerDetailModal.close}>
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
          on:click={() => EditAlertPlayerModal.openForEdit(alertPlayer)}
          >アラートプレイヤーの編集</button
        >

        <button
          class="btn btn-error"
          on:click={() => {
            DeleteAlertPlayerModal.open(alertPlayer);
          }}
        >
          アラートプレイヤーの削除
        </button>
      {:else}
        <button
          class="btn btn-primary"
          on:click={() =>
            EditAlertPlayerModal.openForSpecify(accountID, playerName)}
          >アラートプレイヤーの編集</button
        >
      {/if}

      <button
        class="btn"
        on:click={() =>
          BrowserOpenURL(NumbersURL.player(accountID, playerName))}
      >
        プレイヤーページ(wows-numbers.com)<i class="bi bi-box-arrow-in-up-right"
        ></i>
      </button>

      {#if clan.id !== 0}
        <button
          class="btn"
          on:click={() => BrowserOpenURL(NumbersURL.clan(clan.id))}
        >
          クランページ(wows-numbers.com)<i class="bi bi-box-arrow-in-up-right"
          ></i>
        </button>
      {/if}
    </div>
  </ModalCommon>
{/if}
