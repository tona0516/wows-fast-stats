<script lang="ts">
  import {
    storedPlayerDetail,
    showToast,
  } from "@libs/stores";
  import { ClipboardSetText, BrowserOpenURL } from "@wails/runtime/runtime";
  import ModalCommon from "./ModalCommon.svelte";
  import { ModalManager } from "@libs/ModalManager";
  import { NumbersURL } from "@libs/NumbersURL";
</script>

{#if $storedPlayerDetail}
  {@const accountID = $storedPlayerDetail.playerInfo.id}
  {@const playerName = $storedPlayerDetail.playerInfo.name}
  {@const clan = $storedPlayerDetail.playerInfo.clan}
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
      <button
        class="btn"
        on:click={() =>
          BrowserOpenURL(NumbersURL.getPlayer(accountID, playerName))}
      >
        プレイヤーページ(wows-numbers.com)<i class="bi bi-box-arrow-in-up-right"
        ></i>
      </button>

      {#if clan.id !== 0}
        <button
          class="btn"
          on:click={() => BrowserOpenURL(NumbersURL.getClan(clan.id))}
        >
          クランページ(wows-numbers.com)<i class="bi bi-box-arrow-in-up-right"
          ></i>
        </button>
      {/if}
    </div>
  </ModalCommon>
{/if}