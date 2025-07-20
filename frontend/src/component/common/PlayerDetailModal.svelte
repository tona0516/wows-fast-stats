<script lang="ts">
  import { NumbersURL } from "src/lib/NumbersURL";
  import {
    storedAlertPlayers,
    showUpdateAlertPlayerModal,
    showDeleteAlertPlayerModal,
    storedIsShowPlayerDetailModal,
    storedPlayerDetail,
    closeModal,
  } from "src/stores";
  import ExternalLink from "./ExternalLink.svelte";

  const isAlertPlayer = $storedAlertPlayers.some(
    (alertPlayer) => alertPlayer.account_id === $storedPlayerDetail.id,
  );

  function openUpdateAlertPlayerModal() {
    showUpdateAlertPlayerModal(
      $storedPlayerDetail.id,
      $storedPlayerDetail.name,
      $storedAlertPlayers[$storedPlayerDetail.id].pattern,
      $storedAlertPlayers[$storedPlayerDetail.id].message,
    );
  }

  function openDeleteAlertPlayerModal() {
    showDeleteAlertPlayerModal($storedPlayerDetail.id);
  }
</script>

{#if $storedIsShowPlayerDetailModal}
  <dialog class="modal modal-open">
    <form method="dialog" class="modal-box">
      <h2 class="text-lg font-bold">
        {#if $storedPlayerDetail.clan}
          [{$storedPlayerDetail.clan.tag}]
        {/if}
        {$storedPlayerDetail.name}
      </h2>

      <div class="mt-4">
        <ExternalLink
          url={NumbersURL.player(
            $storedPlayerDetail.id,
            $storedPlayerDetail.name,
          )}
        >
          プレイヤーページ<i class="bi bi-box-arrow-in-up-right"></i>
        </ExternalLink>
      </div>

      {#if $storedPlayerDetail.clan}
        <div class="mt-2">
          <ExternalLink url={NumbersURL.clan($storedPlayerDetail.clan.id)}>
            クランページ<i class="bi bi-box-arrow-in-up-right"></i>
          </ExternalLink>
        </div>
      {/if}

      <div class="mt-4">
        <button
          class="btn btn-primary"
          on:click={() => openUpdateAlertPlayerModal()}
        >
          アラートプレイヤーの編集
        </button>
      </div>

      {#if isAlertPlayer}
        <div class="mt-2">
          <button
            class="btn btn-error"
            on:click={() => openDeleteAlertPlayerModal()}
          >
            アラートプレイヤーの削除
          </button>
        </div>
      {/if}

      <div class="modal-action">
        <button type="button" class="btn" on:click={() => closeModal()}
          >キャンセル</button
        >
      </div>
    </form>
  </dialog>
{/if}
