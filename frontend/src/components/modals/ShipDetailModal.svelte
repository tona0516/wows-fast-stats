<script lang="ts">
  import { showToast, storedPlayerShipDetail } from "@libs/stores";
  import type { core } from "@wails/go/models";
  import { ClipboardSetText, BrowserOpenURL } from "@wails/runtime/runtime";
  import ModalCommon from "./ModalCommon.svelte";
  import { ModalManager } from "@libs/ModalManager";
  import type { ColorCode } from "@libs/ColorCode";
  import { NumbersURL } from "@libs/NumbersURL";
  import { RATING_COLORS, RATING_NAMES } from "@libs/constants";

  interface DamageRating {
    displayName: string;
    colorCode: ColorCode;
    value: string;
  }

  $: damageRatings = getDamageRatings($storedPlayerShipDetail);

  function getDamageRatings(
    player: core.Player | undefined,
  ): DamageRating[] | undefined {
    if (!player) {
      return undefined;
    }

    return player.warship.damageRatings?.map((dr) => {
      return {
        displayName: RATING_NAMES[dr.rating],
        colorCode: RATING_COLORS[dr.rating]?.getFixedTextColor(),
        value: `${dr.value.toFixed()}~`,
      };
    });
  }
</script>

{#if $storedPlayerShipDetail}
  {@const shipName = $storedPlayerShipDetail.warship.name}
  {@const shipID = $storedPlayerShipDetail.warship.id}

  <ModalCommon zValue={50} close={ModalManager.instance.closeShipDetail}>
    <h2 class="text-lg font-bold">
      <span>
        {shipName}
      </span>
      <button
        class="btn btn-xs btn-circle"
        on:click={() => {
          ClipboardSetText(shipName);
          showToast("クリップボードにコピーしました！");
        }}
      >
        <i class="bi bi-clipboard"></i>
      </button>
    </h2>

    <div class="grid xl:grid-cols-1 gap-4 mt-2">
      <button
        class="btn"
        on:click={() => BrowserOpenURL(NumbersURL.getShip(shipID))}
      >
        艦艇ページ(wows-numbers.com)<i class="bi bi-box-arrow-in-up-right"></i>
      </button>

      {#if damageRatings}
        <div class="flex flex-col items-center">
          <table class="table text-nowrap" style="width: 1px;">
            <thead>
              {#each ["ダメージレーティング", "値"] as column}
                <th class="p-1 text-center">{column}</th>
              {/each}
            </thead>
            <tbody>
              {#each damageRatings as dr}
                <tr>
                  <td
                    class="p-1 text-center font-bold"
                    style="color: {dr.colorCode?.raw}">{dr.displayName}</td
                  >
                  <td class="p-1 text-right">{dr.value}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}
    </div>
  </ModalCommon>
{/if}
