<script lang="ts">
  import { ratingFactors } from "src/lib/rating/RatingConst";
  import { RatingConverterFactory } from "src/lib/rating/RatingConverter";
  import type { domain } from "wailsjs/go/models";
  import ExternalLink from "./ExternalLink.svelte";

  export let userConfig: domain.UserConfig;
</script>

<div class="center mt-2 mx-2">
  {#if userConfig.custom_color}
    <table class="table table-sm table-text-color w-auto td-multiple">
      <thead>
        <th>スキル</th>
        <th colspan="2">PR</th>
        <th>ダメージ(平均比)</th>
        <th>勝率</th>
      </thead>
      <tbody>
        {#each ratingFactors() as factor}
          <tr>
            <td>{factor.level}</td>
            <td
              style="background-color: {RatingConverterFactory.fromPR(
                factor.minPR,
                userConfig,
              ).bgColorCode()};">player_name</td
            >
            <td
              style="color: {RatingConverterFactory.fromPR(
                factor.minPR,
                userConfig,
              ).textColorCode()};">{factor.minPR} ~ {factor.maxPR}</td
            >
            <td
              style="color: {RatingConverterFactory.fromDamage(
                factor.minDamage,
                1.0,
                userConfig,
              ).textColorCode()};"
              >{factor.minDamage}倍 ~ {factor.maxDamage}倍</td
            >
            <td
              style="color: {RatingConverterFactory.fromWinRate(
                factor.minWin,
                userConfig,
              ).textColorCode()};">{factor.minWin}% ~ {factor.maxWin}%</td
            >
          </tr>
        {/each}
      </tbody>
    </table>

    <ul>
      <li>
        <ExternalLink
          url="https://asia.wows-numbers.com/personal/rating"
          text="PRについて"
        />
      </li>
      <li>
        <ExternalLink
          url="https://asia.wows-numbers.com/personal/rating/expected/preview/"
          text="艦種別平均値について"
        />
      </li>
    </ul>
  {/if}
</div>
