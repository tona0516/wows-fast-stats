<script lang="ts">
  import type { model } from "wailsjs/go/models";
  import ExternalLink from "src/component/common/ExternalLink.svelte";
  import UkTable from "src/component/common/uikit/UkTable.svelte";
  import { DispName } from "src/lib/DispName";
  import {
    RATING_FACTORS,
    RatingColorFactory,
  } from "src/lib/rating/RatingColorFactory";

  export let config: model.UserConfig;

  const COLUMNS = [
    { text: "スキル", colspan: 1 },
    { text: "PR", colspan: 2 },
    { text: "ダメージ(平均比)", colspan: 1 },
    { text: "勝率", colspan: 1 },
  ];

  const LINKS = [
    {
      url: "https://asia.wows-numbers.com/personal/rating",
      text: "PRについて",
    },
    {
      url: "https://asia.wows-numbers.com/personal/rating/expected/preview/",
      text: "艦種別平均値について",
    },
  ];

  const descriptions = RATING_FACTORS.map((rating) => {
    return {
      ratingText: DispName.SKILL_LEVELS.get(rating.level),
      playerName: {
        text: "player_name",
        textColor: RatingColorFactory.fromPR(
          rating.pr.min,
          config,
        ).getTextColorCode(),
      },
      pr: {
        text: `${rating.pr.min} ~ ${rating.pr.max}`,
        textColor: RatingColorFactory.fromDamage(
          rating.damage.min,
          1.0,
          config,
        ).getTextColorCode(),
      },
      damage: {
        text: `${rating.damage.min}倍 ~ ${rating.damage.max}倍`,
        textColor: RatingColorFactory.fromDamage(
          rating.damage.min,
          1.0,
          config,
        ).getTextColorCode(),
      },
      win: {
        text: `${rating.winRate.min}% ~ ${rating.winRate.max}%`,
        textColor: RatingColorFactory.fromWinRate(
          rating.winRate.min,
          config,
        ).getTextColorCode(),
      },
    };
  });
</script>

{#if config.color}
  <div class="uk-flex uk-flex-center">
    <UkTable>
      <thead>
        {#each COLUMNS as column}
          <th class="uk-text-center" colspan={column.colspan}>{column.text}</th>
        {/each}
      </thead>

      <tbody>
        {#each descriptions as desc}
          <tr>
            <td class="uk-text-center">{desc.ratingText}</td>
            <td
              class="uk-text-center"
              style="color: {desc.playerName.textColor};"
              >{desc.playerName.text}</td
            >
            <td class="uk-text-center" style="color: {desc.pr.textColor};"
              >{desc.pr.text}</td
            >
            <td class="uk-text-center" style="color: {desc.damage.textColor};"
              >{desc.damage.text}</td
            >
            <td class="uk-text-center" style="color: {desc.win.textColor};"
              >{desc.win.text}</td
            >
          </tr>
        {/each}
      </tbody>
    </UkTable>
  </div>

  <div class="uk-flex uk-flex-center">
    {#each LINKS as link}
      <div class="uk-margin-small-right">
        <ExternalLink url={link.url}>{link.text}</ExternalLink>
      </div>
    {/each}
  </div>
{/if}
