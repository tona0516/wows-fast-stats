<script lang="ts">
  import {
    showToast,
    storedPlayerDetailForShip,
    ShipDetailModal,
  } from "src/stores";
  import { BrowserOpenURL, ClipboardSetText } from "wailsjs/runtime/runtime";
  import ModalCommon from "./ModalCommon.svelte";
  import { NumbersURL } from "src/lib/NumbersURL";
  import { RATING_DEFS } from "src/lib/RatingLevel";
  import type { data } from "wailsjs/go/models";
  import { Color, type ColorPair } from "src/lib/Color";
  import { DispName } from "src/lib/DispName";

  interface DamageRating {
    level: string;
    color?: ColorPair;
    value: string;
  }

  $: damageRatings = getDamageRatings($storedPlayerDetailForShip);

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

    return RATING_DEFS.map((rating) => {
      const value = serverAvdgDamage * rating.damage;

      return {
        level: DispName.SKILL_LEVELS.get(rating.level) ?? "",
        color: Color.Rating.getFixed(rating.level),
        value: `${Math.floor(value).format(0)}~`,
      };
    });
  }
</script>

{#if $storedPlayerDetailForShip}
  {@const shipName = $storedPlayerDetailForShip.ship_info.name}
  {@const shipID = $storedPlayerDetailForShip.ship_info.id}

  <ModalCommon zValue={50} close={ShipDetailModal.close}>
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
        on:click={() => BrowserOpenURL(NumbersURL.ship(shipID))}
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
                    style="color: {dr.color?.text}">{dr.level}</td
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
