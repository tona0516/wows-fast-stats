<script lang="ts">
  import type { data } from "wailsjs/go/models";
  import ExternalLink from "src/component/common/ExternalLink.svelte";
  import UkTable from "src/component/common/uikit/UkTable.svelte";
  import { DispName } from "src/lib/DispName";
  import { Rating } from "src/lib/Rating";
  import { RATING_DEFS } from "src/lib/RatingDef";

  export let config: data.UserConfigV2;

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

  const skillColorCode = config.color.skill.text;

  const descriptions = RATING_DEFS.map((rating) => {
    return {
      ratingText: DispName.SKILL_LEVELS.get(rating.level),
      playerName: {
        text: "player_name",
        textColor: Rating.fromPR(rating.pr.min, skillColorCode).colorCode(),
      },
      pr: {
        text: `${rating.pr.min} ~ ${rating.pr.max}`,
        textColor: Rating.fromDamage(
          rating.damage.min,
          1.0,
          skillColorCode,
        ).colorCode(),
      },
      damage: {
        text: `${rating.damage.min}倍 ~ ${rating.damage.max}倍`,
        textColor: Rating.fromDamage(
          rating.damage.min,
          1.0,
          skillColorCode,
        ).colorCode(),
      },
      win: {
        text: `${rating.winRate.min}% ~ ${rating.winRate.max}%`,
        textColor: Rating.fromWinRate(
          rating.winRate.min,
          skillColorCode,
        ).colorCode(),
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
