<script lang="ts">
  import { showToast, storedPlayerShipDetail } from "@libs/stores";
  import type { data } from "@wails/go/models";
  import { ClipboardSetText, BrowserOpenURL } from "@wails/runtime/runtime";
  import ModalCommon from "./ModalCommon.svelte";
  import { shipNumbersURL } from "@libs/utils";
  import { ModalManager } from "@libs/ModalManager";
  import type { ColorCode } from "@libs/ColorCode";
  import { Rating } from "@libs/Rating";

  interface DamageRating {
    displayName: string;
    colorCode: ColorCode;
    value: string;
  }

  $: damageRatings = getDamageRatings($storedPlayerShipDetail);

  function getDamageRatings(
    player: data.Player | undefined,
  ): DamageRating[] | undefined {
    if (!player) {
      return undefined;
    }

    const serverAvdgDamage = player.ship_info.avg_damage;
    if (serverAvdgDamage === 0) {
      return undefined;
    }

    return Rating.getThresholds().map((threshold) => {
      const rating = new Rating(threshold.raw);
      const value = serverAvdgDamage * threshold.shipDamageRatio;

      return {
        displayName: rating.getDisplayName(),
        colorCode: rating.getColor().getFixedTextColor(),
        value: `${value.toFixed()}~`,
      };
    });
  }
</script>

{#if $storedPlayerShipDetail}
  {@const shipName = $storedPlayerShipDetail.ship_info.name}
  {@const shipID = $storedPlayerShipDetail.ship_info.id}

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
        on:click={() => BrowserOpenURL(shipNumbersURL(shipID))}
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
