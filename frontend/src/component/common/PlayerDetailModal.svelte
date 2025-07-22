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
  import ExternalLink from "./ExternalLink.svelte";
  import { ClipboardSetText } from "wailsjs/runtime/runtime";

  $: alertPlayer = $storedAlertPlayers.find(
    (ap) => ap.account_id === $storedPlayerDetail?.id,
  );
</script>

{#if $storedPlayerDetail}
  <dialog class="modal modal-open z-50">
    <form method="dialog" class="modal-box">
      <h2 class="text-lg font-bold">
        {#if $storedPlayerDetail.clan}
          [{$storedPlayerDetail.clan.tag}]
        {/if}
        {$storedPlayerDetail.name}
        <a
          href="#"
          class="mx-2"
          on:click={() => {
            ClipboardSetText($storedPlayerDetail.name);
            showToast("コピーしました！");
          }}
        >
          <i class="bi bi-clipboard"></i>
        </a>
      </h2>

      <div class="mt-4">
        <ExternalLink
          url={NumbersURL.player(
            $storedPlayerDetail.id,
            $storedPlayerDetail.name,
          )}
        >
          プレイヤーページ(wows-numbers.com)<i
            class="bi bi-box-arrow-in-up-right"
          ></i>
        </ExternalLink>
      </div>

      {#if $storedPlayerDetail.clan}
        <div class="mt-2">
          <ExternalLink url={NumbersURL.clan($storedPlayerDetail.clan.id)}>
            クランページ(wows-numbers.com)<i class="bi bi-box-arrow-in-up-right"
            ></i>
          </ExternalLink>
        </div>
      {/if}

      <div class="mt-4">
        <button
          class="btn btn-primary"
          on:click={() => {
            if (alertPlayer) {
              EditAlertPlayerModal.openForEdit(alertPlayer);
            } else {
              EditAlertPlayerModal.openForSpecify(
                $storedPlayerDetail.id,
                $storedPlayerDetail.name,
              );
            }
          }}
        >
          アラートプレイヤーの編集
        </button>
      </div>

      {#if alertPlayer}
        <div class="mt-2">
          <button
            class="btn btn-error"
            on:click={() => {
              DeleteAlertPlayerModal.open(alertPlayer);
            }}
          >
            アラートプレイヤーの削除
          </button>
        </div>
      {/if}

      <div class="modal-action">
        <button
          type="button"
          class="btn"
          on:click={() => PlayerDetailModal.close()}>キャンセル</button
        >
      </div>
    </form>
  </dialog>
{/if}
