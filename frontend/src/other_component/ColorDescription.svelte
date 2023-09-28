<script lang="ts">
  import { ratingFactors } from "src/lib/rating/RatingConst";
  import { RatingConverterFactory } from "src/lib/rating/RatingConverter";
  import type { domain } from "wailsjs/go/models";
  import ExternalLink from "./ExternalLink.svelte";
  import { DispName } from "src/lib/DispName";

  export let userConfig: domain.UserConfig;
</script>

{#if userConfig.custom_color}
  <div class="uk-flex uk-flex-center">
    <div class="uk-overflow-auto">
      <table
        class="uk-table uk-table-shrink uk-table-divider uk-table-small uk-table-middle uk-text-nowrap"
      >
        <thead>
          <th class="uk-text-center">スキル</th>
          <th class="uk-text-center" colspan="2">PR</th>
          <th class="uk-text-center">ダメージ(平均比)</th>
          <th class="uk-text-center">勝率</th>
        </thead>
        <tbody>
          {#each ratingFactors() as factor}
            <tr>
              <td class="uk-text-center"
                >{DispName.SKILL_LEVELS.get(factor.level)}</td
              >
              <td
                class="uk-text-center"
                style="background-color: {RatingConverterFactory.fromPR(
                  factor.minPR,
                  userConfig,
                ).bgColorCode()};">player_name</td
              >
              <td
                class="uk-text-center"
                style="color: {RatingConverterFactory.fromPR(
                  factor.minPR,
                  userConfig,
                ).textColorCode()};">{factor.minPR} ~ {factor.maxPR}</td
              >
              <td
                class="uk-text-center"
                style="color: {RatingConverterFactory.fromDamage(
                  factor.minDamage,
                  1.0,
                  userConfig,
                ).textColorCode()};"
                >{factor.minDamage}倍 ~ {factor.maxDamage}倍</td
              >
              <td
                class="uk-text-center"
                style="color: {RatingConverterFactory.fromWinRate(
                  factor.minWin,
                  userConfig,
                ).textColorCode()};">{factor.minWin}% ~ {factor.maxWin}%</td
              >
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>

  <div class="uk-flex uk-flex-center">
    <div class="uk-margin-small-right">
      <ExternalLink url="https://asia.wows-numbers.com/personal/rating"
        >PRについて</ExternalLink
      >
    </div>
    <div class="uk-margin-small-right">
      <ExternalLink
        url="https://asia.wows-numbers.com/personal/rating/expected/preview/"
        >艦種別平均値について</ExternalLink
      >
    </div>
  </div>
{/if}
